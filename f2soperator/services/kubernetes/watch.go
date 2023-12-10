package kubernetesservice

import (
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func WatchKubernetesResource(resource string, namespace string, callback func(obj interface{})) {
	//dynamic informer needs to be told which type to watch
	seldoninformer, _ := GetDynamicInformer(resource, namespace)
	stopper := make(chan struct{})
	defer close(stopper)
	runCRDInformer(stopper, seldoninformer.Informer(), callback)
}

// crd informer. watch changes to f2s functions in k8s namespace f2s
func runCRDInformer(stopCh <-chan struct{}, s cache.SharedIndexInformer, callback func(obj interface{})) {
	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			logging.Info("A new F2S Function has been added / configured")

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
			callback(obj)
		},
		DeleteFunc: func(obj interface{}) {
			logging.Info("A F2S Function has been removed")
			callback(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			callback(newObj)
		},
	}
	s.AddEventHandler(handlers)
	s.Run(stopCh)
}

func GetDynamicInformer(resource string, namespace string) (informers.GenericInformer, error) {
	cfg, err := getInClusterConfig()
	if err != nil {
		logging.Error(fmt.Errorf("[GetDynamicInformer] Failed to get in-cluster config: %s\n", err.Error()))
		os.Exit(1)
	}

	// Grab a dynamic interface that we can create informers from
	dc, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	// Create a factory object that can generate informers for resource types
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dc, 0, namespace, nil)
	gvr, _ := schema.ParseResourceArg(resource)
	// Finally, create our informer for deployments!
	informer := factory.ForResource(*gvr)
	return informer, nil
}
