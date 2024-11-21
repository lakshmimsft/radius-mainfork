/*
Copyright 2023 The Radius Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package terraform

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	install "github.com/hashicorp/hc-install"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/radius-project/radius/pkg/recipes/terraform/config/backends"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Deploy triggers a Kubernetes Job to run the Terraform commands for deploying the infrastructure.
func (e *executor) Deploy(ctx context.Context, options Options) (*tfjson.State, error) {
	// Validate options.
	if options.RootDir == "" || options.EnvRecipe.TemplatePath == "" || options.Namespace == "" {
		return nil, fmt.Errorf("missing required options")
	}

	// Check if the working directory exists.
	if _, err := os.Stat(options.RootDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("working directory does not exist")
	}

	// Install Terraform
	i := install.NewInstaller()
	tf, err := Install(ctx, i, options.RootDir)
	defer func() {
		if err := i.Remove(ctx); err != nil {
			log.Println(fmt.Errorf("failed to cleanup terraform installation: %s", err.Error()))
		}
	}()
	if err != nil {
		return nil, err
	}

	// Create Terraform config in the working directory
	kubernetesBackendSuffix, err := e.generateConfig(ctx, tf, options)
	if err != nil {
		return nil, err
	}

	if err != nil {
		// Create the Persistent Volume if it does not exist.
		if err := e.createPersistentVolume(ctx); err != nil {
			return nil, fmt.Errorf("failed to create Persistent Volume: %v", err)
		}
	}

	// Step 1: Create the 'terraform-config-pvc' if it doesn't exist
	_, err = e.k8sClientSet.CoreV1().PersistentVolumeClaims(options.Namespace).Get(ctx, "terraform-config-pvc", metav1.GetOptions{})
	if err != nil {
		err := e.createConfigPVC(ctx, options.Namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to create config PVC: %v", err)
		}
	}

	// Check if the Persistent Volume Claim exists.
	_, err = e.k8sClientSet.CoreV1().PersistentVolumeClaims(options.Namespace).Get(ctx, "terraform-pvc", metav1.GetOptions{})
	if err != nil {
		// Create the Persistent Volume Claim if it does not exist.
		if err = e.createPersistentVolumeClaim(ctx, options.Namespace); err != nil {
			return nil, fmt.Errorf("failed to create Persistent Volume Claim: %v", err)
		}
	}

	// Step 2: Create the data loader pod
	err = e.createDataLoaderPod(ctx, options.Namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create data loader pod: %v", err)
	}

	// Step 3: Wait for the data loader pod to be running
	err = e.waitForPodRunning(ctx, options.Namespace, "data-loader-pod")
	if err != nil {
		return nil, fmt.Errorf("data loader pod is not running: %v", err)
	}

	// Step 4: Copy data to the pod
	localPath := options.RootDir // Adjust based on your struct
	err = e.copyToPod(options.Namespace, "data-loader-pod", "data-loader", localPath, "/mnt/terraform-config")
	if err != nil {
		return nil, fmt.Errorf("failed to copy data to pod: %v", err)
	}

	// Step 5: Optionally verify the data in the pod
	stdout, stderr, err := e.execCommand(ctx, options.Namespace, "data-loader-pod", "data-loader", []string{"ls", "/mnt/terraform-config"})
	if err != nil {
		return nil, fmt.Errorf("failed to execute command in pod: %v", err)
	}
	fmt.Printf("Data in pod:\n%s\nErrors:\n%s\n", stdout, stderr)

	// Step 6: Delete the data loader pod if it's no longer needed
	err = e.k8sClientSet.CoreV1().Pods(options.Namespace).Delete(ctx, "data-loader-pod", metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return nil, fmt.Errorf("failed to delete data loader pod: %v", err)
	}

	// Create a Kubernetes Job to run the Terraform apply command.
	jobName := fmt.Sprintf("terraform-apply-job-%s", options.EnvRecipe.Name)

	if err := e.createTerraformJob(ctx, jobName, "apply ", options); err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes Job: %v", err)
	}

	// Wait for the Job to complete.
	if err := e.waitForJobCompletion(ctx, jobName, options.Namespace); err != nil {
		return nil, fmt.Errorf("deployment failed: %v", err)
	}

	// Validate that the terraform state file backend source exists.
	// Currently only Kubernetes secret backend is supported, which is created by Terraform as a part of Terraform apply.
	backendExists, err := backends.NewKubernetesBackend(e.k8sClientSet).ValidateBackendExists(ctx, backends.KubernetesBackendNamePrefix+kubernetesBackendSuffix)
	if err != nil {
		return nil, fmt.Errorf("error retrieving kubernetes secret for terraform state: %w", err)
	} else if !backendExists {
		return nil, fmt.Errorf("expected kubernetes secret for terraform state is not found")
	}

	state, err := e.retrieveTerraformOutputs(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("error retrieving terraform outputs: %v", err)
	}

	// Clean up resources.
	//if err := e.cleanupJobAndConfigMap(ctx, jobName, configMapName, options.Namespace); err != nil {
	//    return fmt.Errorf("cleanup failed: %v", err)
	//}

	return state, nil
}

// Delete triggers a Kubernetes Job to run the Terraform commands for destroying the infrastructure.
func (e *executor) Delete(ctx context.Context, options Options) error {
	// Validate options.
	if options.RootDir == "" || options.EnvRecipe.TemplatePath == "" || options.Namespace == "" {
		return fmt.Errorf("missing required options")
	}

	// Check if the working directory exists.
	if _, err := os.Stat(options.RootDir); os.IsNotExist(err) {
		return fmt.Errorf("working directory does not exist")
	}

	/*
		// Create a ConfigMap with the Terraform configuration.
		configMapName := fmt.Sprintf("terraform-config-%s", options.EnvRecipe.Name)
		if err := e.createConfigMap(ctx, configMapName, options); err != nil {
			return fmt.Errorf("failed to create ConfigMap: %v", err)
		}*/

	// Create a Kubernetes Job to run the Terraform destroy command.
	jobName := fmt.Sprintf("terraform-destroy-job-%s", options.EnvRecipe.Name)
	if err := e.createTerraformJob(ctx, jobName, "destroy", options); err != nil {
		return fmt.Errorf("failed to create Kubernetes Job: %v", err)
	}

	// Wait for the Job to complete.
	if err := e.waitForJobCompletion(ctx, jobName, options.Namespace); err != nil {
		return fmt.Errorf("deletion failed: %v", err)
	}

	// Retrieve Terraform outputs
	outputs, err := e.retrieveTerraformOutputs(ctx, options)
	if err != nil {
		return fmt.Errorf("failed to retrieve Terraform outputs: %v", err)
	}
	fmt.Printf("Terraform Outputs:\n%v\n", outputs)

	// Clean up resources.
	if err := e.cleanupJobAndConfigMap(ctx, jobName, options.Namespace); err != nil {
		return fmt.Errorf("cleanup failed: %v", err)
	}

	return nil
}

// createConfigMap creates a ConfigMap containing the Terraform configuration files.
func (e *executor) createConfigMap(ctx context.Context, configMapName string, options Options) error {
	configMapClient := e.k8sClientSet.CoreV1().ConfigMaps(options.Namespace)
	// Read all files in the working directory recursively.
	data := make(map[string]string)
	err := filepath.Walk(options.RootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories.
		if info.IsDir() {
			return nil
		}
		// Read file content.
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %v", path, err)
		}
		// Use relative path as key in ConfigMap data.
		relPath, err := filepath.Rel(options.RootDir, path)
		if err != nil {
			return err
		}
		data[relPath] = string(content)
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to read files: %v", err)
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: options.Namespace,
		},
		Data: data,
	}
	/*
		// Delete the existing ConfigMap if it exists
		err = configMapClient.Delete(ctx, configMapName, metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			return fmt.Errorf("failed to delete existing ConfigMap: %v", err)
		}
	*/
	// Create the ConfigMap in Kubernetes.
	_, err = configMapClient.Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create ConfigMap: %v", err)
	}

	return nil
}

// createTerraformJob creates a Kubernetes Job that runs Terraform commands.
func (e *executor) createTerraformJob(ctx context.Context, jobName, action string, options Options) error {
	jobClient := e.k8sClientSet.BatchV1().Jobs(options.Namespace)

	// action can be "apply" or "destroy"
	// command := fmt.Sprintf("cd /app/terraform && terraform init && terraform %s -auto-approve", action)
	tfCommand := fmt.Sprintf(`
	echo "Starting Terraform %s" &&
	echo "Environment Variables:" && env &&
    echo "Service Account Files:" && ls -l /var/run/secrets/kubernetes.io/serviceaccount/ &&
   	cd /app/terraform/deploy/ &&
	echo "Listing files in /app/terraform/deploy:" &&
	ls -l /app/terraform/deploy &&
	echo "Running terraform init" &&
	terraform init -reconfigure &&
	echo "Running terraform %s" &&
	terraform %s -auto-approve &&
	echo "Pulling state file" &&
    terraform state pull > /app/terraform/state/terraform.tfstate &&
	echo "Terraform %s completed" &&
    sleep 300`, action, action, action, action)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: options.Namespace,
			Labels: map[string]string{
				"app": "terraform",
			},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					ServiceAccountName: "terraform-sa",
					Containers: []corev1.Container{
						{
							Name:  "terraform",
							Image: "hashicorp/terraform:latest",
							Command: []string{
								"sh",
								"-c",
								tfCommand,
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "terraform-config",
									MountPath: "/app/terraform",
								},
								{
									Name:      "terraform-state",
									MountPath: "/app/terraform/state",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "terraform-config",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "terraform-config-pvc",
								},
							},
						},
						{
							Name: "terraform-state",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "terraform-pvc",
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyOnFailure,
				},
			},
		},
	}

	// Create the Job
	_, err := jobClient.Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes Job: %v", err)
	}

	return nil
}

// createTerraformJob creates a Kubernetes Job that runs Terraform commands.
func (e *executor) createTerraformJob_NotCalledAnymore(ctx context.Context, jobName, configMapName, action string, options Options) error {
	// action can be "apply" or "destroy"
	// command := fmt.Sprintf("cd /app/terraform && terraform init && terraform %s -auto-approve", action)
	command := fmt.Sprintf(`
	echo "Starting Terraform %s" &&
	cd /app/terraform &&
	echo "Listing files in /app/terraform:" &&
	ls -l /app/terraform &&
	echo "Running terraform init" &&
	terraform init &&
	echo "Running terraform %s" &&
	terraform %s -state=/app/terraform/state/terraform.tfstate -auto-approve &&
	echo "Terraform %s completed"`, action, action, action, action)

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
					Containers: []corev1.Container{
						{
							Name:  "terraform",
							Image: "hashicorp/terraform:light",
							Command: []string{
								"sh",
								"-c",
								command,
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "terraform-config",
									MountPath: "/app/terraform",
								},
								{
									Name:      "terraform-state",
									MountPath: "/app/terraform/state",
								},
							},
							WorkingDir: "/app/terraform",
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
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "terraform-pvc",
								},
							},
						},
					},
				},
			},
		},
	}

	// Create the Job in Kubernetes.
	_, err := e.k8sClientSet.BatchV1().Jobs(options.Namespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes Job: %v", err)
	}

	return nil
}

// WaitForJobCompletion waits for the Kubernetes Job to complete. - Add Logger
func (e *executor) waitForJobCompletion(ctx context.Context, jobName, namespace string) error {
	watchInterface, err := e.k8sClientSet.BatchV1().Jobs(namespace).Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", jobName),
	})
	if err != nil {
		return fmt.Errorf("failed to start watch on Job: %v", err)
	}
	defer watchInterface.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled or timed out")
		case event, ok := <-watchInterface.ResultChan():
			if !ok {
				return fmt.Errorf("watch channel closed")
			}
			job, ok := event.Object.(*batchv1.Job)
			if !ok {
				return fmt.Errorf("unexpected type received from watch")
			}

			if job.Status.Succeeded > 0 {
				fmt.Println("Job completed successfully")
				return nil
			}

			if job.Status.Failed > 0 {
				return fmt.Errorf("Job failed")
			}
		}
	}
}

func (e *executor) cleanupJobAndConfigMap(ctx context.Context, jobName, namespace string) error {
	// Delete the Job
	err := e.k8sClientSet.BatchV1().Jobs(namespace).Delete(ctx, jobName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete Job: %v", err)
	}
	/*
		// Delete the ConfigMap
		err = e.k8sClientSet.CoreV1().ConfigMaps(namespace).Delete(ctx, configMapName, metav1.DeleteOptions{})
		if err != nil {
			return fmt.Errorf("failed to delete ConfigMap: %v", err)
		}
	*/
	return nil
}

func (e *executor) createPersistentVolume(ctx context.Context) error {
	pv := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: "terraform-pv",
		},
		Spec: corev1.PersistentVolumeSpec{
			Capacity: corev1.ResourceList{
				corev1.ResourceStorage: resource.MustParse("1Gi"),
			},
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/mnt/data/terraform",
				},
			},
		},
	}

	_, err := e.k8sClientSet.CoreV1().PersistentVolumes().Create(ctx, pv, metav1.CreateOptions{})
	return err
}

func (e *executor) createPersistentVolumeClaim(ctx context.Context, namespace string) error {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "terraform-pvc",
			Namespace: namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
		},
	}

	_, err := e.k8sClientSet.CoreV1().PersistentVolumeClaims(namespace).Create(ctx, pvc, metav1.CreateOptions{})
	return err
}

func (e *executor) retrieveTerraformOutputs(ctx context.Context, options Options) (*tfjson.State, error) {
	// Define the path to the state file
	stateFilePath := "/app/terraform/state/terraform.tfstate"

	// Read the state file
	stateFile, err := os.ReadFile(stateFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %v", err)
	}

	// Parse the state file
	var state tfjson.State
	if err := json.Unmarshal(stateFile, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %v", err)
	}

	return &state, nil
}

func (e *executor) createConfigPVC(ctx context.Context, namespace string) error {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "terraform-config-pvc",
			Namespace: namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany, // RWX to allow multiple pods to mount if needed
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
		},
	}

	_, err := e.k8sClientSet.CoreV1().PersistentVolumeClaims(namespace).Create(ctx, pvc, metav1.CreateOptions{})
	return err
}

/*
func NewExecutor(clientset *kubernetes.Clientset) *executor {
    return &executor{
        clientset: clientset,
        // Initialize other fields if necessary.
    }
}

// driver
type TerraformDriver struct {
    clientset *kubernetes.Clientset
}

func NewTerraformDriver(clientset *kubernetes.Clientset) *TerraformDriver {
    return &TerraformDriver{
        clientset: clientset,
    }
}

func (d *TerraformDriver) Execute(ctx context.Context, cfg recipes.RecipeConfig) error {
    exec := terraform.NewExecutor(d.clientset)

    options := terraform.Options{
        RootDir: cfg.RootDir,
        RecipeName: cfg.RecipeName,
        Namespace:  cfg.Namespace, // Assuming Namespace is part of RecipeConfig
        // Set other options as necessary.
    }

    return exec.Deploy(ctx, options)
}

func (d *TerraformDriver) Destroy(ctx context.Context, cfg recipes.RecipeConfig) error {
    exec := terraform.NewExecutor(d.clientset)

    options := terraform.Options{
        RootDir: cfg.RootDir,
        RecipeName: cfg.RecipeName,
        Namespace:  cfg.Namespace,
        // Set other options as necessary.
    }

    return exec.Delete(ctx, options)
}

// engine
func RunRecipe(ctx context.Context, cfg recipes.RecipeConfig) error {
    // Initialize the Kubernetes client.
    kubeClient, err := kubernetes.NewClient()
    if err != nil {
        return err
    }

    // Initialize the TerraformDriver.
    terraformDriver := driver.NewTerraformDriver(kubeClient.Clientset)

    // Execute the recipe.
    return terraformDriver.Execute(ctx, cfg)
}


// deleteStateFile deletes the Terraform state file.
func (e *executor) deleteStateFile(ctx context.Context) error {
    // Define the path to the state file
    stateFilePath := "/app/terraform/state/terraform.tfstate"

    // Check if the state file exists
    if _, err := os.Stat(stateFilePath); os.IsNotExist(err) {
        return fmt.Errorf("state file does not exist at path: %s", stateFilePath)
    }

    // Delete the state file
    if err := os.Remove(stateFilePath); err != nil {
        return fmt.Errorf("failed to delete state file: %v", err)
    }

    fmt.Printf("State file deleted successfully: %s\n", stateFilePath)
    return nil
}


func (e *executor) createOrReplaceConfigMap(ctx context.Context, configMapName string, options Options) error {
    configMapClient := e.k8sClientSet.CoreV1().ConfigMaps(options.Namespace)

    // Define the ConfigMap data
    configMapData := map[string]string{
        // Add your configuration data here
        "main.tf": "content of your main.tf file",
    }

    // Create the ConfigMap object
    configMap := &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name:      configMapName,
            Namespace: options.Namespace,
        },
        Data: configMapData,
    }

    // Delete the existing ConfigMap if it exists
    err := configMapClient.Delete(ctx, configMapName, metav1.DeleteOptions{})
    if err != nil && !errors.IsNotFound(err) {
        return fmt.Errorf("failed to delete existing ConfigMap: %v", err)
    }

    // Create the new ConfigMap
    _, err = configMapClient.Create(ctx, configMap, metav1.CreateOptions{})
    if err != nil {
        return fmt.Errorf("failed to create ConfigMap: %v", err)
    }

    return nil
}

*/
