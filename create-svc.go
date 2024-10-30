package main

import (
    "context"
    "flag"
    "fmt"
    "path/filepath"

    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/util/intstr"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"
)

func main() {
    fmt.Printf("Created Service %q.\n", nil )
    var kubeconfig *string
    if home := homedir.HomeDir(); home != "" {
        kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
    } else {
        kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
    }
    flag.Parse()

    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
        panic(err.Error())
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err.Error())
    }

    service := &corev1.Service{
        ObjectMeta: metav1.ObjectMeta{
            Name: "mar2-svc",
        },
        Spec: corev1.ServiceSpec{
            Selector: map[string]string{
                "app": "nginx",
            },
            Ports: []corev1.ServicePort{
                {
                    Protocol:   corev1.ProtocolTCP,
                    Port:       80,
                    TargetPort: intstr.FromInt(80),
                },
            },
        },
    }

    servicesClient := clientset.CoreV1().Services(corev1.NamespaceDefault)

    result, err := servicesClient.Create(context.TODO(), service, metav1.CreateOptions{})
    if err != nil {
        panic(err)
    }

    fmt.Printf("Created Service %q.\n", result.GetObjectMeta().GetName())
}
