Create RBAC Roles and Bindings

//Role for Namespace-Scoped Resources (PVCs, ConfigMaps, Secrets, Jobs)
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: default
  name: terraform-executor-role
rules:
- apiGroups: [""]
  resources: ["configmaps", "secrets", "persistentvolumeclaims", "pods/log"]
  verbs: ["get", "create", "update", "delete", "list"]
- apiGroups: ["batch"]
  resources: ["jobs"]
  verbs: ["get", "create", "update", "delete", "list"]

  // ClusterRole for Cluster-Scoped Resources (PersistentVolumes)
  apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: terraform-executor-clusterrole
rules:
- apiGroups: [""]
  resources: ["persistentvolumes"]
  verbs: ["get", "create", "update", "delete", "list"]

// RoleBinding and ClusterRoleBinding
  apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: terraform-executor-rolebinding
  namespace: default
subjects:
- kind: ServiceAccount
  name: terraform-executor-sa
  namespace: default
roleRef:
  kind: Role
  name: terraform-executor-role
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: terraform-executor-clusterrolebinding
subjects:
- kind: ServiceAccount
  name: terraform-executor-sa
  namespace: default
roleRef:
  kind: ClusterRole
  name: terraform-executor-clusterrole
  apiGroup: rbac.authorization.k8s.io

  //Replace default with your target namespace and ensure that the Service Account terraform-executor-sa exists in that namespace.

// Running the Application
Ensure that your application is running with the correct Service Account and that it can access the Kubernetes API. If running inside Kubernetes, associate your Pod with the terraform-executor-sa Service Account.
  apiVersion: v1
kind: Pod
metadata:
  name: terraform-deployer
  namespace: default
spec:
  serviceAccountName: terraform-executor-sa
  containers:
    - name: deployer
      image: your-deployer-image
      command: ["your-deployer-command"]
