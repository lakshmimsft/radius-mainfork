# Radius.Security/secrets - Sensitive Data Handling Implementation

## Overview

This document outlines the approach for implementing sensitive data protection for the `Radius.Security/secrets` resource type handled by dynamic-rp.

## Background

### Current State

- **Radius.Security/secrets**: A new resource type defined in [resource-types-contrib](https://github.com/radius-project/resource-types-contrib) that stores sensitive data like tokens, passwords, keys, and certificates
- **Dynamic-RP**: A type-agnostic handler that provides CRUD operations for user-defined types without prior knowledge of their schemas
- **Recipe-based deployment**: Secrets can be associated with recipes to create various backends (Azure Key Vault, HashiCorp Vault, Kubernetes secrets, etc.)

### Problem Statement

Currently, dynamic-rp stores all resource properties (including sensitive `data` fields) in the database permanently. For `Radius.Security/secrets`, we need to:
1. Allow recipes to access sensitive data during deployment
2. Prevent sensitive data from persisting in Radius database after deployment
3. Maintain flexibility to support multiple secret backends via recipes

## Investigation Findings

### Existing Pattern: Applications.Core/secretStores

**Location**: [pkg/corerp/frontend/controller/secretstores/kubernetes.go:273-274](../../../pkg/corerp/frontend/controller/secretstores/kubernetes.go#L273-L274)

**Pattern**:
- Type: Imperative (no recipe)
- Flow: Synchronous frontend operation
- Nullification point: In `UpsertSecret()` UpdateFilter, **before** first DB save

```go
// Creates K8s secret immediately
ksecret.Data[k] = []byte(val)
updateRequired = true
// Remove secret from metadata before storing it to data store.
secret.Value = nil
```

**Limitation**: Tightly coupled to Kubernetes secrets only; cannot support recipe-based multi-backend deployment.

### Dynamic-RP Architecture

**Key Files**:
- Frontend: [pkg/dynamicrp/frontend/routes.go:99-100](../../../pkg/dynamicrp/frontend/routes.go#L99-L100)
- Backend Controller: [pkg/dynamicrp/backend/controller/dynamicresource.go:59-77](../../../pkg/dynamicrp/backend/controller/dynamicresource.go#L59-L77)
- Recipe Controller: [pkg/dynamicrp/backend/controller/putrecipe.go:49-56](../../../pkg/dynamicrp/backend/controller/putrecipe.go#L49-L56)
- Data Model: [pkg/dynamicrp/datamodel/dynamicresource.go:37-42](../../../pkg/dynamicrp/datamodel/dynamicresource.go#L37-L42)

**Flow**:
```
Frontend (Sync):
  1. DefaultAsyncPut receives request
  2. Validates and prepares resource
  3. Saves to DB via PrepareAsyncOperation (line 167)
  4. Queues async operation
  5. Returns 202 Accepted

Backend (Async):
  1. DynamicResourceController determines operation type
  2. Routes to RecipePutController
  3. Delegates to CreateOrUpdateResource controller
  4. Executes recipe with full resource data
  5. Saves final state to DB (line 120-129)
```

### Recipe Controller Flow

**Location**: [pkg/portableresources/backend/controller/createorupdateresource.go](../../../pkg/portableresources/backend/controller/createorupdateresource.go)

**Key Points**:
- Line 83-87: Loads configuration
- Line 92-97: Executes recipe if needed
- Line 103-106: Processes the resource
- Line 109-118: Updates recipe status
- **Line 120-129: Final save to database** ← Nullification must happen here

## Key Difference: secretStores vs Security/secrets

| Aspect | Applications.Core/secretStores | Radius.Security/secrets |
|--------|-------------------------------|------------------------|
| Operation Model | Synchronous/Imperative | Asynchronous/Recipe-based |
| Backend Support | Kubernetes only | Multiple (Azure KV, HashiCorp, K8s, etc.) |
| Nullification Point | Frontend UpdateFilter (before save) | Backend controller (after recipe) |
| Data Flow | Frontend creates secret → nullify → save | Frontend saves → recipe reads → recipe deploys → nullify → save |
| Exposure Window | None (nullified before DB) | Seconds to minutes (during recipe execution) |

## Recommended Approach

### Implementation Strategy

**Option: Backend Post-Recipe Nullification with Safeguards**

Nullify sensitive data in the backend controller **after** successful recipe execution, with proper cleanup handlers.

### Why This Approach?

1. ✅ **Recipe compatibility**: Recipe needs full data to deploy to various backends
2. ✅ **Type-agnostic**: Dynamic-rp remains unaware of specific property semantics
3. ✅ **Extensible**: Can support any recipe backend (Azure, AWS, HashiCorp, etc.)
4. ✅ **Follows existing patterns**: Similar to how recipe outputs are stored
5. ✅ **Pragmatic**: Gets feature working with acceptable risk profile

### Where to Implement

**Primary Location**: Post-processing in recipe controller after line 118

```go
// pkg/portableresources/backend/controller/createorupdateresource.go
// After line 118 (recipe status update)

// NEW: Clean sensitive data for Radius.Security/secrets
if isSecuritySecretsType(resource) {
    if err := nullifySecretData(resource); err != nil {
        return ctrl.Result{}, err
    }
}

// Existing: Save to database (line 120-129)
update := &database.Object{...}
```

**Helper Functions** (in dynamic-rp processor):

```go
// pkg/dynamicrp/backend/processor/secrets.go

func isSecuritySecretsType(resource rpv1.RadiusResourceModel) bool {
    baseResource := resource.GetBaseResource()
    return strings.EqualFold(baseResource.Type, "Radius.Security/secrets")
}

func nullifySecretData(resource rpv1.RadiusResourceModel) error {
    // Cast to DynamicResource to access Properties map
    dynResource, ok := resource.(*datamodel.DynamicResource)
    if !ok {
        return fmt.Errorf("expected DynamicResource, got %T", resource)
    }

    if dynResource.Properties != nil {
        if _, exists := dynResource.Properties["data"]; exists {
            dynResource.Properties["data"] = nil
        }
    }

    return nil
}
```

**Cleanup Safety**: Add defer block to ensure cleanup on failures

```go
// Track secret exposure duration
secretStoredAt := time.Now()

defer func() {
    if isSecuritySecretsType(resource) {
        // Always nullify, even on error
        _ = nullifySecretData(resource)

        // Best-effort save on error path
        if err != nil {
            _ = c.DatabaseClient().Save(ctx, &database.Object{
                Metadata: database.Metadata{ID: req.ResourceID},
                Data:     resource,
            })
        }

        duration := time.Since(secretStoredAt)
        logger.Info("Secret data lifecycle completed",
            "resourceID", req.ResourceID,
            "duration", duration.String())
    }
}()
```

## Security Considerations

### The Temporary Storage Risk

**Exposure Window**:
```
T0: Frontend saves: {data: {password: "secret123"}} ← PLAINTEXT IN DB
T1: Async job starts (seconds later)
T2: Recipe executes (minutes)
T3: Backend saves: {data: null} ← SECRETS REMOVED
```

**Risk**: Sensitive data exists in database for seconds to minutes.

### Potential Threats

1. **Database backup/snapshot** during T0-T3 contains plaintext
2. **Database read access** by services with DB permissions
3. **Audit logs** might capture queries with data
4. **Recipe failures** might leave data in DB
5. **Etcd storage** (Kubernetes) may not have encryption at rest

### Mitigation Strategies

#### Implemented in First Pass

1. **Automatic cleanup on failure**: Defer block ensures nullification on all exit paths
2. **Duration tracking**: Monitor how long secrets live in database
3. **Logging**: Audit when secrets are nullified
4. **RBAC**: Leverage existing database access controls

#### Infrastructure-Level (Required)

1. **Encryption at rest**: Enable etcd encryption (Kubernetes) or database encryption
2. **Network encryption**: TLS for all database connections
3. **Access controls**: Restrict database access to minimum required services

#### Future Enhancements (Optional)

1. **In-memory cache**: Store secrets in Redis/cache with TTL instead of primary DB
   - Pros: Never in persistent storage, automatic TTL
   - Cons: Additional infrastructure, race conditions

2. **Database encryption**: Encrypt sensitive fields before storage
   - Pros: Defense in depth
   - Cons: Key management complexity

3. **Streaming secrets**: Pass directly from frontend to backend without DB
   - Pros: Minimal DB exposure
   - Cons: Operation queue still persists data

### Security Posture Compared to secretStores

| Aspect | secretStores | Security/secrets |
|--------|--------------|------------------|
| DB exposure | None | Seconds to minutes |
| Encryption at rest | Required for K8s etcd | Required for Radius DB |
| Cleanup guarantee | Immediate | Deferred (on success/failure) |
| Backend flexibility | Single (K8s) | Multiple (recipe-based) |

### Compliance Considerations

Organizations should ensure:
- ✅ Database encryption at rest is enabled
- ✅ Database access is restricted via RBAC
- ✅ Backup strategies account for short-lived sensitive data
- ✅ Audit logging captures secret lifecycle events
- ✅ Monitoring alerts on abnormal secret storage durations

## Implementation Plan

### Phase 1: Core Functionality (First Pass)

1. **Add nullification logic** in recipe controller
   - Location: `pkg/portableresources/backend/controller/createorupdateresource.go`
   - Add type detection and nullification after recipe execution

2. **Add helper functions** in dynamic-rp processor
   - Location: `pkg/dynamicrp/backend/processor/secrets.go`
   - Implement `isSecuritySecretsType()` and `nullifySecretData()`

3. **Add safety mechanisms**
   - Defer block for cleanup on errors
   - Duration tracking
   - Audit logging

4. **Testing**
   - Unit tests for nullification logic
   - Integration tests for recipe-based deployment
   - Failure scenario tests (recipe fails, async timeout, etc.)

### Phase 2: Security Hardening (Future)

1. **Metrics and monitoring**
   - Secret storage duration metrics
   - Alert on unusually long exposure windows
   - Dashboard for secret lifecycle tracking

2. **Enhanced cleanup**
   - Background job to find and clean orphaned secrets
   - TTL-based automatic cleanup

3. **Documentation**
   - Security best practices guide
   - Encryption at rest setup instructions
   - Compliance checklist

### Phase 3: Advanced Options (Optional)

1. **Cache-based storage**
   - Evaluate Redis/Valkey for secret staging
   - Implement TTL-based cleanup
   - Handle cache failures gracefully

2. **Field-level encryption**
   - Key management setup
   - Encrypt data field before DB save
   - Decrypt in backend for recipe execution

## Open Questions for Review

1. **Encryption at rest**: Should we mandate it or document it as recommended?
2. **Backup strategy**: Should we exclude secrets from backups or document the risk?
3. **Audit requirements**: What level of audit logging is required?
4. **Failure handling**: What happens if cleanup fails? Retry? Alert?
5. **Migration path**: If we move to cache-based storage later, how do we migrate?

## References

- [Radius.Security/secrets type definition](https://github.com/radius-project/resource-types-contrib/blob/main/Security/secrets/secrets.yaml)
- [Applications.Core/secretStores implementation](../../../pkg/corerp/frontend/controller/secretstores/kubernetes.go)
- [Dynamic-RP architecture](../../../pkg/dynamicrp/)
- [Recipe controller](../../../pkg/portableresources/backend/controller/createorupdateresource.go)
