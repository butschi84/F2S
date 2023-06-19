package kubernetesservice

import (
	"context"
	"errors"
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

func CreateDeployment(name string, image string, labels map[string]string) (*appsv1.Deployment, error) {
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
	logging.Println("request to delete a K8S deployment:", uid)

	// initialize clientset
	logging.Println("initializing k8s clientset")
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
			logging.Println("deleting f2sfunction in k8s")
			err = clientset.AppsV1().Deployments("f2s-containers").Delete(context.Background(), d.Name, metav1.DeleteOptions{})

			if err != nil {
				logging.Println("error during deletion: ", err)
			}

			return nil
		}
	}

	return errors.New("deployment with this uid not found")
}
