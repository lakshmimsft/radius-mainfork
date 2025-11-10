Yes! Let me create a diagram showing how all these components interact in your workload identity setup:
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              Azure Entra ID                                      │
│                                                                                  │
│  ┌────────────────────────────────────────────────────────────────────────┐    │
│  │ Entra ID Application: TestSecretsInRadiusWI-radius-app                 │    │
│  │ Client ID: 8482c334-f8ca-4359-8634-38caff14fd4e                        │    │
│  │                                                                          │    │
│  │  Federated Credentials:                                                 │    │
│  │   ┌──────────────────────────────────────────────────────────────┐    │    │
│  │   │ Subject: system:serviceaccount:secretsazureappbicep5:        │    │    │
│  │   │          workload-identity-sa                                 │    │    │
│  │   │ Issuer: https://eastus.oic.prod-aks.azure.com/.../...       │    │    │
│  │   │ Audience: api://AzureADTokenExchange                         │    │    │
│  │   └──────────────────────────────────────────────────────────────┘    │    │
│  │                                                                          │    │
│  │  RBAC Role Assignment:                                                  │    │
│  │   - Role: "Key Vault Secrets User"                                     │    │
│  │   - Scope: /subscriptions/.../KeyVault/lak-azure-secret               │    │
│  └────────────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────────────┘
                                      ▲
                                      │ (3) Exchange JWT for
                                      │     Azure access token
                                      │
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         AKS Cluster (Kubernetes)                                 │
│                                                                                  │
│  ┌────────────────────────────────────────────────────────────────────────┐    │
│  │ Namespace: secretsazureappbicep5                                        │    │
│  │                                                                          │    │
│  │  ┌────────────────────────────────────────────────────────────┐        │    │
│  │  │ ServiceAccount: workload-identity-sa                        │        │    │
│  │  │ Annotations:                                                │        │    │
│  │  │   azure.workload.identity/client-id: 8482c334-...          │        │    │
│  │  └────────────────────────────────────────────────────────────┘        │    │
│  │                          │                                              │    │
│  │                          │ (1) Webhook sees annotation                 │    │
│  │                          ▼                                              │    │
│  │  ┌────────────────────────────────────────────────────────────┐        │    │
│  │  │ Pod: myapp5-77474fd944-wt4kb                               │        │    │
│  │  │ Labels: azure.workload.identity/use: "true"                │        │    │
│  │  │ ServiceAccountName: workload-identity-sa                   │        │    │
│  │  │                                                              │        │    │
│  │  │ ┌──────────────────────────────────────────────────────┐  │        │    │
│  │  │ │ Injected by Workload Identity Webhook:               │  │        │    │
│  │  │ │                                                        │  │        │    │
│  │  │ │ Volumes:                                              │  │        │    │
│  │  │ │   - name: azure-identity-token                        │  │        │    │
│  │  │ │     projected:                                         │  │        │    │
│  │  │ │       serviceAccountToken:                            │  │        │    │
│  │  │ │         audience: api://AzureADTokenExchange          │  │        │    │
│  │  │ │         path: azure-identity-token                    │  │        │    │
│  │  │ │                                                        │  │        │    │
│  │  │ │ VolumeMounts:                                         │  │        │    │
│  │  │ │   - /var/run/secrets/azure/tokens/                   │  │        │    │
│  │  │ │                                                        │  │        │    │
│  │  │ │ Environment Variables:                                │  │        │    │
│  │  │ │   AZURE_CLIENT_ID=8482c334-...                       │  │        │    │
│  │  │ │   AZURE_TENANT_ID=bdfed5a8-...                       │  │        │    │
│  │  │ │   AZURE_FEDERATED_TOKEN_FILE=                        │  │        │    │
│  │  │ │     /var/run/secrets/azure/tokens/azure-identity-... │  │        │    │
│  │  │ │   AZURE_AUTHORITY_HOST=                              │  │        │    │
│  │  │ │     https://login.microsoftonline.com/               │  │        │    │
│  │  │ └──────────────────────────────────────────────────────┘  │        │    │
│  │  │                                                              │        │    │
│  │  │ ┌──────────────────────────────────────────────────────┐  │        │    │
│  │  │ │ Container: web                                        │  │        │    │
│  │  │ │                                                        │  │        │    │
│  │  │ │ VolumeMounts:                                         │  │        │    │
│  │  │ │   - /mnt/secrets-store (CSI volume)                  │  │        │    │
│  │  │ │   - /var/run/secrets/azure/tokens/ ───────────┐      │  │        │    │
│  │  │ │                                                │      │  │        │    │
│  │  │ │ File: /var/run/secrets/azure/tokens/          │      │  │        │    │
│  │  │ │       azure-identity-token                     │      │  │        │    │
│  │  │ │ Contains: Kubernetes Service Account JWT ─────┘      │  │        │    │
│  │  │ │          (signed by AKS OIDC issuer)                 │  │        │    │
│  │  │ └──────────────────────────────────────────────────────┘  │        │    │
│  │  │                                                              │        │    │
│  │  │ ┌──────────────────────────────────────────────────────┐  │        │    │
│  │  │ │ Volume: secrets-store (CSI)                          │  │        │    │
│  │  │ │ Driver: secrets-store.csi.k8s.io                     │  │        │    │
│  │  │ │ SecretProviderClass: myapp5                          │  │        │    │
│  │  │ └──────────────────────────────────────────────────────┘  │        │    │
│  │  └────────────────────────────────────────────────────────────┘        │    │
│  │                          │                                              │    │
│  │                          │ (2) CSI driver mounts volume                 │    │
│  │                          ▼                                              │    │
│  │  ┌────────────────────────────────────────────────────────────┐        │    │
│  │  │ SecretProviderClass: myapp5                                 │        │    │
│  │  │ provider: azure                                             │        │    │
│  │  │ parameters:                                                 │        │    │
│  │  │   clientID: 8482c334-f8ca-4359-8634-38caff14fd4e          │        │    │
│  │  │   keyvaultName: lak-azure-secret                           │        │    │
│  │  │   tenantId: bdfed5a8-...                                   │        │    │
│  │  │   objects: [username, password, apikey]                    │        │    │
│  │  └────────────────────────────────────────────────────────────┘        │    │
│  └────────────────────────────────────────────────────────────────────────┘    │
│                                                                                  │
│  ┌────────────────────────────────────────────────────────────────────────┐    │
│  │ Node (kube-system namespace)                                            │    │
│  │                                                                          │    │
│  │  ┌────────────────────────────────────────────────────────────┐        │    │
│  │  │ Secrets Store CSI Driver                                    │        │    │
│  │  │ DaemonSet: csi-secrets-store-secrets-store-csi-driver      │        │    │
│  │  │                                                              │        │    │
│  │  │ CSIDriver Resource:                                         │        │    │
│  │  │   tokenRequests:                                            │        │    │
│  │  │   - audience: api://AzureADTokenExchange                   │        │    │
│  │  └────────────────────────────────────────────────────────────┘        │    │
│  │                          │                                              │    │
│  │                          │ (2a) Passes pod info + token                 │    │
│  │                          ▼                                              │    │
│  │  ┌────────────────────────────────────────────────────────────┐        │    │
│  │  │ Azure Key Vault Provider                                    │        │    │
│  │  │ DaemonSet: csi-azure-provider                              │        │    │
│  │  │                                                              │        │    │
│  │  │ Reads from pod:                                             │        │    │
│  │  │   - AZURE_CLIENT_ID                                        │        │    │
│  │  │   - AZURE_TENANT_ID                                        │        │    │
│  │  │   - AZURE_FEDERATED_TOKEN_FILE                            │        │    │
│  │  │   - Service Account JWT token                              │        │    │
│  │  └────────────────────────────────────────────────────────────┘        │    │
│  └────────────────────────────────────────────────────────────────────────┘    │
│                                      │                                           │
│                                      │ (3) Uses JWT to get Azure token          │
└──────────────────────────────────────┼───────────────────────────────────────────┘
                                       │
                                       ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         Azure Key Vault                                          │
│                                                                                  │
│  ┌────────────────────────────────────────────────────────────────────────┐    │
│  │ Key Vault: lak-azure-secret                                             │    │
│  │                                                                          │    │
│  │ Secrets:                                                                │    │
│  │   - username                                                            │    │
│  │   - password                                                            │    │
│  │   - apikey                                                              │    │
│  │                                                                          │    │
│  │ (4) Provider retrieves secrets using Azure access token                │    │
│  │ (5) Secrets mounted to pod at /mnt/secrets-store/                      │    │
│  │ (6) K8s Secret "myapp5" created with secret values                     │    │
│  │ (7) Container env vars (USERNAME, PASSWORD, APIKEY) reference secret   │    │
│  └────────────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────────────┘
Flow Summary:
Workload Identity Webhook sees pod with label azure.workload.identity/use: "true" and injects:
Azure identity token volume with Kubernetes JWT
Environment variables (AZURE_CLIENT_ID, etc.)
CSI Driver mounts the secrets-store volume and contacts the Azure Provider
Azure Provider reads:
Pod's environment variables
Kubernetes service account JWT from /var/run/secrets/azure/tokens/azure-identity-token
SecretProviderClass configuration (clientID, keyvaultName)
Azure Provider exchanges the Kubernetes JWT for an Azure access token by calling Azure Entra ID (which validates the federated credential)
Azure Provider uses the Azure access token to retrieve secrets from Key Vault
Secrets are mounted to pod at /mnt/secrets-store/ and synced to Kubernetes Secret myapp5
Container environment variables reference the Kubernetes Secret