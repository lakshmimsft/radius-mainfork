extension radius
extension testresources

@description('Specifies the location for resources.')
param location string = 'westus2'

@description('Specifies the oidc issuer URL.')
#disable-next-line no-hardcoded-env-urls
param oidcIssuer string = 'https://eastus.oic.prod-aks.azure.com/bdfed5a8-05f1-40c5-9539-813b88cae8fd/efb919f6-1337-4540-8e53-5093b23259fe/'


resource secretenv 'Applications.Core/environments@2023-10-01-preview' = {
  name: 'secret-azure-env5'
  location: location
  properties: {
    compute: {
      kind: 'kubernetes'
      resourceId: 'self'
      namespace: 'secret-azure-env5'
      identity: {
        kind: 'azure.com.workload'
        oidcIssuer: oidcIssuer
      }
    }
    providers: {
      azure: {
        scope: resourceGroup().id
      }
    }
    recipes: {
      'Test.Resources/secrets': {
        default: {
          templateKind: 'bicep'
          templatePath: 'lakacr4.azurecr.io/azsecretrecipe:latest'
        }
      }
    }
  }
}

resource secretsazureapp 'Applications.Core/applications@2023-10-01-preview' = {
  name: 'secretsazureappbicep5'
  location: location
  properties: {
    environment: secretenv.id
    extensions: [
      {
        kind: 'kubernetesNamespace'
        namespace: 'secretsazureappbicep5'
      }
    ]
  }
}

resource secret 'Test.Resources/secrets@2025-08-01-preview' = {
  name: 'lak-azure-secret'
  properties: {
    environment: secretenv.id
    application: secretsazureapp.id
    data: {
      username: {
        value: 'admin'
      }
      password: {
        value: 'c2VjcmV0cGFzc3dvcmQ='
        encoding: 'base64'
      }
      apikey: {
        value: 'abc123xyz'
      }
    }
  }
}
