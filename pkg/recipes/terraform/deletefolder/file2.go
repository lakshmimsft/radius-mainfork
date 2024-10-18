import (
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/errors"
    "k8s.io/apimachinery/pkg/apis/meta/v1"
    resource "k8s.io/apimachinery/pkg/api/resource"
)

// Define the PersistentVolume and PersistentVolumeClaim

// createPersistentVolume creates a PV if it doesn't exist
func createPersistentVolume(ctx context.Context, clientset *kubernetes.Clientset, pv *corev1.PersistentVolume) error {
    existingPV, err := clientset.CoreV1().PersistentVolumes().Get(ctx, pv.Name, metav1.GetOptions{})
    if err != nil {
        if errors.IsNotFound(err) {
            // Create PV
            _, err = clientset.CoreV1().PersistentVolumes().Create(ctx, pv, metav1.CreateOptions{})
            if err != nil {
                return fmt.Errorf("failed to create PersistentVolume: %v", err)
            }
            fmt.Printf("PersistentVolume %s created.\n", pv.Name)
            return nil
        }
        return fmt.Errorf("failed to get PersistentVolume: %v", err)
    }

    // PV already exists
    fmt.Printf("PersistentVolume %s already exists.\n", existingPV.Name)
    return nil
}

// createPersistentVolumeClaim creates a PVC if it doesn't exist
func createPersistentVolumeClaim(ctx context.Context, clientset *kubernetes.Clientset, pvc *corev1.PersistentVolumeClaim) error {
    existingPVC, err := clientset.CoreV1().PersistentVolumeClaims(pvc.Namespace).Get(ctx, pvc.Name, metav1.GetOptions{})
    if err != nil {
        if errors.IsNotFound(err) {
            // Create PVC
            _, err = clientset.CoreV1().PersistentVolumeClaims(pvc.Namespace).Create(ctx, pvc, metav1.CreateOptions{})
            if err != nil {
                return fmt.Errorf("failed to create PersistentVolumeClaim: %v", err)
            }
            fmt.Printf("PersistentVolumeClaim %s created in namespace %s.\n", pvc.Name, pvc.Namespace)
            return nil
        }
        return fmt.Errorf("failed to get PersistentVolumeClaim: %v", err)
    }

    // PVC already exists
    fmt.Printf("PersistentVolumeClaim %s already exists in namespace %s.\n", existingPVC.Name, existingPVC.Namespace)
    return nil
}
