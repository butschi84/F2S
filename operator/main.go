package operator

import (
	"butschi84/f2s/configuration"
	"butschi84/f2s/logger"
	kubernetesservice "butschi84/f2s/services/kubernetes"
	"fmt"
	"sync"
	"time"

	"log"
)

var logging *log.Logger

func init() {
	// initialize logging
	logging = logger.Initialize("operator")

}

func RunOperator(config *configuration.F2SConfiguration, wg *sync.WaitGroup) {
	defer wg.Done()

	// subscribe to configuration changes
	logging.Println("subscribing to config package events")
	config.EventManager.Subscribe(handleEvent)

	for {
		// Perform the desired task
		logging.Println("rebalancing...")

		Rebalance()

		// Sleep for 30 seconds
		time.Sleep(10 * time.Second)
	}
}

func stringArrayContains(target string, arr []string) bool {
	found := false
	for _, str := range arr {
		if str == target {
			found = true
			break
		}
	}
	return found
}

// manage k8s deployments in namespace f2s-containers
func Rebalance() {
	logging.Println("starting rebalance")

	// check for surplus deployments in f2s-containers namespace
	logging.Println("checking for k8s f2s-containers surplus deployments")
	removeSurplusDeployments()

	// check which deployments are missing in k8s namespace f2s-containers
	logging.Println("checking for k8s f2s-containers missing deployments")
	addMissingDeployments()
}

// check which deployments are missing in k8s namespace f2s-containers
func addMissingDeployments() {
	functions := configuration.ActiveConfiguration.Functions
	deployments, err := kubernetesservice.GetDeployments()
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range functions.Items {
		// check if deployment can be found
		deploymentExisting := false
		for _, d := range deployments.Items {
			if d.Name == f.Name {
				deploymentExisting = true
			}
		}
		if !deploymentExisting {
			logging.Println(fmt.Sprintf("deployment for function %s (%s) has to be created", f.Name, f.UID))
			kubernetesservice.CreateDeployment(f.Name, f.Target.ContainerImage)
		}
	}
}

// check which deployments in k8s namespace f2s-containers have no corresponding f2sfunction
func removeSurplusDeployments() {
	functions := configuration.ActiveConfiguration.Functions
	deployments, err := kubernetesservice.GetDeployments()
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range deployments.Items {
		// check if deployment can be found in functions
		functionExisting := stringArrayContains(d.Name, functions.GetNames())
		logging.Println(fmt.Sprintf("search result for deployment %s %v", d.Name, functionExisting))

		if !functionExisting {
			logging.Println(fmt.Sprintf("delete surplus deployment %s (%s)", d.Name, d.UID))
			kubernetesservice.DeleteDeployment(string(d.UID))
		}
	}
}
