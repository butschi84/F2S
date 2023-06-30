package kubernetesservice

import (
	"context"
	"fmt"
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

// search a k8s endpoint by name
func GetEndpointWithName(functionName string) (corev1.Endpoints, error) {
	// get all endpoints from k8s
	endpoints, _ := ListEndpoints()

	// search specific endpoint by name
	for _, endpoint := range endpoints.Items {
		if endpoint.Name == functionName {
			return endpoint, nil
		}
	}

	// endpoint was not found
	return corev1.Endpoints{}, fmt.Errorf("the specified endpoint does not exist")
}
