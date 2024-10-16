/*

// FetchLogs retrieves logs from the Job's Pod after completion
func (e *executor) FetchLogs(ctx context.Context, jobName, namespace string) (string, error) {
    // List Pods with label "job-name=jobName"
    podList, err := e.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
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
    req := e.clientset.CoreV1().Pods(namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
        Container: containerName,
        Follow:    false, // set to true if you want to stream logs
    })

    podLogs, err := req.Do(ctx).Raw()
    if err != nil {
        return "", fmt.Errorf("failed to get logs for Pod %s: %v", pod.Name, err)
    }

    return string(podLogs), nil
}


func (e *executor) Deploy(ctx context.Context, options Options) error {
    // Existing deployment steps...

    // Wait for the Job to complete
    if err := e.waitForJobCompletion(ctx, jobName, options.Namespace); err != nil {
        // Fetch logs even if the Job failed for debugging
        logs, logErr := e.FetchLogs(ctx, jobName, options.Namespace)
        if logErr != nil {
            return fmt.Errorf("deployment failed: %v; additionally, failed to fetch logs: %v", err, logErr)
        }
        // Handle logs (e.g., log them, return them, etc.)
        fmt.Printf("Deployment failed. Logs:\n%s\n", logs)
        return fmt.Errorf("deployment failed: %v", err)
    }

    // Fetch logs after successful completion
    logs, err := e.FetchLogs(ctx, jobName, options.Namespace)
    if err != nil {
        return fmt.Errorf("deployment succeeded but failed to fetch logs: %v", err)
    }

    // Optionally, log or process the logs
    fmt.Printf("Deployment succeeded. Logs:\n%s\n", logs)

    // Clean up resources
    if err := e.cleanupJobAndConfigMap(ctx, jobName, configMapName, secretName, options.Namespace); err != nil {
        return fmt.Errorf("cleanup failed: %v", err)
    }

    return nil
}

func (e *executor) FetchAllLogs(ctx context.Context, jobName, namespace string) (map[string]string, error) {
    podList, err := e.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
        LabelSelector: fmt.Sprintf("job-name=%s", jobName),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to list Pods for Job %s: %v", jobName, err)
    }

    if len(podList.Items) == 0 {
        return nil, fmt.Errorf("no Pods found for Job %s", jobName)
    }

    logsMap := make(map[string]string)

    for _, pod := range podList.Items {
        for _, container := range pod.Spec.Containers {
            req := e.clientset.CoreV1().Pods(namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
                Container: container.Name,
                Follow:    false,
            })

            podLogs, err := req.Do(ctx).Raw()
            if err != nil {
                logsMap[container.Name] = fmt.Sprintf("failed to get logs: %v", err)
                continue
            }

            logsMap[container.Name] = string(podLogs)
        }
    }

    return logsMap, nil
}

----------------------------

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

func (e *executor) Deploy(ctx context.Context, options Options) error {
    // Existing deployment steps...

    // Wait for the Job to complete
    if err := e.waitForJobCompletion(ctx, jobName, options.Namespace); err != nil {
        // Attempt to fetch logs with retries
        logs, logErr := e.FetchLogsWithRetry(ctx, jobName, options.Namespace, 3, 5*time.Second)
        if logErr != nil {
            return fmt.Errorf("deployment failed: %v; additionally, failed to fetch logs: %v", err, logErr)
        }
        // Handle logs (e.g., log them, return them, etc.)
        fmt.Printf("Deployment failed. Logs:\n%s\n", logs)
        return fmt.Errorf("deployment failed: %v", err)
    }

    // Fetch logs after successful completion
    logs, err := e.FetchLogsWithRetry(ctx, jobName, options.Namespace, 3, 5*time.Second)
    if err != nil {
        return fmt.Errorf("deployment succeeded but failed to fetch logs: %v", err)
    }

    // Optionally, log or process the logs
    fmt.Printf("Deployment succeeded. Logs:\n%s\n", logs)

    // Clean up resources
    if err := e.cleanupJobAndConfigMap(ctx, jobName, configMapName, secretName, options.Namespace); err != nil {
        return fmt.Errorf("cleanup failed: %v", err)
    }

    return nil
}

*/