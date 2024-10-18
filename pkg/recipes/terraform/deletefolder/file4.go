//Integrate PV and PVC Creation into the Deploy() Method

func (e *executor) Deploy(ctx context.Context, options Options) error {
    // Initialize Kubernetes client
    clientset, err := getKubernetesClient()
    if err != nil {
        return fmt.Errorf("failed to get Kubernetes client: %v", err)
    }

    // Define PV and PVC
    pv := definePersistentVolume()
    pvc := definePersistentVolumeClaim(options.Namespace)

    // Create PV if it doesn't exist
    if err := createPersistentVolume(ctx, clientset, pv); err != nil {
        return err
    }

    // Create PVC if it doesn't exist
    if err := createPersistentVolumeClaim(ctx, clientset, pvc); err != nil {
        return err
    }

    // Proceed with generating config, creating ConfigMap, Secret, and Job
    if err := e.generateConfig(ctx, options); err != nil {
        return fmt.Errorf("failed to generate config: %v", err)
    }

    configMapName := fmt.Sprintf("terraform-config-%s", options.RecipeName)
    if err := e.createConfigMap(ctx, configMapName, options); err != nil {
        return fmt.Errorf("failed to create ConfigMap: %v", err)
    }

    envVars, err := e.getEnvironmentVariables(ctx, options)
    if err != nil {
        return fmt.Errorf("failed to get environment variables: %v", err)
    }

    secretName := fmt.Sprintf("terraform-secret-%s", options.RecipeName)
    if err := e.createSecret(ctx, secretName, envVars, options.Namespace); err != nil {
        return fmt.Errorf("failed to create Secret: %v", err)
    }

    // Create Kubernetes Job
    jobName := fmt.Sprintf("terraform-apply-job-%s", options.RecipeName)
    if err := e.createTerraformJob(ctx, jobName, configMapName, secretName, "apply", options); err != nil {
        return fmt.Errorf("failed to create Kubernetes Job: %v", err)
    }

    // Wait for the Job to complete
    if err := e.waitForJobCompletion(ctx, jobName, options.Namespace); err != nil {
        logs, logErr := e.FetchLogsWithRetry(ctx, jobName, options.Namespace, 3, 5*time.Second)
        if logErr != nil {
            return fmt.Errorf("deployment failed: %v; additionally, failed to fetch logs: %v", err, logErr)
        }
        fmt.Printf("Deployment failed. Logs:\n%s\n", logs)
        return fmt.Errorf("deployment failed: %v", err)
    }

    // Fetch outputs from S3 or retrieve state file from PV
    // Depending on your implementation choice (remote backend or external upload)

    // Example: Download state file from PV (assuming it's mounted in your application Pod)
    // You might need to access the file via shared storage or another mechanism

    // Clean up Kubernetes resources
    if err := e.cleanupJobAndConfigMap(ctx, jobName, configMapName, secretName, options.Namespace); err != nil {
        return fmt.Errorf("cleanup failed: %v", err)
    }

    return nil
}
