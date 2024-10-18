
// Dynamic Provisioning:

// Recommendation: Instead of manually creating PVs, leverage Kubernetes StorageClasses for dynamic provisioning. This approach is more scalable and manageable, especially in cloud environments.

apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: ebs-sc
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp2
  fsType: ext4
reclaimPolicy: Retain


// Update PVC to Use StorageClass

pvc := &corev1.PersistentVolumeClaim{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "terraform-state-pvc",
        Namespace: namespace,
    },
    Spec: corev1.PersistentVolumeClaimSpec{
        AccessModes: []corev1.PersistentVolumeAccessMode{
            corev1.ReadWriteOnce,
        },
        StorageClassName: func() *string {
            sc := "ebs-sc"
            return &sc
        }(),
        Resources: corev1.ResourceRequirements{
            Requests: corev1.ResourceList{
                corev1.ResourceStorage: resource.MustParse("1Gi"),
            },
        },
    },
}

