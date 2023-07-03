package kubernetesservice

import (
	typesV1alpha1 "butschi84/f2s/state/configuration/api/types/v1alpha1"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// get all f2s functions
func GetF2SFunctions() (*typesV1alpha1.FunctionList, error) {
	logging.Info("GetF2SFunctions: request to get all F2SFunctions (crd's in k8s namespace 'f2s')")

	// initialize clientset
	logging.Info("initializing k8s clientset")
	clientSet, err := GetV1Alpha1ClientSet()
	if err != nil {
		logging.Error(fmt.Errorf("error during clientset initialisation: %s", err))
		panic(err)
	}

	functions, err := clientSet.Functions("f2s").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	logging.Info(fmt.Sprintf("number of configured functions: %d\n", len(functions.Items)))
	return functions, err
}

// create a new f2s function (crd in k8s namespace f2s)
func CreateF2SFunction(prettyFunction *typesV1alpha1.PrettyFunction) (*typesV1alpha1.Function, error) {
	logging.Info("CreateF2SFunction: request to create a new F2S Function (crd in k8s namespace 'f2s')")

	// initialize clientset
	logging.Info("initializing k8s clientset")
	clientSet, err := GetV1Alpha1ClientSet()
	if err != nil {
		logging.Error(fmt.Errorf("error during clientset initialisation: %s", err))
		panic(err)
	}

	// prepare metadata
	logging.Info("preparing metadata for new f2sfunction creation")
	newFunction := &typesV1alpha1.Function{
		ObjectMeta: metav1.ObjectMeta{
			Name: prettyFunction.Name,
		},
		Spec: typesV1alpha1.FunctionSpec{
			Endpoint: prettyFunction.Spec.Endpoint,
			Method:   prettyFunction.Spec.Method,
		},
		Target: typesV1alpha1.FunctionTarget{
			ContainerImage: prettyFunction.Target.ContainerImage,
			Endpoint:       prettyFunction.Target.Endpoint,
			Port:           prettyFunction.Target.Port,
			MinReplicas:    prettyFunction.Target.MinReplicas,
			MaxReplicas:    prettyFunction.Target.MaxReplicas,
		},
	}

	// Create new function CRD Object
	logging.Info("creating function in k8s")
	function, err := clientSet.Functions("f2s").Create(newFunction)
	if err != nil {
		logging.Error(fmt.Errorf("error during function creation: %s", err))
		log.Fatal(err)
	}

	return function, err
}

// delete a f2s function (crd in k8s namespace f2s)
func DeleteF2SFunction(uid string) error {
	logging.Info("request to delete a F2SFunction in K8S: ", uid)

	// initialize clientset
	logging.Info("initializing k8s clientset")
	clientSet, err := GetV1Alpha1ClientSet()
	if err != nil {
		logging.Error(fmt.Errorf("error during clientset initialisation: %s", err))
		panic(err)
	}

	logging.Info("deleting f2sfunction in k8s")
	err = clientSet.Functions("f2s").Delete(uid, metav1.DeleteOptions{})

	if err != nil {
		logging.Error(fmt.Errorf("error during deletion: %s", err))
	}

	return err
}
