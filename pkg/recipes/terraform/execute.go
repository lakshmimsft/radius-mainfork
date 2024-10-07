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
    batchv1 "k8s.io/api/batch/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

job := &batchv1.Job{
    ObjectMeta: metav1.ObjectMeta{
        Name: "my-job",
    },
    Spec: batchv1.JobSpec{
        Template: corev1.PodTemplateSpec{
            Spec: corev1.PodSpec{
                Containers: []corev1.Container{
                    {
                        Name:    "my-container",
                        Image:   "my-image",
                        Command: []string{"sh", "-c", "run your deploy command here"},
                    },
                },
                RestartPolicy: corev1.RestartPolicyNever,
            },
        },
        BackoffLimit: int32Ptr(4),
    },
}

func Deploy() error {
    config, err := rest.InClusterConfig()
    if err != nil {
        return err
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return err
    }

    job, err := clientset.BatchV1().Jobs("your-namespace").Create(context.TODO(), job, metav1.CreateOptions{})
    if err != nil {
        return err
    }

    return nil
}

func WaitForJobCompletion(clientset *kubernetes.Clientset, jobName string, namespace string) error {
    for {
        job, err := clientset.BatchV1().Jobs(namespace).Get(context.TODO(), jobName, metav1.GetOptions{})
        if err != nil {
            return err
        }

        if job.Status.Succeeded > 0 {
            fmt.Println("Job completed successfully")
            break
        }

        if job.Status.Failed > 0 {
            return fmt.Errorf("Job failed")
        }

        time.Sleep(5 * time.Second)
    }
    return nil
}
