package configuration

import (
	clientV1alpha1 "butschi84/f2s/configuration/api/clientset/v1alpha1"
	typesV1alpha1 "butschi84/f2s/configuration/api/types/v1alpha1"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// get k8s cluster config
func getInClusterConfig() (*rest.Config, error) {
	var kubeconfig string
	var config *rest.Config
	var err error

	kubeconfig = "/Users/roman/.kube/config"
	if kubeconfig == "" {
		logging.Printf("using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		logging.Printf("using configuration from '%s'", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get in-cluster config: %v", err)
	}
	return config, nil
}

func GetCRDs() *typesV1alpha1.FunctionList {
	logging.Println("querying all crds")

	// Retrieve the in-cluster configuration
	config, err := getInClusterConfig()
	if err != nil {
		logging.Printf("Failed to get in-cluster config: %v\n", err)
		os.Exit(1)
	}

	clientSet, err := clientV1alpha1.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// ctx := context.TODO()
	functions, err := clientSet.Functions("f2s").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	logging.Printf("number of configured functions: %+v\n", len(functions.Items))
	return functions
}
