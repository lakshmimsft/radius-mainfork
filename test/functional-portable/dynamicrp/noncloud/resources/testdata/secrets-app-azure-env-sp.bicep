extension radius
extension testresources

@description('Specifies the location for resources.')
param location string = 'westus2'


//@description('Specifies the oidc issuer URL.')
//#disable-next-line no-hardcoded-env-urls
//param oidcIssuer string = 'https://eastus.oic.prod-aks.azure.com/bdfed5a8-05f1-40c5-9539-813b88cae8fd/efb919f6-1337-4540-8e53-5093b23259fe/'
//param oidcIssuer string = 'https://radiusoidc.blob.core.windows.net/kubeoidc/'

resource secretenv 'Applications.Core/environments@2023-10-01-preview' = {
  name: 'secret-azure-env4'
  location: location
  properties: {
    compute: {
      kind: 'kubernetes'
      resourceId: 'self'
      namespace: 'secretsazureappbicep4-env'
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
      'Test.Resources/containers': {
        default: {
          templateKind: 'bicep'
          templatePath: 'lakacr4.azurecr.io/containerenv:latest'
        }
      }
    }
  }
}

resource secretsazureapp 'Applications.Core/applications@2023-10-01-preview' = {
  name: 'secretsazureappbicep4'
  location: location
  properties: {
    environment: secretenv.id
    extensions: [
      {
        kind: 'kubernetesNamespace'
        namespace: 'secretsazureappbicep4'
      }
    ]
  }
}

resource secret 'Test.Resources/secrets@2025-08-01-preview' = {
  name: 'lak-azure-secret-env4'
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

// Create a container that mounts the persistent volume
resource myContainer 'Test.Resources/containers@2025-08-01-preview' = {
  name: 'myapp4'
  properties: {
    environment: secretenv.id
    application: secretsazureapp.id
    connections: {
      data: {
        source: secret.id
      }
    }
    containers: {
      web: {
        image: 'nginx:alpine'
        command: ['/bin/sh', '-c']
        args: ['nginx -g "daemon off;"']
        workingDir: '/usr/share/nginx/html'
        ports: {
          http: {
            containerPort: 80
            protocol: 'TCP'
          }
        }
        env: {
          TEST_ENV: {
            value: 'testenv'
          }
          MAX_CONNECTIONS: {
            value: '100'
          }
          NGINX_HOST: {
            value: 'localhost'
          }
        }
        volumeMounts: [] 
        resources: {
          requests: {
            cpu: '0.1'       
            memoryInMib: 128   
          }
          limits: {
            cpu: '0.5'
            memoryInMib: 512
          }
        }
        livenessProbe: {
          httpGet: {
            path: '/'
            port: 80
            scheme: 'http'
          }
          initialDelaySeconds: 10
          periodSeconds: 30
          timeoutSeconds: 5
          failureThreshold: 3
          successThreshold: 1
        }
        readinessProbe: {
          httpGet: {
            path: '/'
            port: 80
          }
          initialDelaySeconds: 5
          periodSeconds: 10
        }
      }
      init: {
        initContainer: true
        image: 'busybox:latest'
        command: ['sh', '-c']
        args: ['echo "Initializing..." && sleep 5']
        workingDir: '/tmp'
        env: {
          INIT_MESSAGE: {
            value: 'Starting initialization'
          }
        }
        resources: {
          requests: {
            cpu: '0.1'
            memoryInMib: 64
          }
        }
      }
    }
    restartPolicy: 'Always'
    volumes: {}
    extensions: {}
    replicas: 1
    autoScaling: {
      maxReplicas: 3
      metrics: [
        {
          kind: 'cpu'
          target: {
            averageUtilization: 50
          }
        }
      ]
    }
  }
}
