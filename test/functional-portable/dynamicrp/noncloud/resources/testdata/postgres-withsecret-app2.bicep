extension testresources
extension radius
extension kubernetes with {
  kubeConfig: ''
  namespace: 'postgres-secret5'
} as kubernetes

/*
param registry string

param version string

@description('PostgreSQL password')
@secure()
param password string = newGuid()
*/

resource udtenv 'Applications.Core/environments@2023-10-01-preview' = {
  name: 'postgres-secret5'
  location: 'global'
  properties: {
    compute: {
      kind: 'kubernetes'
      resourceId: 'self'
      namespace: 'postgres-secret5'
    }
    recipes: {
      'Test.Resources/postgres2': {
        default: {
          templateKind: 'terraform'
          templatePath: 'github.com/lakshmimsft/lak-temp-public//postgres-secret2'
        }
      }
      'Test.Resources/secrets2': {
        default: {
          templateKind: 'terraform'
          templatePath: 'github.com/lakshmimsft/lak-temp-public//hashivault-secret'
        }
      }
    }
  }
}

resource udtapp 'Applications.Core/applications@2023-10-01-preview' = {
  name: 'postgres-app5'
  location: 'global'
  properties: {
    environment: udtenv.id
    extensions: [
      {
        kind: 'kubernetesNamespace'
        namespace: 'postgres-app5'
      }
    ]
  }
}


resource udtpg 'Test.Resources/postgres2@2025-01-01-preview' = {
  name: 'postgres-new5'
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

resource secret 'Test.Resources/secrets2@2025-08-01-preview' = {
  name: 'lak-postgres-secret-tf6'
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
