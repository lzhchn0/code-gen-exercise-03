package main

import (
	"context"
	"fmt"
	"log"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Load kubernetes config
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		// Try in-cluster config if local config fails
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("Failed to create config: %v", err)
		}
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Define the ingress
	pathType := networkingv1.PathTypePrefix
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: "example-ingress",
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/",
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: "example.com",
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/api",
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: "api-service",
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Create the ingress
	namespace := "default"
	result, err := clientset.NetworkingV1().Ingresses(namespace).Create(
		context.Background(),
		ingress,
		metav1.CreateOptions{},
	)
	if err != nil {
		log.Fatalf("Failed to create ingress: %v", err)
	}
	fmt.Printf("Created ingress %q in namespace %q\n", result.GetName(), namespace)
} 


