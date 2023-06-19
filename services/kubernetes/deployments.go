package kubernetesservice

import (
	"context"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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

func CreateDeployment() (*appsv1.Deployment, error) {
	// get clientset
	clientset, err := GetV1ClientSet()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new deployment
	newDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "my-app",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "my-app",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "my-container",
							Image: "nginx",
						},
					},
				},
			},
		},
	}

	deployment, err := clientset.AppsV1().Deployments("f2s-containers").Create(context.Background(), newDeployment, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}

	return deployment, err
}
