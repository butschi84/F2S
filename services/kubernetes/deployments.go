package kubernetesservice

import (
	"context"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetDeployments() (*appsv1.DeploymentList, error) {

	// get clientset
	clientset, err := GetV1ClientSet()
	if err != nil {
		log.Fatal(err)
	}

	// Get all deployments in the cluster
	deployments, err := clientset.AppsV1().Deployments("f2s-containers").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	return deployments, err
}
