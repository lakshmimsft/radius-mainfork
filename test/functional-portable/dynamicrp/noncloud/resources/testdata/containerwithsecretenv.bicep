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

// Extract connection data from linked resources
var resourceConnections = context.resource.connections ?? {}

// // Extract Key Vault reference from connections if it exists
// var keyVaultConnections = reduce(items(resourceConnections), {}, (acc, conn) =>
//   contains(conn.value, 'id') && contains(string(conn.value.id), 'Microsoft.KeyVault/vaults')
//     ? union(acc, { '${conn.key}': conn.value.id })
//     : acc
// )

// TESTING: Hardcoded Key Vault ID
var keyVaultId = '/subscriptions/bb48e769-817e-4890-ab84-db911c92c791/resourceGroups/lak-azrg/providers/Microsoft.KeyVault/vaults/lak-azure-secret'
var hasKeyVault = true

// Get Key Vault details for CSI driver (no need to reference the resource, just extract the name)
var keyVaultName = last(split(keyVaultId, '/'))
var keyVaultTenantId = subscription().tenantId

// Define the secrets to mount from Key Vault (hardcoded for testing)
var keyVaultSecretNames = ['username', 'password', 'apikey']

// Workload identity configuration
var workloadIdentityClientId = '<YOUR_MANAGED_IDENTITY_CLIENT_ID>'  // Replace with output from: az identity show --name radius-keyvault-reader --resource-group lak-azrg --query clientId -o tsv
var workloadIdentityServiceAccount = 'workload-identity-sa'

var daprSidecar = resourceProperties.?extensions.?daprSidecar
var hasDaprSidecar = daprSidecar != null
var effectiveDaprAppId = hasDaprSidecar && daprSidecar.?appId != null && string(daprSidecar.?appId) != '' ? string(daprSidecar.?appId) : normalizedName
var podAnnotations = hasDaprSidecar ? union(
  {
    'dapr.io/enabled': 'true'
    'dapr.io/app-id': effectiveDaprAppId
  },
  daprSidecar.?appPort != null ? { 'dapr.io/app-port': string(daprSidecar.?appPort) } : {},
  (daprSidecar.?config != null && string(daprSidecar.?config) != '') ? { 'dapr.io/config': string(daprSidecar.?config) } : {}
) : {}

var environmentSegments = context.resource.properties.environment != null ? split(string(context.resource.properties.environment), '/') : []
var environmentLabel = length(environmentSegments) > 0 ? last(environmentSegments) : ''

// Labels
var labels = {
  'radapp.io/resource': resourceName
  'radapp.io/environment': environmentLabel
  'radapp.io/application': context.application == null ? '' : context.application.name
}

// Extract connection data from linked resources
var resourceConnections = context.resource.connections ?? {}

// Use replicas from properties, default to 1 if not specified
var replicaCount = resourceProperties.?replicas != null ? int(resourceProperties.replicas) : 1

// Build container specs once and partition into workload/init sets
var containerSpecs = reduce(containerItems, [], (acc, item) => concat(acc, [{
  isInit: item.value.?initContainer ?? false
  spec: union(
    {
      name: item.key
      image: item.value.image
    },
    // Add ports if they exist
    contains(item.value, 'ports') ? {
      ports: reduce(items(item.value.ports), [], (portAcc, port) => concat(portAcc, [{
        name: port.key
        containerPort: port.value.containerPort
        protocol: port.value.?protocol ?? 'TCP'
      }]))
    } : {},
    // Add environment variables from container spec and Key Vault secrets
    contains(item.value, 'env') || hasKeyVault ? {
      env: concat(
        // Add regular environment variables from container spec
        contains(item.value, 'env') ? reduce(items(item.value.env), [], (envAcc, envItem) => concat(envAcc, [union(
          {
            name: envItem.key
          },
          contains(envItem.value, 'value') ? { value: envItem.value.value } : {}
        )])) : [],
        // Add environment variables from Key Vault secrets via Kubernetes secret
        hasKeyVault ? reduce(keyVaultSecretNames, [], (acc, secretName) => concat(acc, [{
          name: toUpper(secretName)
          valueFrom: {
            secretKeyRef: {
              name: '${normalizedName}-kv-secret'
              key: secretName
            }
          }
        }])) : []
      )
    } : {},
    // Add volume mounts for CSI driver
    hasKeyVault || contains(item.value, 'volumeMounts') ? {
      volumeMounts: concat(
        // Add CSI volume mount for Key Vault
        hasKeyVault ? [{
          name: 'secrets-store-inline'
          mountPath: '/mnt/secrets-store'
          readOnly: true
        }] : [],
        // Add existing volume mounts
        contains(item.value, 'volumeMounts') ? reduce(item.value.volumeMounts, [], (vmAcc, vm) => concat(vmAcc, [{
          name: vm.volumeName
          mountPath: vm.mountPath
        }])) : []
      )
    } : {},
    // Add command if specified
    contains(item.value, 'command') ? { command: item.value.command } : {},
    // Add args if specified
    contains(item.value, 'args') ? { args: item.value.args } : {},
    // Add working directory if specified
    contains(item.value, 'workingDir') ? { workingDir: item.value.workingDir } : {},
    // Add resources if specified
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
    // Add liveness probe if specified and this is not an init container
    (!(item.value.?initContainer ?? false) && contains(item.value, 'livenessProbe')) ? {
      livenessProbe: union(
        contains(item.value.livenessProbe, 'exec') ? {
          exec: { command: item.value.livenessProbe.exec.command }
        } : {},
        contains(item.value.livenessProbe, 'httpGet') ? {
          httpGet: union(
            { port: item.value.livenessProbe.httpGet.port },
            contains(item.value.livenessProbe.httpGet, 'path') ? { path: item.value.livenessProbe.httpGet.path } : {},
            contains(item.value.livenessProbe.httpGet, 'scheme') ? { scheme: toUpper(string(item.value.livenessProbe.httpGet.scheme)) } : {},
            contains(item.value.livenessProbe.httpGet, 'httpHeaders') ? { httpHeaders: item.value.livenessProbe.httpGet.httpHeaders } : {}
          )
        } : {},
        contains(item.value.livenessProbe, 'tcpSocket') ? {
          tcpSocket: { port: item.value.livenessProbe.tcpSocket.port }
        } : {},
        contains(item.value.livenessProbe, 'initialDelaySeconds') ? { initialDelaySeconds: item.value.livenessProbe.initialDelaySeconds } : {},
        contains(item.value.livenessProbe, 'periodSeconds') ? { periodSeconds: item.value.livenessProbe.periodSeconds } : {},
        contains(item.value.livenessProbe, 'timeoutSeconds') ? { timeoutSeconds: item.value.livenessProbe.timeoutSeconds } : {},
        contains(item.value.livenessProbe, 'failureThreshold') ? { failureThreshold: item.value.livenessProbe.failureThreshold } : {},
        contains(item.value.livenessProbe, 'successThreshold') ? { successThreshold: item.value.livenessProbe.successThreshold } : {},
        contains(item.value.livenessProbe, 'terminationGracePeriodSeconds') ? { terminationGracePeriodSeconds: item.value.livenessProbe.terminationGracePeriodSeconds } : {}
      )
    } : {},
    // Add readiness probe if specified and this is not an init container
    (!(item.value.?initContainer ?? false) && contains(item.value, 'readinessProbe')) ? {
      readinessProbe: union(
        contains(item.value.readinessProbe, 'exec') ? {
          exec: { command: item.value.readinessProbe.exec.command }
        } : {},
        contains(item.value.readinessProbe, 'httpGet') ? {
          httpGet: union(
            { port: item.value.readinessProbe.httpGet.port },
            contains(item.value.readinessProbe.httpGet, 'path') ? { path: item.value.readinessProbe.httpGet.path } : {},
            contains(item.value.readinessProbe.httpGet, 'scheme') ? { scheme: toUpper(string(item.value.readinessProbe.httpGet.scheme)) } : {},
            contains(item.value.readinessProbe.httpGet, 'httpHeaders') ? { httpHeaders: item.value.readinessProbe.httpGet.httpHeaders } : {}
          )
        } : {},
        contains(item.value.readinessProbe, 'tcpSocket') ? {
          tcpSocket: { port: item.value.readinessProbe.tcpSocket.port }
        } : {},
        contains(item.value.readinessProbe, 'initialDelaySeconds') ? { initialDelaySeconds: item.value.readinessProbe.initialDelaySeconds } : {},
        contains(item.value.readinessProbe, 'periodSeconds') ? { periodSeconds: item.value.readinessProbe.periodSeconds } : {},
        contains(item.value.readinessProbe, 'timeoutSeconds') ? { timeoutSeconds: item.value.readinessProbe.timeoutSeconds } : {},
        contains(item.value.readinessProbe, 'failureThreshold') ? { failureThreshold: item.value.readinessProbe.failureThreshold } : {},
        contains(item.value.readinessProbe, 'successThreshold') ? { successThreshold: item.value.readinessProbe.successThreshold } : {},
        contains(item.value.readinessProbe, 'terminationGracePeriodSeconds') ? { terminationGracePeriodSeconds: item.value.readinessProbe.terminationGracePeriodSeconds } : {}
      )
    } : {}
  )
}]))

var podContainers = reduce(containerSpecs, [], (acc, container) => !(container.isInit ?? false) ? concat(acc, [container.spec]) : acc)
var podInitContainers = reduce(containerSpecs, [], (acc, container) => (container.isInit ?? false) ? concat(acc, [container.spec]) : acc)

// Add volume mounts
var volumeItems = items(resourceProperties.?volumes ?? {})
var podVolumes = concat(
  // Add CSI volume for Key Vault if connected
  hasKeyVault ? [{
    name: 'secrets-store-inline'
    csi: {
      driver: 'secrets-store.csi.k8s.io'
      readOnly: true
      volumeAttributes: {
        secretProviderClass: '${normalizedName}-kv-sync'
      }
    }
  }] : [],
  // Add regular volumes from resource properties
  reduce(volumeItems, [], (acc, vol) => concat(acc, [union(
    {
      name: vol.key
    },
    contains(vol.value, 'persistentVolume') ? union(
      (contains(vol.value.persistentVolume, 'claimName') && vol.value.persistentVolume.claimName != '') ? {
        persistentVolumeClaim: {
          claimName: vol.value.persistentVolume.claimName
        }
      } : {},
      (!(contains(vol.value.persistentVolume, 'claimName') && vol.value.persistentVolume.claimName != '') && contains(resourceConnections, vol.key) && (resourceConnections[vol.key].?status.?computedValues.?claimName ?? '') != '') ? {
        persistentVolumeClaim: {
          claimName: resourceConnections[vol.key].?status.?computedValues.?claimName
        }
      } : {}
    ) : {},
    contains(vol.value, 'secret') ? {
      secret: {
        secretName: vol.value.secret.secretName
      }
    } : {},
    contains(vol.value, 'emptyDir') ? {
      emptyDir: contains(vol.value.emptyDir, 'medium') ? {
        medium: vol.value.emptyDir.medium
      } : {}
    } : {}
  )]))
)

resource deployment 'apps/Deployment@v1' = {
  metadata: {
    name: normalizedName
    namespace: namespace
    labels: labels
  }
  spec: {
    replicas: replicaCount
    selector: {
      matchLabels: {
        'radapp.io/resource': resourceName
      }
    }
    template: {
      metadata: union(
        {
          labels: union(labels, hasKeyVault ? {
            'azure.workload.identity/use': 'true'
          } : {})
        },
        hasDaprSidecar ? { annotations: podAnnotations } : {}
      )
      spec: union(
        {
          containers: podContainers
        },
        hasKeyVault ? { serviceAccountName: workloadIdentityServiceAccount } : {},
        length(podInitContainers) > 0 ? { initContainers: podInitContainers } : {},
        length(podVolumes) > 0 ? { volumes: podVolumes } : {},
        contains(resourceProperties, 'restartPolicy') ? { restartPolicy: resourceProperties.restartPolicy } : {}
      )
    }
  }
}

// Build service ports - one service per container that has ports
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
    labels: union(labels, {
      container: svc.containerName
    })
  }
  spec: {
    type: 'ClusterIP'
    selector: {
      'radapp.io/resource': resourceName
    }
    ports: svc.ports
  }
}]

// Add Horizontal Pod Autoscaler (if autoScaling specified)
var autoScaling = resourceProperties.?autoScaling
var hasAutoScaling = autoScaling != null

// Build HPA metrics
var hpaMetrics = hasAutoScaling && contains(autoScaling, 'metrics') ? reduce(autoScaling.metrics, [], (acc, metric) => concat(acc, [union(
  {
    type: metric.kind == 'cpu' || metric.kind == 'memory' ? 'Resource' : 'External'
  },
  metric.kind == 'cpu' || metric.kind == 'memory' ? {
    resource: {
      name: metric.kind
      target: union(
        {
          type: contains(metric.target, 'averageUtilization') ? 'Utilization' : contains(metric.target, 'averageValue') ? 'AverageValue' : 'Value'
        },
        contains(metric.target, 'averageUtilization') ? { averageUtilization: metric.target.averageUtilization } : {},
        contains(metric.target, 'averageValue') ? { averageValue: string(metric.target.averageValue) } : {},
        contains(metric.target, 'value') ? { value: string(metric.target.value) } : {}
      )
    }
  } : {},
  (metric.kind == 'custom' && contains(metric, 'customMetric')) ? {
    external: {
      metric: {
        name: metric.customMetric
      }
      target: union(
        {
          type: contains(metric.target, 'averageUtilization') ? 'Utilization' : contains(metric.target, 'averageValue') ? 'AverageValue' : 'Value'
        },
        contains(metric.target, 'averageUtilization') ? { averageUtilization: metric.target.averageUtilization } : {},
        contains(metric.target, 'averageValue') ? { averageValue: string(metric.target.averageValue) } : {},
        contains(metric.target, 'value') ? { value: string(metric.target.value) } : {}
      )
    }
  } : {}
)])) : []

resource hpa 'autoscaling/HorizontalPodAutoscaler@v2' = if (hasAutoScaling) {
  metadata: {
    name: normalizedName
    namespace: namespace
    labels: labels
  }
  spec: {
    scaleTargetRef: {
      apiVersion: 'apps/v1'
      kind: 'Deployment'
      name: normalizedName
    }
    minReplicas: hasAutoScaling && contains(autoScaling, 'minReplicas') ? int(autoScaling.minReplicas) : replicaCount
    maxReplicas: hasAutoScaling && contains(autoScaling, 'maxReplicas') ? int(autoScaling.maxReplicas) : 10
    metrics: hpaMetrics
  }
}

// Create SecretProviderClass for Azure Key Vault CSI driver
resource secretProviderClass 'secrets-store.csi.x-k8s.io/SecretProviderClass@v1' = if (hasKeyVault) {
  metadata: {
    name: '${normalizedName}-kv-sync'
    namespace: namespace
  }
  spec: {
    provider: 'azure'
    parameters: {
      usePodIdentity: 'false'
      useVMManagedIdentity: 'false'
      clientID: workloadIdentityClientId
      keyvaultName: keyVaultName
      tenantId: keyVaultTenantId
      objects: string({
        array: reduce(keyVaultSecretNames, [], (acc, secretName) => concat(acc, [{
          objectName: secretName
          objectType: 'secret'
        }]))
      })
    }
    // Sync Key Vault secrets as Kubernetes secret for environment variable injection
    secretObjects: [
      {
        secretName: '${normalizedName}-kv-secret'
        type: 'Opaque'
        data: reduce(keyVaultSecretNames, [], (acc, secretName) => concat(acc, [{
          objectName: secretName
          key: secretName
        }]))
      }
    ]
  }
}

var deploymentResource = '/planes/kubernetes/local/namespaces/${namespace}/providers/apps/Deployment/${normalizedName}'
var serviceResources = reduce(servicesConfig, [], (acc, svc) => concat(acc, ['/planes/kubernetes/local/namespaces/${namespace}/providers/core/Service/${normalizedName}-${svc.containerName}']))
var hpaResource = hasAutoScaling ? ['/planes/kubernetes/local/namespaces/${namespace}/providers/autoscaling/HorizontalPodAutoscaler/${normalizedName}'] : []
var secretProviderResource = hasKeyVault ? ['/planes/kubernetes/local/namespaces/${namespace}/providers/secrets-store.csi.x-k8s.io/SecretProviderClass/${normalizedName}-kv-sync'] : []

var allResources = concat([deploymentResource], serviceResources, hpaResource, secretProviderResource)

output result object = {
  resources: allResources
}
