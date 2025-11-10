extension testresources
extension radius
extension kubernetes with {
  kubeConfig: ''
  namespace: 'postgres-secret3'
} as kubernetes

/*
param registry string

param version string

@description('PostgreSQL password')
@secure()
param password string = newGuid()
*/

resource udtenv 'Applications.Core/environments@2023-10-01-preview' = {
  name: 'postgres-secret3'
  location: 'global'
  properties: {
    compute: {
      kind: 'kubernetes'
      resourceId: 'self'
      namespace: 'postgres-secret3'
    }
    recipes: {
      'Test.Resources/postgres2': {
        default: {
          templateKind: 'terraform'
          templatePath: 'github.com/lakshmimsft/lak-temp-public//postgres-secret'
        }
      }
      'Test.Resources/secrets': {
        default: {
          templateKind: 'terraform'
          templatePath: 'github.com/lakshmimsft/lak-temp-public//hashivault-secret'
        }
      }
    }
  }
}

resource udtapp 'Applications.Core/applications@2023-10-01-preview' = {
  name: 'postgres-app3'
  location: 'global'
  properties: {
    environment: udtenv.id
    extensions: [
      {
        kind: 'kubernetesNamespace'
        namespace: 'postgres-app3'
      }
    ]
  }
}


resource udtpg 'Test.Resources/postgres2@2025-01-01-preview' = {
  name: 'postgres-new4'
  location: 'global'
  properties: {
    environment: udtenv.id
    application: udtapp.id
    connections: {
      secretstore: {
        source: secret.id
      }
    }
  }
}

resource secret 'Test.Resources/secrets@2025-08-01-preview' = {
  name: 'lak-postgres-secret-tf5'
  properties: {
    environment: udtenv.id
    application: udtapp.id
    data: {
      username: {
        value: 'admin'
      }
      password: {
        value: 'c2VjcmV0cGFzc3dvcmQ='
        encoding: 'base64'
      }
      dbname: {
        value: 'mydatabase'
      }
    }
  }
}
