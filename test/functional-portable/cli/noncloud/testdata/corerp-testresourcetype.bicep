extension laktypes
extension radius

@description('Specifies the location for resources.')
param location string = 'global'


resource env 'Applications.Core/environments@2023-10-01-preview' = {
  name: 'lak-env-test7'
  location: 'global'
  properties: {
    compute: {
      kind: 'kubernetes'
      resourceId: 'self'
      namespace: 'lak-env-test7'
    }
    recipes: {
      'Lakshmi.Resources/lakTypeB': {
        default: {
          templateKind: 'bicep'
          //templatePath: 'ghcr.io/radius-project/test/testrecipes/test-bicep-recipes/rabbitmq-recipe:latest'
          templatePath: 'lakacr2.azurecr.io/laktypeb:latest'
        }
      }
    }
  }
}

resource app 'Applications.Core/applications@2023-10-01-preview' = {
  name: 'lak-test-app7'
  location: location
  properties: {
    environment: env.id
  }
}

resource laktype4 'Lakshmi.Resources/lakTypeB@2023-10-01-preview' = {
  name: 'lak-test-typeB7'
  location: location
  properties: {
    application: app.id
    environment: env.id
  }
}
