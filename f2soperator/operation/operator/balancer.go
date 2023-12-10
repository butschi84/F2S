package operator

import (
	kubernetesservice "butschi84/f2s/services/kubernetes"
	"butschi84/f2s/services/prometheus"
	"butschi84/f2s/state/configuration"
	"butschi84/f2s/state/eventmanager"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"
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

func convertMillisToTime(timestampMillis string) (time.Time, error) {
	timestampMillisInt, err := strconv.ParseInt(timestampMillis, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	timestampSeconds := timestampMillisInt / 1000
	timestamp := time.Unix(timestampSeconds, 0)
	return timestamp, nil
}

func getLastScalingTimestamp(deploymentName string) (time.Time, bool) {
	// get annotations of k8s deployment
	annotations, err := kubernetesservice.GetDeploymentAnnotations(deploymentName)
	if err != nil {
		logging.Error(fmt.Sprintf("could not get annotations of kubernetes deployment: %s", deploymentName))
		return time.Time{}, false
	}

	// find annotation 'f2s/last-scaling'
	timestamp, exists := annotations["f2s/last-scaling"]

	if exists {
		// convert to time.time
		tm, err := convertMillisToTime(timestamp)
		if err != nil {
			logging.Error(fmt.Sprintf("could not convert timestamp %s to time.time", timestamp))
			return time.Time{}, false
		}
		return tm, true
	} else {
		return time.Time{}, false
	}
}

// scale each function deployment according to prometheus metric
// f2sscaling_function_scaling_decision_replicas_difference
func scaleDeployments() {
	functions := configuration.ActiveConfiguration.Functions
	for _, function := range functions.Items {
		var resultScale int
		currentAvailableReplicas, availableReplicasErr := prometheus.ReadCurrentPrometheusMetricValue(&configuration.ActiveConfiguration, fmt.Sprintf("kube_deployment_status_replicas_available{functionname=\"%s\"}", function.Name))
		requiredContainers, requiredContainersErr := prometheus.ReadCurrentPrometheusMetricValue(&configuration.ActiveConfiguration, fmt.Sprintf("job:function_containers_required:containers{functionname=\"%s\"} or vector(0)", function.Name))
		if availableReplicasErr != nil {
			logging.Error("there was an error when trying to read metric [kube_deployment_status_replicas_available]. setting result-scale to 0")
			resultScale = 0
		} else if requiredContainersErr != nil {
			logging.Error("there was an error when trying to read metric [job:function_containers_required:containers]. setting result-scale to 0")
			resultScale = 0
		} else {
			resultScale = int(math.Ceil(requiredContainers))
		}
		logging.Info(fmt.Sprintf("function %s has currently %v replicas available", function.Name, currentAvailableReplicas))
		logging.Info(fmt.Sprintf("function %s has desired %v replicas", function.Name, requiredContainers))

		// get current inflight requests of function
		target, err := f2shub.F2SDispatcherHub.GetDispatcherFunctionByName(function.Name)
		if err != nil {
			logging.Error(fmt.Sprintf("%s", err))
			logging.Error(fmt.Sprintf("[scaling] could not get function target for function-name: %s. skipping scaling of this function...", function.Name))
			continue
		}
		numInflightRequests := target.GetTotalInflightRequests()
		if numInflightRequests > 0 && resultScale == 0 {
			logging.Info(fmt.Sprintf("[scaling] dont scale function %s to zero because there are %d inflight requests. scale to 1", function.Name, numInflightRequests))
			resultScale = 1
		}

		// check if last scaling of the function was less than 15 seconds ago
		fifteenSecondsAgo := time.Now().Add(-15 * time.Second)
		lastScalingTimestamp, exists := getLastScalingTimestamp(function.Name)
		if exists == true && lastScalingTimestamp.After(fifteenSecondsAgo) {
			logging.Info(fmt.Sprintf("[scaling] last scaling of function %s was less than 15 seconds ago. skipping regular scaling", function.Name))
			continue
		}

		// check minumums and maximums
		if resultScale > function.Target.MaxReplicas {
			resultScale = function.Target.MaxReplicas
		}
		if resultScale < function.Target.MinReplicas {
			resultScale = function.Target.MinReplicas
		}

		// do the scaling
		if currentAvailableReplicas != float64(resultScale) {
			logging.Info(fmt.Sprintf("[scaling] scaling %s: before=>%v after=>%v", function.Name, currentAvailableReplicas, resultScale))
			kubernetesservice.ScaleDeployment(function.Name, int32(resultScale))

			// set last-scaling-time in dispatcher state
			dispatcherFunction, err := f2shub.F2SDispatcherHub.GetDispatcherFunctionByName(function.Name)
			if err != nil {
				logging.Warn(fmt.Sprintf("[scaling] could not get functiontarget for function: %s. %s", &function.Name, err.Error()))
			}
			dispatcherFunction.SetLastScaling()

			// annotate deployment with last-scaling timestamp
			timestampMillis := time.Now().UnixNano() / int64(time.Millisecond)
			kubernetesservice.AnnotateDeployment(function.Name, map[string]string{
				"f2s/last-scaling": fmt.Sprintf("%v", timestampMillis),
			})

			// send 'function scaled' event
			f2shub.F2SEventManager.Publish(eventmanager.Event{
				UID:         f2shub.F2SEventManager.GenerateUUID(),
				Data:        function.Prettify(),
				Type:        eventmanager.Event_FunctionScaled,
				Description: fmt.Sprintf("F2SFunction %s scaled from %v to %v", function.Name, currentAvailableReplicas, int(resultScale)),
			})
		}
	}
}
