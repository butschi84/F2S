package kubernetesservice

import (
	"context"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListEndpoints() (*corev1.EndpointsList, error) {
	// Get clientset
	clientset, err := GetV1ClientSet()
	if err != nil {
		log.Fatal(err)
	}

	// Get endpoints in the specified namespace
	endpoints, err := clientset.CoreV1().Endpoints("f2s-containers").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	return endpoints, err
}
