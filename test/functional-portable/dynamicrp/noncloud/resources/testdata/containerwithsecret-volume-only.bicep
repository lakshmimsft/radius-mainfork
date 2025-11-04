@description('Radius context object passed into the recipe.')
param context object

extension kubernetes with {
  kubeConfig: ''
  namespace: context.runtime.kubernetes.namespace
} as kubernetes

var resourceName = context.resource.name
var namespace = context.runtime.kubernetes.namespace
var normalizedName = resourceName

// TESTING: Hardcoded Key Vault configuration
var keyVaultId = '/subscriptions/bb48e769-817e-4890-ab84-db911c92c791/resourceGroups/lak-azrg/providers/Microsoft.KeyVault/vaults/lak-azure-secret'
var keyVaultName = last(split(keyVaultId, '/'))
var keyVaultTenantId = subscription().tenantId
var keyVaultSecretNames = ['username', 'password', 'apikey']

// Workload identity configuration - REPLACE WITH YOUR VALUES
var workloadIdentityClientId = '<YOUR_MANAGED_IDENTITY_CLIENT_ID>'  // az identity show --name radius-keyvault-reader --resource-group lak-azrg --query clientId -o tsv
var workloadIdentityServiceAccount = 'workload-identity-sa'

// USE CASE 2: Volume Mount Only (No Environment Variables)
// Secrets are mounted as files at /mnt/secrets-store/<secretname>

// Create SecretProviderClass for Azure Key Vault CSI driver
resource secretProviderClass 'secrets-store.csi.x-k8s.io/SecretProviderClass@v1' = {
  metadata: {
    name: '${normalizedName}-kv-sync-volume'
    namespace: namespace
  }
  spec: {
    provider: 'azure'
    parameters: {
      usePodIdentity: 'false'
      useVMManagedIdentity: 'false'
      clientID: workloadIdentityClientId
      keyvaultName: keyVaultName
      tenantId: keyVaultTenantId
      objects: string({
        array: reduce(keyVaultSecretNames, [], (acc, secretName) => concat(acc, [{
          objectName: secretName
          objectType: 'secret'
        }]))
      })
    }
    // NOTE: NO secretObjects section - secrets will only be available as volume mount files
  }
}

// Simple deployment with volume mount from Key Vault
resource deployment 'apps/Deployment@v1' = {
  metadata: {
    name: normalizedName
    namespace: namespace
    labels: {
      app: normalizedName
      'azure.workload.identity/use': 'true'
    }
  }
  spec: {
    replicas: 1
    selector: {
      matchLabels: {
        app: normalizedName
      }
    }
    template: {
      metadata: {
        labels: {
          app: normalizedName
          'azure.workload.identity/use': 'true'
        }
      }
      spec: {
        serviceAccountName: workloadIdentityServiceAccount
        containers: [
          {
            name: 'app'
            image: 'nginx:latest'
            // NOTE: NO environment variables from secrets
            volumeMounts: [
              {
                name: 'secrets-store-inline'
                mountPath: '/mnt/secrets-store'
                readOnly: true
              }
            ]
          }
        ]
        volumes: [
          {
            name: 'secrets-store-inline'
            csi: {
              driver: 'secrets-store.csi.k8s.io'
              readOnly: true
              volumeAttributes: {
                secretProviderClass: '${normalizedName}-kv-sync-volume'
              }
            }
          }
        ]
      }
    }
  }
}

output result object = {
  resources: [
    '/planes/kubernetes/local/namespaces/${namespace}/providers/apps/Deployment/${normalizedName}'
    '/planes/kubernetes/local/namespaces/${namespace}/providers/secrets-store.csi.x-k8s.io/SecretProviderClass/${normalizedName}-kv-sync-volume'
  ]
}
