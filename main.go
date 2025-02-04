package main

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	// Użyj kubeconfig, jeśli jesteś lokalnie
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = fmt.Sprintf("%s/.kube/config", home)
	}

	// Użyj configu z kubeconfig lub z klastra (in-cluster)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Nie udało się stworzyć konfiguracji: %v", err)
	}

	// Tworzenie klienta
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Nie udało się stworzyć klienta Kubernetes: %v", err)
	}

	// Przykład: pobieranie listy deploymentów
	deployments, err := clientset.AppsV1().Deployments("kube-system").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Nie udało się pobrać deploymentów: %v", err)
	}

	fmt.Println("Adnotacje dla deploymentów w namespace 'kube-system':")
	for _, deployment := range deployments.Items {
		fmt.Printf("Deployment: %s\n", deployment.Name)
		for key, value := range deployment.Annotations {
			if value == "1" {
				fmt.Println("klops") // Zastąp pustą wartością
			}
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	fmt.Println("Deploymenty w namespace 'default':")
	for _, deployment := range deployments.Items {
		fmt.Printf("- %s\n", deployment.Name)
	}
}
