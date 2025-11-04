param context object

@description('Specifies the Azure location for resources.')
param location string = 'westus2'

//@description('Object ID of the service principal or managed identity to grant Key Vault access')
//param principalId string = '21844a91-2a15-4064-803a-2b93d9e62ad9' // From the error message - your SP object ID

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

// Grant Key Vault Secrets User role to the service principal/managed identity
// Role: Key Vault Secrets User (4633458b-17de-408a-b874-0445c86b69e6)
/*
resource keyVaultRoleAssignment 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  name: guid(keyVault.id, principalId, '4633458b-17de-408a-b874-0445c86b69e6')
  scope: keyVault
  properties: {
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', '4633458b-17de-408a-b874-0445c86b69e6')
    principalId: principalId
    principalType: 'ServicePrincipal'
  }
}
*/
output result object = {
  resources: [
    keyVault.id
  ]
  values: {
    vaultUri: keyVault.properties.vaultUri
  }
}
