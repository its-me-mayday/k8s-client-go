package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"net/http"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createKubeClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			return nil, fmt.Errorf("Error in Kubernetes configuration: %v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Error in Kubernetes clientset %v", err)
	}

	return clientset, nil
}

func getNamespaces(w http.ResponseWriter, r *http.Request) {
	clientset, err := createKubeClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(namespaces.Items)
}

func main() {
	log.Println("k8s-client-go starts")
	http.HandleFunc("/namespaces", getNamespaces)

	log.Printf("Server listen to %d\n", 8085)
	log.Fatal(http.ListenAndServe(":8085", nil))

}
