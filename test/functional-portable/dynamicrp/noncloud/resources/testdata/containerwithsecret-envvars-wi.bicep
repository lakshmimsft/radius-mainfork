@description('Radius context object passed into the recipe.')
param context object

extension kubernetes with {
  kubeConfig: ''
  namespace: context.runtime.kubernetes.namespace
} as kubernetes

var resourceName = context.resource.name
var namespace = context.runtime.kubernetes.namespace
var normalizedName = resourceName
var resourceProperties = context.resource.properties ?? {}
var containerItems = items(resourceProperties.containers ?? {})

// TESTING: Hardcoded Key Vault configuration (will be replaced with connection data later)
var keyVaultId = '/subscriptions/bb48e769-817e-4890-ab84-db911c92c791/resourceGroups/lak-azrg/providers/Microsoft.KeyVault/vaults/lak-azure-secret'
var keyVaultSecretNames = ['username', 'password', 'apikey']
var workloadIdentityClientId = '8482c334-f8ca-4359-8634-38caff14fd4e'

// USE CASE 1: Environment Variables Only with Workload Identity
// Environment variables reference the Key Vault using secretKeyRef (secretId + key)
// Uses Azure Workload Identity instead of service principal credentials

// Build container specs with environment variables from the secret
var containerSpecs = reduce(containerItems, [], (acc, item) => concat(acc, [{
  isInit: item.value.?initContainer ?? false
  spec: union(
    {
      name: item.key
      image: item.value.image
    },
    // Add ports
    contains(item.value, 'ports') ? {
      ports: reduce(items(item.value.ports), [], (portAcc, port) => concat(portAcc, [{
        name: port.key
        containerPort: port.value.containerPort
        protocol: port.value.?protocol ?? 'TCP'
      }]))
    } : {},
    // Add environment variables: combine container env vars with secret references
    {
      env: concat(
        // Existing environment variables from container spec
        contains(item.value, 'env') ? reduce(items(item.value.env), [], (envAcc, envItem) => concat(envAcc, [union(
          {
            name: envItem.key
          },
          // Check if it has a direct value or valueFrom
          contains(envItem.value, 'value') ? { value: envItem.value.value } : {},
          contains(envItem.value, 'valueFrom') ? { valueFrom: envItem.value.valueFrom } : {}
        )])) : [],
        // Add secret references using secretKeyRef
        reduce(keyVaultSecretNames, [], (acc2, secretKey) => concat(acc2, [{
          name: toUpper(secretKey)
          valueFrom: {
            secretKeyRef: {
              name: normalizedName
              key: secretKey
            }
          }
        }]))
      )
    },
    // Add command, args, workingDir
    contains(item.value, 'command') ? { command: item.value.command } : {},
    contains(item.value, 'args') ? { args: item.value.args } : {},
    contains(item.value, 'workingDir') ? { workingDir: item.value.workingDir } : {},
    // Add CSI volume mount to sync secrets
    {
      volumeMounts: [
        {
          name: 'secrets-store'
          mountPath: '/mnt/secrets-store'
          readOnly: true
        }
      ]
    },
    // Add resources
    contains(item.value, 'resources') ? {
      resources: union(
        contains(item.value.resources, 'limits') ? {
          limits: union(
            contains(item.value.resources.limits, 'cpu') ? { cpu: item.value.resources.limits.cpu } : {},
            contains(item.value.resources.limits, 'memoryInMib') ? { memory: '${item.value.resources.limits.memoryInMib}Mi' } : {}
          )
        } : {},
        contains(item.value.resources, 'requests') ? {
          requests: union(
            contains(item.value.resources.requests, 'cpu') ? { cpu: item.value.resources.requests.cpu } : {},
            contains(item.value.resources.requests, 'memoryInMib') ? { memory: '${item.value.resources.requests.memoryInMib}Mi' } : {}
          )
        } : {}
      )
    } : {},
    // Add probes
    (!(item.value.?initContainer ?? false) && contains(item.value, 'livenessProbe')) ? {
      livenessProbe: union(
        contains(item.value.livenessProbe, 'httpGet') ? {
          httpGet: union(
            { port: item.value.livenessProbe.httpGet.port },
            contains(item.value.livenessProbe.httpGet, 'path') ? { path: item.value.livenessProbe.httpGet.path } : {},
            contains(item.value.livenessProbe.httpGet, 'scheme') ? { scheme: toUpper(string(item.value.livenessProbe.httpGet.scheme)) } : {}
          )
        } : {},
        contains(item.value.livenessProbe, 'initialDelaySeconds') ? { initialDelaySeconds: item.value.livenessProbe.initialDelaySeconds } : {},
        contains(item.value.livenessProbe, 'periodSeconds') ? { periodSeconds: item.value.livenessProbe.periodSeconds } : {},
        contains(item.value.livenessProbe, 'timeoutSeconds') ? { timeoutSeconds: item.value.livenessProbe.timeoutSeconds } : {},
        contains(item.value.livenessProbe, 'failureThreshold') ? { failureThreshold: item.value.livenessProbe.failureThreshold } : {},
        contains(item.value.livenessProbe, 'successThreshold') ? { successThreshold: item.value.livenessProbe.successThreshold } : {}
      )
    } : {},
    (!(item.value.?initContainer ?? false) && contains(item.value, 'readinessProbe')) ? {
      readinessProbe: union(
        contains(item.value.readinessProbe, 'httpGet') ? {
          httpGet: union(
            { port: item.value.readinessProbe.httpGet.port },
            contains(item.value.readinessProbe.httpGet, 'path') ? { path: item.value.readinessProbe.httpGet.path } : {},
            contains(item.value.readinessProbe.httpGet, 'scheme') ? { scheme: toUpper(string(item.value.readinessProbe.httpGet.scheme)) } : {}
          )
        } : {},
        contains(item.value.readinessProbe, 'initialDelaySeconds') ? { initialDelaySeconds: item.value.readinessProbe.initialDelaySeconds } : {},
        contains(item.value.readinessProbe, 'periodSeconds') ? { periodSeconds: item.value.readinessProbe.periodSeconds } : {},
        contains(item.value.readinessProbe, 'timeoutSeconds') ? { timeoutSeconds: item.value.readinessProbe.timeoutSeconds } : {}
      )
    } : {}
  )
}]))

var podContainers = reduce(containerSpecs, [], (acc, container) => !(container.isInit ?? false) ? concat(acc, [container.spec]) : acc)
var podInitContainers = reduce(containerSpecs, [], (acc, container) => (container.isInit ?? false) ? concat(acc, [container.spec]) : acc)
var replicaCount = resourceProperties.?replicas != null ? int(resourceProperties.replicas) : 1

// Extract Key Vault name from the full resource ID
var keyVaultName = last(split(keyVaultId, '/'))

// Build the YAML objects string for SecretProviderClass
var objectsYaml = 'array:\n${reduce(keyVaultSecretNames, '', (acc, secretKey) => '${acc}  - |\n    objectName: ${secretKey}\n    objectType: secret\n')}'

// Create SecretProviderClass for Azure Key Vault CSI driver with Workload Identity
resource secretProviderClass 'secrets-store.csi.x-k8s.io/SecretProviderClass@v1' = {
  metadata: {
    name: normalizedName
    namespace: namespace
    labels: {
      'radapp.io/resource': resourceName
    }
  }
  spec: {
    provider: 'azure'
    secretObjects: [
      {
        secretName: normalizedName
        type: 'Opaque'
        data: [
          for secretKey in keyVaultSecretNames: {
            objectName: secretKey
            key: secretKey
          }
        ]
      }
    ]
    parameters: {
      clientID: workloadIdentityClientId
      keyvaultName: keyVaultName
      cloudName: 'AzurePublicCloud'
      tenantId: 'bdfed5a8-05f1-40c5-9539-813b88cae8fd'
      objects: objectsYaml
    }
  }
}

// Deployment with secret environment variables and workload identity
resource deployment 'apps/Deployment@v1' = {
  metadata: {
    name: normalizedName
    namespace: namespace
    labels: {
      'radapp.io/resource': resourceName
    }
  }
  spec: {
    replicas: replicaCount
    selector: {
      matchLabels: {
        'radapp.io/resource': resourceName
      }
    }
    template: {
      metadata: {
        labels: union(
          {
            'radapp.io/resource': resourceName
          },
          {
            'azure.workload.identity/use': 'true'
          }
        )
      }
      spec: union(
        {
          serviceAccountName: 'workload-identity-sa'
          containers: podContainers
          volumes: [
            {
              name: 'secrets-store'
              csi: {
                driver: 'secrets-store.csi.k8s.io'
                readOnly: true
                volumeAttributes: {
                  secretProviderClass: normalizedName
                }
              }
            }
          ]
        },
        length(podInitContainers) > 0 ? { initContainers: podInitContainers } : {},
        contains(resourceProperties, 'restartPolicy') ? { restartPolicy: resourceProperties.restartPolicy } : {}
      )
    }
  }
}

// Build service ports
var servicesConfig = reduce(containerItems, [], (acc, item) =>
  contains(item.value, 'ports') && length(items(item.value.ports)) > 0 ? concat(acc, [{
    containerName: item.key
    ports: reduce(items(item.value.ports), [], (portAcc, port) => concat(portAcc, [{
      name: port.key
      port: port.value.containerPort
      targetPort: port.value.containerPort
      protocol: port.value.?protocol ?? 'TCP'
    }]))
  }]) : acc
)

resource services 'core/Service@v1' = [for svc in servicesConfig: {
  metadata: {
    name: '${normalizedName}-${svc.containerName}'
    namespace: namespace
    labels: {
      'radapp.io/resource': resourceName
      container: svc.containerName
    }
  }
  spec: {
    type: 'ClusterIP'
    selector: {
      'radapp.io/resource': resourceName
    }
    ports: svc.ports
  }
}]

var deploymentResource = '/planes/kubernetes/local/namespaces/${namespace}/providers/apps/Deployment/${normalizedName}'
var serviceResources = reduce(servicesConfig, [], (acc, svc) => concat(acc, ['/planes/kubernetes/local/namespaces/${namespace}/providers/core/Service/${normalizedName}-${svc.containerName}']))

output result object = {
  resources: concat([deploymentResource], serviceResources)
}
