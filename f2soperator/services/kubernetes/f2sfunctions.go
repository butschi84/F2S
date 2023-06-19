package kubernetesservice

import (
	typesV1alpha1 "butschi84/f2s/configuration/api/types/v1alpha1"
	"log"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

// get all f2s functions
func GetF2SFunctions() (*typesV1alpha1.FunctionList, error) {
	logging.Println("GetF2SFunctions: request to get all F2SFunctions (crd's in k8s namespace 'f2s')")

	// initialize clientset
	logging.Println("initializing k8s clientset")
	clientSet, err := GetV1Alpha1ClientSet()
	if err != nil {
		logging.Println("error during clientset initialisation: ", err)
		panic(err)
	}

	functions, err := clientSet.Functions("f2s").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	logging.Printf("number of configured functions: %+v\n", len(functions.Items))
	return functions, err
}

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

func GetDynamicInformer() (informers.GenericInformer, error) {
	cfg, err := getInClusterConfig()
	if err != nil {
		logging.Printf("Failed to get in-cluster config: %v\n", err)
		os.Exit(1)
	}

	// Grab a dynamic interface that we can create informers from
	dc, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	// Create a factory object that can generate informers for resource types
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dc, 0, "f2s", nil)
	gvr, _ := schema.ParseResourceArg("functions.v1alpha1.f2s.opensight.ch")
	// Finally, create our informer for deployments!
	informer := factory.ForResource(*gvr)
	return informer, nil
}

func WatchF2SFunctions(callback func()) {
	//dynamic informer needs to be told which type to watch
	seldoninformer, _ := GetDynamicInformer()
	stopper := make(chan struct{})
	defer close(stopper)
	runCRDInformer(stopper, seldoninformer.Informer(), callback)
}

// crd informer. watch changes to f2s functions in k8s namespace f2s
func runCRDInformer(stopCh <-chan struct{}, s cache.SharedIndexInformer, callback func()) {
	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			logging.Println("A new F2S Function has been added / configured")

			// logging.Println("trying to parse event to F2SFunction Obj")
			// d := &typesV1alpha1.Function{}
			// err := runtime.DefaultUnstructuredConverter.
			// 	FromUnstructured(obj.(*unstructured.Unstructured).UnstructuredContent(), d)
			// if err != nil {
			// 	logging.Println("could not convert event to F2SFunction")
			// 	logging.Print(err)
			// 	return
			// }
			// logging.Println(fmt.Sprintf("added function %s (%s)", d.Name, d.UID))
			callback()
		},
		DeleteFunc: func(obj interface{}) {
			logging.Println("A F2S Function has been removed")
			callback()
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			logging.Println("A F2S Function has been updated")
			callback()
		},
	}
	s.AddEventHandler(handlers)
	s.Run(stopCh)
}
