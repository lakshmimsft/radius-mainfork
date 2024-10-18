// Define the PV and PVC Specifications

func definePersistentVolume() *corev1.PersistentVolume {
    pv := &corev1.PersistentVolume{
        ObjectMeta: metav1.ObjectMeta{
            Name: "terraform-state-pv",
        },
        Spec: corev1.PersistentVolumeSpec{
            Capacity: corev1.ResourceList{
                corev1.ResourceStorage: resource.MustParse("1Gi"),
            },
            AccessModes: []corev1.PersistentVolumeAccessMode{
                corev1.ReadWriteOnce,
            },
            PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimRetain,
            StorageClassName:             "manual", // Ensure your PVC uses this StorageClass
            HostPath: &corev1.HostPathVolumeSource{
                Path: "/mnt/data", // Path on the host node
            },
        },
    }
    return pv
}

func definePersistentVolumeClaim(namespace string) *corev1.PersistentVolumeClaim {
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
                sc := "manual"
                return &sc
            }(),
            Resources: corev1.ResourceRequirements{
                Requests: corev1.ResourceList{
                    corev1.ResourceStorage: resource.MustParse("1Gi"),
                },
            },
        },
    }
    return pvc
}
