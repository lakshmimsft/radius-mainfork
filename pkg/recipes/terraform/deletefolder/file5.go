// Ensure Proper Cleanup
func (e *executor) cleanupJobAndConfigMap(ctx context.Context, jobName, configMapName, secretName, namespace string) error {
    clientset, err := getKubernetesClient()
    if err != nil {
        return fmt.Errorf("failed to get Kubernetes client: %v", err)
    }

    // Delete Job
    err = clientset.BatchV1().Jobs(namespace).Delete(ctx, jobName, metav1.DeleteOptions{})
    if err != nil && !errors.IsNotFound(err) {
        return fmt.Errorf("failed to delete Job %s: %v", jobName, err)
    }
    fmt.Printf("Job %s deleted.\n", jobName)

    // Delete ConfigMap
    err = clientset.CoreV1().ConfigMaps(namespace).Delete(ctx, configMapName, metav1.DeleteOptions{})
    if err != nil && !errors.IsNotFound(err) {
        return fmt.Errorf("failed to delete ConfigMap %s: %v", configMapName, err)
    }
    fmt.Printf("ConfigMap %s deleted.\n", configMapName)

    // Delete Secret
    err = clientset.CoreV1().Secrets(namespace).Delete(ctx, secretName, metav1.DeleteOptions{})
    if err != nil && !errors.IsNotFound(err) {
        return fmt.Errorf("failed to delete Secret %s: %v", secretName, err)
    }
    fmt.Printf("Secret %s deleted.\n", secretName)

    // Optionally, delete PVC if it's no longer needed
    /*
    pvcName := "terraform-state-pvc"
    err = clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, pvcName, metav1.DeleteOptions{})
    if err != nil && !errors.IsNotFound(err) {
        return fmt.Errorf("failed to delete PersistentVolumeClaim %s: %v", pvcName, err)
    }
    fmt.Printf("PersistentVolumeClaim %s deleted.\n", pvcName)
    */

    // Optionally, delete PV if it's no longer needed (be cautious as PVs are cluster-scoped)
    /*
    pvName := "terraform-state-pv"
    err = clientset.CoreV1().PersistentVolumes().Delete(ctx, pvName, metav1.DeleteOptions{})
    if err != nil && !errors.IsNotFound(err) {
        return fmt.Errorf("failed to delete PersistentVolume %s: %v", pvName, err)
    }
    fmt.Printf("PersistentVolume %s deleted.\n", pvName)
    */

    return nil
}
