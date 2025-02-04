package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	// Define the workload flag
	workload := flag.String("workload", "all", "Specify the workload to fetch annotations for (deployments, daemonsets, cronjobs, all)")
	flag.Parse()

	// Use kubeconfig if available locally
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = fmt.Sprintf("%s/.kube/config", home)
	}

	// Use the config from kubeconfig or in-cluster config
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Failed to create configuration: %v", err)
	}

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	// Get the list of all namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get namespaces: %v", err)
	}

	for _, ns := range namespaces.Items {
		namespace := ns.Name
		fmt.Printf("Namespace: %s\n", namespace)

		switch *workload {
		case "deployments":
			printDeploymentAnnotations(clientset, namespace)
		case "daemonsets":
			printDaemonSetAnnotations(clientset, namespace)
		case "statefulset":
			printDaemonSetAnnotations(clientset, namespace)
		case "cronjobs":
			printCronJobAnnotations(clientset, namespace)
		case "all":
			printDeploymentAnnotations(clientset, namespace)
			printDaemonSetAnnotations(clientset, namespace)
			printCronJobAnnotations(clientset, namespace)
		default:
			log.Fatalf("Unknown workload type: %s", *workload)
		}
	}
}

func printDeploymentAnnotations(clientset *kubernetes.Clientset, namespace string) {
	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get Deployments: %v", err)
	}
	fmt.Println("Annotations for Deployments:")
	for _, deployment := range deployments.Items {
		printAnnotations("Deployment", deployment.Name, deployment.Annotations)
	}
}

func printDaemonSetAnnotations(clientset *kubernetes.Clientset, namespace string) {
	daemonsets, err := clientset.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get DaemonSets: %v", err)
	}
	fmt.Println("Annotations for DaemonSets:")
	for _, ds := range daemonsets.Items {
		printAnnotations("DaemonSet", ds.Name, ds.Annotations)
	}
}

func printStatefulSetAnnotations(clientset *kubernetes.Clientset, namespace string) {
	statefulset, err := clientset.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get StatefulSets: %v", err)
	}
	fmt.Println("Annotations for StatefulSets:")
	for _, sfs := range statefulset.Items {
		printAnnotations("StatefulSet", sfs.Name, sfs.Annotations)
	}
}

func printCronJobAnnotations(clientset *kubernetes.Clientset, namespace string) {
	cronjobs, err := clientset.BatchV1().CronJobs(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get CronJobs: %v", err)
	}
	fmt.Println("Annotations for CronJobs:")
	for _, cj := range cronjobs.Items {
		printAnnotations("CronJob", cj.Name, cj.Annotations)
	}
}

// Helper function to print annotations
func printAnnotations(resourceType, name string, annotations map[string]string) {
	fmt.Printf("%s: %s\n", resourceType, name)
	if len(annotations) == 0 {
		fmt.Println("  No annotations")
	}
	for key, value := range annotations {
		fmt.Printf("  %s: %s\n", key, value)
	}
}
