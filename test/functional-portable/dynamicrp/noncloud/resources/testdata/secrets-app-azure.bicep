extension radius
extension testresources

@description('Specifies the location for resources.')
param location string = 'westus2'

resource secretenv 'Applications.Core/environments@2023-10-01-preview' = {
  name: 'secret-azure-env'
  location: location
  properties: {
    compute: {
      kind: 'kubernetes'
      resourceId: 'self'
      namespace: 'secret-azure-env'
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
  name: 'secretsazureappbicep'
  location: location
  properties: {
    environment: secretenv.id
    extensions: [
      {
        kind: 'kubernetesNamespace'
        namespace: 'secretsazureappbicep'
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
