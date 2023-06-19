package configuration

import (
	typesV1alpha1 "butschi84/f2s/configuration/api/types/v1alpha1"
	kubernetesservice "butschi84/f2s/services/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

func GetCRDs() *typesV1alpha1.FunctionList {
	logging.Println("querying all crds")

	metav1.AddMetaToScheme(scheme.Scheme)

	clientSet, err := kubernetesservice.GetV1Alpha1ClientSet()
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
