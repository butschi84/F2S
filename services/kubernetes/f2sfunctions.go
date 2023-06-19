package kubernetesservice

import (
	typesV1alpha1 "butschi84/f2s/configuration/api/types/v1alpha1"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// create a new f2s function (crd in k8s namespace f2s)
func CreateF2SFunction(prettyFunction *typesV1alpha1.PrettyFunction) (*typesV1alpha1.Function, error) {
	logging.Println("CreateF2SFunction: request to create a new F2S Function (crd in k8s namespace 'f2s')")

	// initialize clientset
	logging.Println("initializing k8s clientset")
	clientSet, err := GetV1Alpha1ClientSet()
	if err != nil {
		logging.Println("error during clientset initialisation: ", err)
		panic(err)
	}

	// prepare metadata
	logging.Println("preparing metadata for new f2sfunction creation")
	newFunction := &typesV1alpha1.Function{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-deployment",
		},
		Spec:   typesV1alpha1.FunctionSpec{},
		Target: typesV1alpha1.FunctionTarget{},
	}

	// Create new function CRD Object
	logging.Println("creating function in k8s")
	function, err := clientSet.Functions("f2s").Create(newFunction)
	if err != nil {
		logging.Println("error during function creation: ", err)
		log.Fatal(err)
	}

	return function, err
}

// delete a f2s function (crd in k8s namespace f2s)
func DeleteF2SFunction(uid string) error {
	logging.Println("request to delete a F2SFunction in K8S: ", uid)

	// initialize clientset
	logging.Println("initializing k8s clientset")
	clientSet, err := GetV1Alpha1ClientSet()
	if err != nil {
		logging.Println("error during clientset initialisation: ", err)
		panic(err)
	}

	logging.Println("deleting f2sfunction in k8s")
	err = clientSet.Functions("f2s").Delete(uid, metav1.DeleteOptions{})

	if err != nil {
		logging.Println("error during deletion: ", err)
	}

	return err
}
