package kubernetesservice

import (
	typesV1alpha1 "butschi84/f2s/configuration/api/types/v1alpha1"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateF2SFunction(prettyFunction *typesV1alpha1.PrettyFunction) (*typesV1alpha1.Function, error) {
	clientSet, err := GetV1Alpha1ClientSet()
	if err != nil {
		panic(err)
	}

	// Create a new deployment
	newDeployment := &typesV1alpha1.Function{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-deployment",
		},
		Spec:   typesV1alpha1.FunctionSpec{},
		Target: typesV1alpha1.FunctionTarget{},
	}

	function, err := clientSet.Functions("f2s").Create(newDeployment)
	if err != nil {
		log.Fatal(err)
	}

	return function, err
}
