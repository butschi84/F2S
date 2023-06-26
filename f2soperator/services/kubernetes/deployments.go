package kubernetesservice

import (
	"context"
	"errors"
	"fmt"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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

func CreateDeployment(name string, image string, labels map[string]string, port int) (*appsv1.Deployment, error) {
	// get clientset
	clientset, err := GetV1ClientSet()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new deployment
	newDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
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
							Name:  name,
							Image: image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: int32(port),
								},
							},
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

func DeleteDeployment(uid string) error {
	logging.Info("request to delete a K8S deployment:", uid)

	// initialize clientset
	logging.Info("initializing k8s clientset")
	clientset, err := GetV1ClientSet()
	if err != nil {
		log.Fatal(err)
	}

	deployments, err := GetDeployments()
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range deployments.Items {
		if d.UID == types.UID(uid) {
			logging.Info("deleting f2sfunction in k8s")
			err = clientset.AppsV1().Deployments("f2s-containers").Delete(context.Background(), d.Name, metav1.DeleteOptions{})

			if err != nil {
				logging.Error(fmt.Errorf("error during deletion: %s", err))
			}

			return nil
		}
	}

	return errors.New("deployment with this uid not found")
}

func ScaleDeployment(deploymentName string, replicas int32) error {
	logging.Info("request to scale a K8S deployment:", deploymentName)

	// Initialize clientset
	logging.Info("initializing K8S clientset")
	clientset, err := GetV1ClientSet()
	if err != nil {
		log.Fatal(err)
	}

	deployment, err := clientset.AppsV1().Deployments("f2s-containers").Get(context.Background(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to retrieve deployment: %w", err)
	}

	// Update the replicas value
	deployment.Spec.Replicas = &replicas

	// Apply the scaling changes
	_, err = clientset.AppsV1().Deployments("f2s-containers").Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}
