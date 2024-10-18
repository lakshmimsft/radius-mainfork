// helper methods
func (e *executor) generateConfig(ctx context.Context, options Options) error {
    // Implement your logic to generate Terraform configuration files
    // and write them to options.WorkingDir
    // Example:
    mainTF := `
    provider "aws" {
        region = "us-west-2"
    }

    resource "aws_instance" "example" {
        ami           = "ami-0c55b159cbfafe1f0"
        instance_type = "t2.micro"
    }

    output "instance_id" {
        value = aws_instance.example.id
    }
    `
    err := os.WriteFile(filepath.Join(options.WorkingDir, "main.tf"), []byte(mainTF), 0644)
    if err != nil {
        return fmt.Errorf("failed to write main.tf: %v", err)
    }

    return nil
}

func (e *executor) createConfigMap(ctx context.Context, configMapName string, options Options) error {
    clientset, err := getKubernetesClient()
    if err != nil {
        return fmt.Errorf("failed to get Kubernetes client: %v", err)
    }

    // Read all files in the working directory
    files, err := os.ReadDir(options.WorkingDir)
    if err != nil {
        return fmt.Errorf("failed to read working directory: %v", err)
    }

    data := make(map[string]string)
    for _, file := range files {
        if file.IsDir() {
            continue
        }
        filePath := filepath.Join(options.WorkingDir, file.Name())
        content, err := os.ReadFile(filePath)
        if err != nil {
            return fmt.Errorf("failed to read file %s: %v", filePath, err)
        }
        data[file.Name()] = string(content)
    }

    configMap := &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name:      configMapName,
            Namespace: options.Namespace,
        },
        Data: data,
    }

    // Create ConfigMap
    _, err = clientset.CoreV1().ConfigMaps(options.Namespace).Create(ctx, configMap, metav1.CreateOptions{})
    if err != nil {
        if errors.IsAlreadyExists(err) {
            fmt.Printf("ConfigMap %s already exists in namespace %s.\n", configMapName, options.Namespace)
            return nil
        }
        return fmt.Errorf("failed to create ConfigMap: %v", err)
    }

    fmt.Printf("ConfigMap %s created in namespace %s.\n", configMapName, options.Namespace)
    return nil
}

func (e *executor) createSecret(ctx context.Context, secretName string, data map[string]string, namespace string) error {
    clientset, err := getKubernetesClient()
    if err != nil {
        return fmt.Errorf("failed to get Kubernetes client: %v", err)
    }

    secretData := make(map[string][]byte)
    for key, value := range data {
        secretData[key] = []byte(value)
    }

    secret := &corev1.Secret{
        ObjectMeta: metav1.ObjectMeta{
            Name:      secretName,
            Namespace: namespace,
        },
        Data: secretData,
    }

    // Create Secret
    _, err = clientset.CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
    if err != nil {
        if errors.IsAlreadyExists(err) {
            fmt.Printf("Secret %s already exists in namespace %s.\n", secretName, namespace)
            return nil
        }
        return fmt.Errorf("failed to create Secret: %v", err)
    }

    fmt.Printf("Secret %s created in namespace %s.\n", secretName, namespace)
    return nil
}

func (e *executor) createTerraformJob(ctx context.Context, jobName, configMapName, secretName, action string, options Options) error {
    clientset, err := getKubernetesClient()
    if err != nil {
        return fmt.Errorf("failed to get Kubernetes client: %v", err)
    }

    // Define the Terraform command based on the action
    command := fmt.Sprintf("terraform init && terraform %s -auto-approve && terraform output -json > /app/terraform/outputs.json && aws s3 cp /app/terraform/outputs.json s3://your-outputs-bucket/outputs.json", action)

    job := &batchv1.Job{
        ObjectMeta: metav1.ObjectMeta{
            Name:      jobName,
            Namespace: options.Namespace,
        },
        Spec: batchv1.JobSpec{
            Template: corev1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: map[string]string{
                        "job-name": jobName,
                    },
                },
                Spec: corev1.PodSpec{
                    RestartPolicy: corev1.RestartPolicyNever,
                    ServiceAccountName: "terraform-executor-sa", // Ensure this SA has necessary permissions
                    Containers: []corev1.Container{
                        {
                            Name:    "terraform",
                            Image:   "hashicorp/terraform:light",
                            Command: []string{"sh", "-c", command},
                            WorkingDir: "/app/terraform",
                            VolumeMounts: []corev1.VolumeMount{
                                {
                                    Name:      "terraform-config",
                                    MountPath: "/app/terraform",
                                },
                            },
                            EnvFrom: []corev1.EnvFromSource{
                                {
                                    SecretRef: &corev1.SecretEnvSource{
                                        LocalObjectReference: corev1.LocalObjectReference{
                                            Name: secretName,
                                        },
                                    },
                                },
                            },
                        },
                    },
                    Volumes: []corev1.Volume{
                        {
                            Name: "terraform-config",
                            VolumeSource: corev1.VolumeSource{
                                ConfigMap: &corev1.ConfigMapVolumeSource{
                                    LocalObjectReference: corev1.LocalObjectReference{
                                        Name: configMapName,
                                    },
                                },
                            },
                        },
                        {
                            Name: "terraform-state",
                            PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
                                ClaimName: "terraform-state-pvc",
                            },
                        },
                    },
                },
            },
        },
    }

    // Create Job
    _, err = clientset.BatchV1().Jobs(options.Namespace).Create(ctx, job, metav1.CreateOptions{})
    if err != nil {
        if errors.IsAlreadyExists(err) {
            fmt.Printf("Job %s already exists in namespace %s.\n", jobName, options.Namespace)
            return nil
        }
        return fmt.Errorf("failed to create Job: %v", err)
    }

    fmt.Printf("Job %s created in namespace %s.\n", jobName, options.Namespace)
    return nil
}

func (e *executor) waitForJobCompletion(ctx context.Context, jobName, namespace string) error {
    clientset, err := getKubernetesClient()
    if err != nil {
        return fmt.Errorf("failed to get Kubernetes client: %v", err)
    }

    for {
        job, err := clientset.BatchV1().Jobs(namespace).Get(ctx, jobName, metav1.GetOptions{})
        if err != nil {
            return fmt.Errorf("failed to get Job: %v", err)
        }

        if job.Status.Succeeded > 0 {
            fmt.Printf("Job %s succeeded.\n", jobName)
            return nil
        }

        if job.Status.Failed > 0 {
            fmt.Printf("Job %s failed.\n", jobName)
            return fmt.Errorf("job %s failed", jobName)
        }

        fmt.Printf("Job %s is still running. Waiting...\n", jobName)
        time.Sleep(5 * time.Second)
    }
}

func (e *executor) FetchLogsWithRetry(ctx context.Context, jobName, namespace string, retries int, delay time.Duration) (string, error) {
    var logs string
    var err error

    for i := 0; i < retries; i++ {
        logs, err = e.FetchLogs(ctx, jobName, namespace)
        if err == nil {
            return logs, nil
        }

        // Log the retry attempt
        fmt.Printf("Attempt %d to fetch logs failed: %v. Retrying in %v...\n", i+1, err, delay)

        select {
        case <-ctx.Done():
            return "", fmt.Errorf("context canceled while fetching logs")
        case <-time.After(delay):
            // Continue to the next retry
        }
    }

    return "", fmt.Errorf("failed to fetch logs after %d attempts: last error: %v", retries, err)
}

func (e *executor) FetchLogs(ctx context.Context, jobName, namespace string) (string, error) {
    clientset, err := getKubernetesClient()
    if err != nil {
        return "", fmt.Errorf("failed to get Kubernetes client: %v", err)
    }

    // List Pods with label "job-name=jobName"
    podList, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
        LabelSelector: fmt.Sprintf("job-name=%s", jobName),
    })
    if err != nil {
        return "", fmt.Errorf("failed to list Pods for Job %s: %v", jobName, err)
    }

    if len(podList.Items) == 0 {
        return "", fmt.Errorf("no Pods found for Job %s", jobName)
    }

    // Assuming only one Pod per Job for simplicity
    pod := podList.Items[0]

    // Identify the container name, assuming "terraform" as defined earlier
    containerName := "terraform"

    // Fetch logs
    req := clientset.CoreV1().Pods(namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
        Container: containerName,
        Follow:    false, // set to true if you want to stream logs
    })

    podLogs, err := req.Do(ctx).Raw()
    if err != nil {
        return "", fmt.Errorf("failed to get logs for Pod %s: %v", pod.Name, err)
    }

    return string(podLogs), nil
}
