//go:build go1.18
// +build go1.18

// Licensed under the Apache License, Version 2.0 . See LICENSE in the repository root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package v20231001preview

import "time"

// APIVersionProperties - The properties of an API version.
type APIVersionProperties struct {
	// READ-ONLY; The status of the asynchronous operation.
	ProvisioningState *ProvisioningState
}

// APIVersionResource - The resource type for defining an API version of a resource type supported by the containing resource
// provider.
type APIVersionResource struct {
	// The resource-specific properties for this resource.
	Properties *APIVersionProperties

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// AwsAccessKeyCredentialProperties - AWS credential properties for Access Key
type AwsAccessKeyCredentialProperties struct {
	// REQUIRED; Access key ID for AWS identity
	AccessKeyID *string

	// REQUIRED; The AWS credential kind
	Kind *AWSCredentialKind

	// REQUIRED; Secret Access Key for AWS identity
	SecretAccessKey *string

	// REQUIRED; The storage properties
	Storage CredentialStoragePropertiesClassification

	// READ-ONLY; The status of the asynchronous operation.
	ProvisioningState *ProvisioningState
}

// GetAwsCredentialProperties implements the AwsCredentialPropertiesClassification interface for type AwsAccessKeyCredentialProperties.
func (a *AwsAccessKeyCredentialProperties) GetAwsCredentialProperties() *AwsCredentialProperties {
	return &AwsCredentialProperties{
		Kind: a.Kind,
		ProvisioningState: a.ProvisioningState,
	}
}

// AwsCredentialProperties - AWS Credential properties
type AwsCredentialProperties struct {
	// REQUIRED; The AWS credential kind
	Kind *AWSCredentialKind

	// READ-ONLY; The status of the asynchronous operation.
	ProvisioningState *ProvisioningState
}

// GetAwsCredentialProperties implements the AwsCredentialPropertiesClassification interface for type AwsCredentialProperties.
func (a *AwsCredentialProperties) GetAwsCredentialProperties() *AwsCredentialProperties { return a }

// AwsCredentialResource - Concrete tracked resource types can be created by aliasing this type using a specific property
// type.
type AwsCredentialResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// REQUIRED; The resource-specific properties for this resource.
	Properties AwsCredentialPropertiesClassification

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// AwsCredentialResourceTagsUpdate - The type used for updating tags in AwsCredentialResource resources.
type AwsCredentialResourceTagsUpdate struct {
	// Resource tags.
	Tags map[string]*string
}

// AwsIRSACredentialProperties - AWS credential properties for IAM Roles for Service Accounts (IRSA)
type AwsIRSACredentialProperties struct {
	// REQUIRED; The AWS credential kind
	Kind *AWSCredentialKind

	// REQUIRED; RoleARN for AWS IRSA identity
	RoleARN *string

	// REQUIRED; The storage properties
	Storage CredentialStoragePropertiesClassification

	// READ-ONLY; The status of the asynchronous operation.
	ProvisioningState *ProvisioningState
}

// GetAwsCredentialProperties implements the AwsCredentialPropertiesClassification interface for type AwsIRSACredentialProperties.
func (a *AwsIRSACredentialProperties) GetAwsCredentialProperties() *AwsCredentialProperties {
	return &AwsCredentialProperties{
		Kind: a.Kind,
		ProvisioningState: a.ProvisioningState,
	}
}

// AwsPlaneResource - The AWS plane resource
type AwsPlaneResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// REQUIRED; The resource-specific properties for this resource.
	Properties *AwsPlaneResourceProperties

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// AwsPlaneResourceListResult - The response of a AwsPlaneResource list operation.
type AwsPlaneResourceListResult struct {
	// REQUIRED; The AwsPlaneResource items on this page
	Value []*AwsPlaneResource

	// The link to the next page of items
	NextLink *string
}

// AwsPlaneResourceProperties - The Plane properties.
type AwsPlaneResourceProperties struct {
	// READ-ONLY; The status of the asynchronous operation.
	ProvisioningState *ProvisioningState
}

// AwsPlaneResourceTagsUpdate - The type used for updating tags in AwsPlaneResource resources.
type AwsPlaneResourceTagsUpdate struct {
	// Resource tags.
	Tags map[string]*string
}

// AzureCredentialProperties - The base properties of Azure Credential
type AzureCredentialProperties struct {
	// REQUIRED; The kind of Azure credential
	Kind *AzureCredentialKind

	// READ-ONLY; The status of the asynchronous operation.
	ProvisioningState *ProvisioningState
}

// GetAzureCredentialProperties implements the AzureCredentialPropertiesClassification interface for type AzureCredentialProperties.
func (a *AzureCredentialProperties) GetAzureCredentialProperties() *AzureCredentialProperties { return a }

// AzureCredentialResource - Represents Azure Credential Resource
type AzureCredentialResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// REQUIRED; The resource-specific properties for this resource.
	Properties AzureCredentialPropertiesClassification

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// AzureCredentialResourceTagsUpdate - The type used for updating tags in AzureCredentialResource resources.
type AzureCredentialResourceTagsUpdate struct {
	// Resource tags.
	Tags map[string]*string
}

// AzurePlaneResource - The Azure plane resource.
type AzurePlaneResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// REQUIRED; The resource-specific properties for this resource.
	Properties *AzurePlaneResourceProperties

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// AzurePlaneResourceListResult - The response of a AzurePlaneResource list operation.
type AzurePlaneResourceListResult struct {
	// REQUIRED; The AzurePlaneResource items on this page
	Value []*AzurePlaneResource

	// The link to the next page of items
	NextLink *string
}

// AzurePlaneResourceProperties - The Plane properties.
type AzurePlaneResourceProperties struct {
	// REQUIRED; The URL used to proxy requests.
	URL *string

	// READ-ONLY; The status of the asynchronous operation.
	ProvisioningState *ProvisioningState
}

// AzurePlaneResourceTagsUpdate - The type used for updating tags in AzurePlaneResource resources.
type AzurePlaneResourceTagsUpdate struct {
	// Resource tags.
	Tags map[string]*string
}

// AzureServicePrincipalProperties - The properties of Azure Service Principal credential storage
type AzureServicePrincipalProperties struct {
	// REQUIRED; clientId for ServicePrincipal
	ClientID *string

	// REQUIRED; secret for ServicePrincipal
	ClientSecret *string

	// REQUIRED; The kind of Azure credential
	Kind *AzureCredentialKind

	// REQUIRED; The storage properties
	Storage CredentialStoragePropertiesClassification

	// REQUIRED; tenantId for ServicePrincipal
	TenantID *string

	// READ-ONLY; The status of the asynchronous operation.
	ProvisioningState *ProvisioningState
}

// GetAzureCredentialProperties implements the AzureCredentialPropertiesClassification interface for type AzureServicePrincipalProperties.
func (a *AzureServicePrincipalProperties) GetAzureCredentialProperties() *AzureCredentialProperties {
	return &AzureCredentialProperties{
		Kind: a.Kind,
		ProvisioningState: a.ProvisioningState,
	}
}

// AzureWorkloadIdentityProperties - The properties of Azure Workload Identity credential storage
type AzureWorkloadIdentityProperties struct {
	// REQUIRED; clientId for WorkloadIdentity
	ClientID *string

	// REQUIRED; The kind of Azure credential
	Kind *AzureCredentialKind

	// REQUIRED; The storage properties
	Storage CredentialStoragePropertiesClassification

	// REQUIRED; tenantId for WorkloadIdentity
	TenantID *string

	// READ-ONLY; The status of the asynchronous operation.
	ProvisioningState *ProvisioningState
}

// GetAzureCredentialProperties implements the AzureCredentialPropertiesClassification interface for type AzureWorkloadIdentityProperties.
func (a *AzureWorkloadIdentityProperties) GetAzureCredentialProperties() *AzureCredentialProperties {
	return &AzureCredentialProperties{
		Kind: a.Kind,
		ProvisioningState: a.ProvisioningState,
	}
}

// ComponentsKhmx01SchemasGenericresourceAllof0 - Concrete proxy resource types can be created by aliasing this type using
// a specific property type.
type ComponentsKhmx01SchemasGenericresourceAllof0 struct {
	// The resource-specific properties for this resource.
	Properties map[string]any

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// CredentialStorageProperties - The base credential storage properties
type CredentialStorageProperties struct {
	// REQUIRED; The kind of credential storage
	Kind *CredentialStorageKind
}

// GetCredentialStorageProperties implements the CredentialStoragePropertiesClassification interface for type CredentialStorageProperties.
func (c *CredentialStorageProperties) GetCredentialStorageProperties() *CredentialStorageProperties { return c }

// ErrorAdditionalInfo - The resource management error additional info.
type ErrorAdditionalInfo struct {
	// READ-ONLY; The additional info.
	Info map[string]any

	// READ-ONLY; The additional info type.
	Type *string
}

// ErrorDetail - The error detail.
type ErrorDetail struct {
	// READ-ONLY; The error additional info.
	AdditionalInfo []*ErrorAdditionalInfo

	// READ-ONLY; The error code.
	Code *string

	// READ-ONLY; The error details.
	Details []*ErrorDetail

	// READ-ONLY; The error message.
	Message *string

	// READ-ONLY; The error target.
	Target *string
}

// ErrorResponse - Common error response for all Azure Resource Manager APIs to return error details for failed operations.
// (This also follows the OData error response format.).
type ErrorResponse struct {
	// The error object.
	Error *ErrorDetail
}

// GenericPlaneResource - The generic representation of a plane resource
type GenericPlaneResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// REQUIRED; The resource-specific properties for this resource.
	Properties *GenericPlaneResourceProperties

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// GenericPlaneResourceProperties - The properties of the generic representation of a plane resource.
type GenericPlaneResourceProperties struct {
	// READ-ONLY; The status of the asynchronous operation.
	ProvisioningState *ProvisioningState
}

// GenericResource - Represents resource data.
type GenericResource struct {
	// The resource-specific properties for this resource.
	Properties map[string]any

	// READ-ONLY; The name of resource
	Name *string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// InternalCredentialStorageProperties - Internal credential storage properties
type InternalCredentialStorageProperties struct {
	// REQUIRED; The kind of credential storage
	Kind *CredentialStorageKind

	// READ-ONLY; The name of secret stored.
	SecretName *string
}

// GetCredentialStorageProperties implements the CredentialStoragePropertiesClassification interface for type InternalCredentialStorageProperties.
func (i *InternalCredentialStorageProperties) GetCredentialStorageProperties() *CredentialStorageProperties {
	return &CredentialStorageProperties{
		Kind: i.Kind,
	}
}

// LocationProperties - The properties of a location.
type LocationProperties struct {
	// Address of a resource provider implementation.
	Address *string

	// Configuration for resource types supported by the location.
	ResourceTypes map[string]*LocationResourceType

	// READ-ONLY; The status of the asynchronous operation.
	ProvisioningState *ProvisioningState
}

// LocationResource - The resource type for defining a location of the containing resource provider. The location resource
// represents a logical location where the resource provider operates.
type LocationResource struct {
	// The resource-specific properties for this resource.
	Properties *LocationProperties

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// LocationResourceType - The configuration for a resource type in a specific location.
type LocationResourceType struct {
	// The configuration for API versions of a resource type supported by the location.
	APIVersions map[string]map[string]any
}

// PagedResourceProviderSummary - Paged collection of ResourceProviderSummary items
type PagedResourceProviderSummary struct {
	// REQUIRED; The ResourceProviderSummary items on this page
	Value []*ResourceProviderSummary

	// The link to the next page of items
	NextLink *string
}

// PlaneNameParameter - The Plane Name parameter.
type PlaneNameParameter struct {
	// REQUIRED; The name of the plane
	PlaneName *string
}

// ProxyResource - The resource model definition for a Azure Resource Manager proxy resource. It will not have tags and a
// location
type ProxyResource struct {
	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// RadiusPlaneResource - The Radius plane resource.
type RadiusPlaneResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// REQUIRED; The resource-specific properties for this resource.
	Properties *RadiusPlaneResourceProperties

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// RadiusPlaneResourceListResult - The response of a RadiusPlaneResource list operation.
type RadiusPlaneResourceListResult struct {
	// REQUIRED; The RadiusPlaneResource items on this page
	Value []*RadiusPlaneResource

	// The link to the next page of items
	NextLink *string
}

// RadiusPlaneResourceProperties - The Plane properties.
type RadiusPlaneResourceProperties struct {
	// REQUIRED; Resource Providers for UCP Native Plane
	ResourceProviders map[string]*string

	// READ-ONLY; The status of the asynchronous operation.
	ProvisioningState *ProvisioningState
}

// RadiusPlaneResourceTagsUpdate - The type used for updating tags in RadiusPlaneResource resources.
type RadiusPlaneResourceTagsUpdate struct {
	// Resource tags.
	Tags map[string]*string
}

// Resource - Common fields that are returned in the response for all Azure Resource Manager resources
type Resource struct {
	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// ResourceGroupProperties - The resource group resource properties
type ResourceGroupProperties struct {
	// READ-ONLY; The status of the asynchronous operation.
	ProvisioningState *ProvisioningState
}

// ResourceGroupResource - The resource group resource
type ResourceGroupResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// The resource-specific properties for this resource.
	Properties *ResourceGroupProperties

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// ResourceGroupResourceTagsUpdate - The type used for updating tags in ResourceGroupResource resources.
type ResourceGroupResourceTagsUpdate struct {
	// Resource tags.
	Tags map[string]*string
}

// ResourceListResult - The response of a Resource list operation.
type ResourceListResult struct {
	// REQUIRED; The Resource items on this page
	Value []*Resource

	// The link to the next page of items
	NextLink *string
}

// ResourceProviderProperties - The properties of a resource provider.
type ResourceProviderProperties struct {
	// READ-ONLY; The status of the asynchronous operation.
	ProvisioningState *ProvisioningState
}

// ResourceProviderResource - The resource type for defining a resource provider.
type ResourceProviderResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// The resource-specific properties for this resource.
	Properties *ResourceProviderProperties

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// ResourceProviderSummary - The summary of a resource provider configuration. This type is optimized for querying resource
// providers and supported types.
type ResourceProviderSummary struct {
	// REQUIRED; The resource provider locations.
	Locations map[string]map[string]any

	// REQUIRED; The resource provider name.
	Name *string

	// REQUIRED; The resource types supported by the resource provider.
	ResourceTypes map[string]*ResourceProviderSummaryResourceType
}

// ResourceProviderSummaryResourceType - A resource type and its versions.
type ResourceProviderSummaryResourceType struct {
	// REQUIRED; API versions supported by the resource type.
	APIVersions map[string]map[string]any

	// The default api version for the resource type.
	DefaultAPIVersion *string
}

// ResourceTypeProperties - The properties of a resource type.
type ResourceTypeProperties struct {
	// The default api version for the resource type.
	DefaultAPIVersion *string

	// READ-ONLY; The status of the asynchronous operation.
	ProvisioningState *ProvisioningState
}

// ResourceTypeResource - The resource type for defining a resource type supported by the containing resource provider.
type ResourceTypeResource struct {
	// The resource-specific properties for this resource.
	Properties *ResourceTypeProperties

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// SystemData - Metadata pertaining to creation and last modification of the resource.
type SystemData struct {
	// The timestamp of resource creation (UTC).
	CreatedAt *time.Time

	// The identity that created the resource.
	CreatedBy *string

	// The type of identity that created the resource.
	CreatedByType *CreatedByType

	// The timestamp of resource last modification (UTC)
	LastModifiedAt *time.Time

	// The identity that last modified the resource.
	LastModifiedBy *string

	// The type of identity that last modified the resource.
	LastModifiedByType *CreatedByType
}

// TrackedResource - The resource model definition for an Azure Resource Manager tracked top level resource which has 'tags'
// and a 'location'
type TrackedResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

