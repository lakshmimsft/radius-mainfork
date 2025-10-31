param context object

@description('Specifies the Azure location for resources.')
param location string = 'westus2'

// Extract properties from context
var secretKind = context.resource.properties.?kind ?? 'generic'
var secretData = context.resource.properties.data
var resourceName = context.resource.name
var appName = context.application == null ? '' : context.application.name

// Validate required fields for secret kinds
var missingFields = secretKind == 'certificate-pem' && (!contains(secretData, 'tls.crt') || !contains(secretData, 'tls.key'))
  ? 'certificate-pem secrets must contain keys `tls.crt` and `tls.key`'
  : secretKind == 'basicAuthentication' && (!contains(secretData, 'username') || !contains(secretData, 'password'))
  ? 'basicAuthentication secrets must contain keys `username` and `password`'
  : secretKind == 'azureWorkloadIdentity' && (!contains(secretData, 'clientId') || !contains(secretData, 'tenantId'))
  ? 'azureWorkloadIdentity secrets must contain keys `clientId` and `tenantId`'
  : secretKind == 'awsIRSA' && !contains(secretData, 'roleARN')
  ? 'awsIRSA secrets must contain key `roleARN`'
  : ''

var vaultName = length(missingFields) > 0 ? missingFields : resourceName

resource keyVault 'Microsoft.KeyVault/vaults@2023-07-01' = {
  name: vaultName
  location: location
  tags: {
    resource: resourceName
    app: appName
  }
  properties: {
    sku: {
      family: 'A'
      name: 'standard'
    }
    tenantId: subscription().tenantId
    enabledForDeployment: true
    enabledForTemplateDeployment: true
    enableRbacAuthorization: true
    accessPolicies: []
  }
}

// Create a secret for each key in secretData
// Note: Azure Key Vault requires separate secret resources, unlike K8s secrets
resource secrets 'Microsoft.KeyVault/vaults/secrets@2023-07-01' = [for item in items(secretData): {
  parent: keyVault
  name: item.key
  properties: {
    value: item.value.value
    contentType: contains(item.value, 'encoding') && item.value.encoding == 'base64'
      ? 'base64'
      : 'string'
  }
}]

output result object = {
  resources: [
    keyVault.id
  ]
}
