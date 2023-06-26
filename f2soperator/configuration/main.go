package configuration

import (
	"butschi84/f2s/logger"
	"butschi84/f2s/services/eventmanager"
	kubernetesservice "butschi84/f2s/services/kubernetes"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var logging *log.Logger

var ActiveConfiguration F2SConfiguration

func init() {
	// initialize logging
	logging = logger.Initialize("configuration")

	logging.Println("getting all f2s functions using kubernetes service")
	functions, err := kubernetesservice.GetF2SFunctions()
	if err != nil {
		logging.Println("Failed to read f2s functions")
		return
	}

	// start eventmanager (other packages will subscribe to this)
	eventManager := eventmanager.NewEventManager()
	eventManager.Start()

	// read f2sconfigmap
	// Read YAML file
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	var config F2SConfigMap
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}

	logging.Println("initializing config")
	ActiveConfiguration = F2SConfiguration{
		Functions:    functions,
		EventManager: eventManager,
		Config:       config,
	}

	// debug output configmap
	if ActiveConfiguration.Config.Debug {
		logging.Println("f2s configuration:")
		logging.Println(fmt.Sprintf("=> debug: %v", true))
	}

	// watch change events of f2sfunction crd in k8s
	logging.Println("starting to watch f2sfunctions in k8s")
	go kubernetesservice.WatchKubernetesResource("functions.v1alpha1.f2s.opensight.ch", "f2s", OnF2SFunctionChanged)

	// watch change events of endpoints in k8s (namespace f2s-containers)
	logging.Println("starting to watch endpoints in k8s")
	go kubernetesservice.WatchKubernetesResource("endpoints.v1.", "f2s-containers", OnF2SEndpointsChanged)
}
