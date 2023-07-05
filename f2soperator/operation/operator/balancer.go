package operator

import (
	kubernetesservice "butschi84/f2s/services/kubernetes"
	"butschi84/f2s/services/prometheus"
	"butschi84/f2s/state/configuration"
	"fmt"
	"log"
	"math"
)

// manage k8s deployments in namespace f2s-containers
func Rebalance() {
	// logging.Println("starting rebalance")

	// check for surplus deployments in f2s-containers namespace
	// logging.Println("checking for k8s f2s-containers surplus deployments")
	removeSurplusItems()

	// check which deployments are missing in k8s namespace f2s-containers
	// logging.Println("checking for k8s f2s-containers missing deployments")
	addMissingDeployments()

	scaleDeployments()
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
			logging.Info(fmt.Sprintf("deployment for function %s (%s) has to be created", f.Name, f.UID))
			kubernetesservice.CreateDeployment(f.Name, f.Target.ContainerImage, map[string]string{"f2sfunction": f.Name}, f.Target.Port)
			kubernetesservice.CreateService(f.Name, f.Target.Port, map[string]string{"f2sfunction": f.Name})
		}
	}
}

// check which deployments in k8s namespace f2s-containers have no corresponding f2sfunction
func removeSurplusItems() {
	functions := configuration.ActiveConfiguration.Functions
	deployments, err := kubernetesservice.GetDeployments()
	if err != nil {
		log.Fatal(err)
	}
	services, err := kubernetesservice.ListServices()
	if err != nil {
		log.Fatal(err)
	}

	// check for surplus deployments
	for _, d := range deployments.Items {
		// check if deployment can be found in functions
		functionExisting := stringArrayContains(d.Name, functions.GetNames())
		// logging.Println(fmt.Sprintf("search result for deployment %s %v", d.Name, functionExisting))

		if !functionExisting {
			logging.Info(fmt.Sprintf("delete surplus deployment %s (%s)", d.Name, d.UID))
			kubernetesservice.DeleteDeployment(string(d.UID))
		}
	}

	// check for surplus services
	for _, s := range services.Items {
		// check if deployment can be found in functions
		functionExisting := stringArrayContains(s.Name, functions.GetNames())
		// logging.Println(fmt.Sprintf("search result for service %s %v", s.Name, functionExisting))

		if !functionExisting {
			logging.Info(fmt.Sprintf("delete surplus service %s (%s)", s.Name, s.UID))
			kubernetesservice.DeleteService(string(s.UID))
		}
	}
}

// scale each function deployment according to prometheus metric
// f2sscaling_function_scaling_decision_replicas_difference
func scaleDeployments() {
	functions := configuration.ActiveConfiguration.Functions
	for _, function := range functions.Items {
		var resultScale int
		currentAvailableReplicas, err := prometheus.ReadPrometheusMetricValue(&configuration.ActiveConfiguration, "kube_deployment_status_replicas_available", map[string]string{"functionname": function.Name})
		requiredContainers, err := prometheus.ReadPrometheusMetricValue(&configuration.ActiveConfiguration, "job:function_containers_required:containers", map[string]string{"functionname": function.Name})
		if err != nil {
			// no invocations / metrics => scale to minimum
			resultScale = 0
		} else {
			resultScale = int(math.Ceil(requiredContainers))
		}

		// check minumums and maximums
		if resultScale > function.Target.MaxReplicas {
			resultScale = function.Target.MaxReplicas
		}
		if resultScale < function.Target.MinReplicas {
			resultScale = function.Target.MinReplicas
		}

		// do the scaling
		logging.Info(fmt.Sprintf("scaling function replicas %s from %v to %v", function.Name, currentAvailableReplicas, resultScale))
		// kubernetesservice.ScaleDeployment(function.Name, int32(resultScale))
	}
}
