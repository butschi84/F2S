package kubernetesservice

import (
	"context"
	"errors"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func CreateService(name string, port int, labels map[string]string) (*corev1.Service, error) {
	// get clientset
	clientset, err := GetV1ClientSet()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new deployment
	newService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Port:       int32(port),
					TargetPort: intstr.FromInt(port),
				},
			},
		},
	}

	deployment, err := clientset.CoreV1().Services("f2s-containers").Create(context.Background(), newService, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}

	return deployment, err
}

func ListServices() (*corev1.ServiceList, error) {

	// get clientset
	clientset, err := GetV1ClientSet()
	if err != nil {
		log.Fatal(err)
	}

	// Get all deployments in the cluster
	services, err := clientset.CoreV1().Services("f2s-containers").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	return services, err
}

func DeleteService(uid string) error {
	// get clientset
	clientset, err := GetV1ClientSet()
	if err != nil {
		log.Fatal(err)
	}

	services, err := ListServices()
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range services.Items {
		if d.UID == types.UID(uid) {
			logging.Println("deleting f2sfunction in k8s")
			err = clientset.CoreV1().Services("f2s-containers").Delete(context.Background(), d.Name, metav1.DeleteOptions{})

			if err != nil {
				logging.Println("error during deletion: ", err)
			}

			return nil
		}
	}

	return errors.New("deployment with this uid not found")
}
