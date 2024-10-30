package main

// It will take several minutes to complete. 

import (
    "context"
    "fmt"
    "sigs.k8s.io/controller-runtime/pkg/client"
    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
    // Create a new client
    cfg, err := config.GetConfig()
    if err != nil {
        fmt.Println("Error getting config:", err)
        return
    }

    k8sClient, err := client.New(cfg, client.Options{})
    if err != nil {
        fmt.Println("Error creating client:", err)
        return
    }

    // Define the labels
    labels := map[string]string{"app": "nginx"}

    // Define the Deployment
    deployment := &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "nginx-deployment",
            Namespace: "default",
        },
        Spec: appsv1.DeploymentSpec{
            Replicas: int32Ptr(3),
            Selector: &metav1.LabelSelector{
                MatchLabels: labels,
            },
            Template: corev1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: labels,
                },
                Spec: corev1.PodSpec{
                    Containers: []corev1.Container{
                        {
                            Name:  "nginx",
                            Image: "nginx:1.14.2",
                            Ports: []corev1.ContainerPort{
                                {
                                    ContainerPort: 80,
                                },
                            },
                        },
                    },
                },
            },
        },
    }

    // Create the Deployment
    ctx := context.Background()
    err = k8sClient.Create(ctx, deployment)
    if err != nil {
        fmt.Println("Error creating deployment:", err)
        return
    }

    fmt.Println("Deployment created successfully")
}

func int32Ptr(i int32) *int32 {
    return &i
}

