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
    // ... existing imports
    "bytes"
    "fmt"
    "k8s.io/client-go/tools/remotecommand"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/client-go/kubernetes/scheme"
)

func execCommand(clientset *kubernetes.Clientset, restconfig *restclient.Config, namespace, podName, containerName string, command []string) (string, string, error) {
    req := clientset.CoreV1().RESTClient().Post().
        Resource("pods").
        Name(podName).
        Namespace(namespace).
        SubResource("exec").
        VersionedParams(&corev1.PodExecOptions{
            Container: containerName,
            Command:   command,
            Stdout:    true,
            Stderr:    true,
            Stdin:     false,
            TTY:       false,
        }, scheme.ParameterCodec)

    exec, err := remotecommand.NewSPDYExecutor(restconfig, "POST", req.URL())
    if err != nil {
        return "", "", err
    }

    var stdout, stderr bytes.Buffer
    err = exec.Stream(remotecommand.StreamOptions{
        Stdout: &stdout,
        Stderr: &stderr,
    })
    if err != nil {
        return "", "", err
    }

    return stdout.String(), stderr.String(), nil
}

import (
    // ... existing imports
    "archive/tar"
    "io"
    "os"
    "path/filepath"
)

func copyToPod(clientset *kubernetes.Clientset, restconfig *restclient.Config, namespace, podName, containerName, srcPath, destPath string) error {
    // Create tar archive of the source path
    tarReader, err := createTarStream(srcPath)
    if err != nil {
        return fmt.Errorf("failed to create tar stream: %v", err)
    }

    // Prepare the exec request
    cmd := []string{"tar", "xvf", "-", "-C", destPath}
    req := clientset.CoreV1().RESTClient().Post().
        Resource("pods").
        Name(podName).
        Namespace(namespace).
        SubResource("exec").
        VersionedParams(&corev1.PodExecOptions{
            Container: containerName,
            Command:   cmd,
            Stdin:     true,
            Stdout:    true,
            Stderr:    true,
        }, scheme.ParameterCodec)

    exec, err := remotecommand.NewSPDYExecutor(restconfig, "POST", req.URL())
    if err != nil {
        return fmt.Errorf("failed to initialize executor: %v", err)
    }

    // Stream the tar archive into the pod
    err = exec.Stream(remotecommand.StreamOptions{
        Stdin:  tarReader,
        Stdout: os.Stdout,
        Stderr: os.Stderr,
    })
    if err != nil {
        return fmt.Errorf("failed to execute command: %v", err)
    }

    return nil
}

func createTarStream(srcPath string) (io.Reader, error) {
    pipeReader, pipeWriter := io.Pipe()
    tarWriter := tar.NewWriter(pipeWriter)

    go func() {
        defer pipeWriter.Close()
        defer tarWriter.Close()

        err := filepath.Walk(srcPath, func(file string, fi os.FileInfo, err error) error {
            if err != nil {
                return err
            }

            relPath, err := filepath.Rel(srcPath, file)
            if err != nil {
                return err
            }

            header, err := tar.FileInfoHeader(fi, fi.Name())
            if err != nil {
                return err
            }

            header.Name = relPath

            if err := tarWriter.WriteHeader(header); err != nil {
                return err
            }

            if fi.Mode().IsRegular() {
                f, err := os.Open(file)
                if err != nil {
                    return err
                }
                defer f.Close()
                if _, err := io.Copy(tarWriter, f); err != nil {
                    return err
                }
            }

            return nil
        })

        if err != nil {
            pipeWriter.CloseWithError(err)
        }
    }()

    return pipeReader, nil
}

import (
    // ... existing imports
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createDataLoaderPod(clientset *kubernetes.Clientset, namespace string) error {
    pod := &corev1.Pod{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "data-loader-pod",
            Namespace: namespace,
        },
        Spec: corev1.PodSpec{
            Containers: []corev1.Container{
                {
                    Name:    "data-loader",
                    Image:   "alpine", // Or any lightweight image
                    Command: []string{"sleep", "3600"},
                    VolumeMounts: []corev1.VolumeMount{
                        {
                            Name:      "terraform-config",
                            MountPath: "/mnt/terraform-config",
                        },
                    },
                },
            },
            Volumes: []corev1.Volume{
                {
                    Name: "terraform-config",
                    PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
                        ClaimName: "terraform-config-pvc",
                    },
                },
            },
            RestartPolicy: corev1.RestartPolicyNever,
        },
    }

    _, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
    if err != nil {
        if !apierrors.IsAlreadyExists(err) {
            return fmt.Errorf("failed to create data-loader-pod: %v", err)
        }
    }

    return nil
}

import (
    // ... existing imports
    "time"
    apierrors "k8s.io/apimachinery/pkg/api/errors"
)

func waitForPodRunning(clientset *kubernetes.Clientset, namespace, podName string) error {
    timeout := time.After(60 * time.Second)
    tick := time.Tick(2 * time.Second)

    for {
        select {
        case <-timeout:
            return fmt.Errorf("timeout waiting for pod %s to be running", podName)
        case <-tick:
            pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
            if err != nil {
                if apierrors.IsNotFound(err) {
                    continue
                }
                return err
            }
            if pod.Status.Phase == corev1.PodRunning {
                return nil
            }
        }
    }
}
