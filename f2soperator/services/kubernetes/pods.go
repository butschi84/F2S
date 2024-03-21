package kubernetesservice

import (
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPodAnnotations(podName string) (annotations map[string]string, err error) {
	// Get clientset
	clientset, err := GetV1ClientSet()
	if err != nil {
		return nil, err
	}

	// Get the deployment
	deployment, err := clientset.CoreV1().Pods("f2s-containers").Get(context.TODO(), podName, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return deployment.ObjectMeta.Annotations, nil
}

func AnnotatePod(podName string, annotations map[string]string) error {
	// get clientset
	clientset, err := GetV1ClientSet()
	if err != nil {
		log.Fatal(err)
	}

	// Get the deployment
	pod, err := clientset.CoreV1().Pods("f2s-containers").Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Add or update annotations
	if pod.ObjectMeta.Annotations == nil {
		pod.ObjectMeta.Annotations = make(map[string]string)
	}
	for key, value := range annotations {
		pod.ObjectMeta.Annotations[key] = value
	}

	// Update the pod with new annotations
	_, err = clientset.CoreV1().Pods("f2s-containers").Update(context.TODO(), pod, v1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
