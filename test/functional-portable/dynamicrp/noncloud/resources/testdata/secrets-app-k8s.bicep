extension radius
extension testresources

resource secretenv 'Applications.Core/environments@2023-10-01-preview' = {
  name: 'secret-k8-env-tf9'
  location: 'global'
  properties: {
    compute: {
      kind: 'kubernetes'
      resourceId: 'self'
      namespace: 'secret-k8-env-tf9'
    }
    recipes: {
      'Test.Resources/secrets': {
        default: {
          templateKind: 'terraform'
          templatePath: 'github.com/lakshmimsft/lak-temp-public//k8s-secret'
        }
      }
    }
  }
}

resource secretsk8sapp 'Applications.Core/applications@2023-10-01-preview' = {
  name: 'secretsk8sappbicep'
  location: 'global'
  properties: {
    environment: secretenv.id
    extensions: [
      {
        kind: 'kubernetesNamespace'
        namespace: 'secretsk8sappbicep'
      }
    ]
  }
}

resource secret 'Test.Resources/secrets@2025-08-01-preview' = {
  name: 'lak-k8s-secret-tf9'
  properties: {
    environment: secretenv.id
    application: secretsk8sapp.id
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
