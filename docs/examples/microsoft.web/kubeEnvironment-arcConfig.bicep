@description('Name of the Arc-enabled Kubernetes Environment.')
param kubeEnvironmentName string = 'arc-enabled-kubeenv'

@description('Azure region for the Kubernetes Environment.')
param location string = resourceGroup().location

@secure()
@description('Base64-encoded kubeconfig for the connected Arc cluster. Treat this as sensitive input and source it from a secure parameter file or key vault reference.')
param kubeConfig string

resource kubeEnv 'Microsoft.Web/kubeEnvironments@2025-03-01' = {
  name: kubeEnvironmentName
  location: location
  properties: {
    // environmentType is optional; include it to be explicit when targeting managed Container Apps scenarios.
    environmentType: 'Managed'
    arcConfiguration: {
      // kubeConfig is marked as {sensitive} in the RP contract. Supplying it via a secure parameter keeps it out of templates and state files.
      kubeConfig: kubeConfig
      artifactsStorageType: 'LocalNode'
      frontEndServiceConfiguration: {
        kind: 'LoadBalancer'
      }
    }
  }
}

output kubeEnvironmentId string = kubeEnv.id
