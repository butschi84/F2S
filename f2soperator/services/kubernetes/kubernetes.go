package kubernetesservice

import (
	"fmt"
	"os"

	clientV1alpha1 "butschi84/f2s/configuration/api/clientset/v1alpha1"
	v1alpha1 "butschi84/f2s/configuration/api/types/v1alpha1"

	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// get k8s cluster config
func getInClusterConfig() (*rest.Config, error) {
	var kubeconfig string
	var config *rest.Config
	var err error

	kubeconfig = os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		logging.Info("using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		// logging.Printf("using configuration from '%s'", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get in-cluster config: %v", err)
	}
	return config, nil
}

func GetV1Alpha1ClientSet() (*clientV1alpha1.V1Alpha1Client, error) {
	// Retrieve the in-cluster configuration
	config, err := getInClusterConfig()
	if err != nil {
		logging.Error(fmt.Errorf("Failed to get in-cluster config: %s\n", err))
		os.Exit(1)
	}

	v1alpha1.AddToScheme(scheme.Scheme)

	clientSet, err := clientV1alpha1.NewForConfig(config, scheme.Scheme)
	if err != nil {
		panic(err)
	}

	return clientSet, nil
}

func GetV1ClientSet() (*k8s.Clientset, error) {
	// Retrieve the in-cluster configuration
	config, err := getInClusterConfig()
	if err != nil {
		logging.Error(fmt.Errorf("Failed to get in-cluster config: %s\n", err))
		os.Exit(1)
	}

	clientSet, err := k8s.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return clientSet, nil
}
