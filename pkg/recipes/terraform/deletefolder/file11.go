// AWS S3 eg. for state file eg.

apiVersion: batch/v1
kind: Job
metadata:
  name: terraform-apply-job
  namespace: default
spec:
  template:
    metadata:
      labels:
        job-name: terraform-apply-job
    spec:
      restartPolicy: Never
      serviceAccountName: terraform-executor-sa  # Ensure this SA has S3 upload permissions
      containers:
        - name: terraform
          image: hashicorp/terraform:light
          command: ["sh", "-c", "terraform init && terraform apply -auto-approve && aws s3 cp /app/terraform/terraform.tfstate s3://your-state-bucket/terraform.tfstate"]
          workingDir: /app/terraform
          volumeMounts:
            - name: terraform-config
              mountPath: /app/terraform
          envFrom:
            - secretRef:
                name: terraform-secret
      volumes:
        - name: terraform-config
          configMap:
            name: terraform-config-map
/*

func (e *executor) Deploy(ctx context.Context, options Options) error {
    // Existing deployment steps...

    // Wait for the Job to complete
    if err := e.waitForJobCompletion(ctx, jobName, options.Namespace); err != nil {
        // Fetch logs even if the Job failed for debugging
        logs, logErr := e.FetchLogs(ctx, jobName, options.Namespace)
        if logErr != nil {
            return fmt.Errorf("deployment failed: %v; additionally, failed to fetch logs: %v", err, logErr)
        }
        fmt.Printf("Deployment failed. Logs:\n%s\n", logs)
        return fmt.Errorf("deployment failed: %v", err)
    }

    // Download the state file from S3
    stateFilePath := "/path/to/downloaded/terraform.tfstate"
    err := downloadFromS3(ctx, "your-state-bucket", "terraform.tfstate", stateFilePath)
    if err != nil {
        return fmt.Errorf("failed to download state file: %v", err)
    }

    // Read and process the state file
    stateData, err := os.ReadFile(stateFilePath)
    if err != nil {
        return fmt.Errorf("failed to read state file: %v", err)
    }

    // Optionally, return the state data or parse it as needed
    // For example, you can parse it into a struct or extract specific outputs

    // Clean up Kubernetes resources
    if err := e.cleanupJobAndConfigMap(ctx, jobName, configMapName, secretName, options.Namespace); err != nil {
        return fmt.Errorf("cleanup failed: %v", err)
    }

    return nil
}

*/

// helper file download from s3

import (
    "context"
    "fmt"
    "os"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

func downloadFromS3(ctx context.Context, bucket, key, downloadPath string) error {
    // Initialize AWS session
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2"),
        // Credentials can be managed via environment variables or IAM roles
    })
    if err != nil {
        return fmt.Errorf("failed to create AWS session: %v", err)
    }

    // Create S3 service client
    svc := s3.New(sess)

    // Get the object
    obj, err := svc.GetObjectWithContext(ctx, &s3.GetObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
    })
    if err != nil {
        return fmt.Errorf("failed to get object from S3: %v", err)
    }
    defer obj.Body.Close()

    // Create the file
    outFile, err := os.Create(downloadPath)
    if err != nil {
        return fmt.Errorf("failed to create file: %v", err)
    }
    defer outFile.Close()

    // Write the body to file
    _, err = io.Copy(outFile, obj.Body)
    if err != nil {
        return fmt.Errorf("failed to copy object data to file: %v", err)
    }

    return nil
}
