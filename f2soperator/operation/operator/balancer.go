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
	removeSurplusItems()

	// check which deployments are missing in k8s namespace f2s-containers
	addMissingDeployments()

	// scale deployments
	scaleDeployments()
}

// check which deployments are missing in k8s namespace f2s-containers
func addMissingDeployments() {
	logging.Debug("checking for k8s f2s-containers missing deployments")

	// get list of function definitions
	functions := configuration.ActiveConfiguration.Functions
	logging.Debug(fmt.Sprintf("found %d function definitions in active configuration", len(functions.Items)))

	// get current deployments from kubernetes (namespace 'f2s-containers')
	deployments, err := kubernetesservice.GetDeployments()
	if err != nil {
		logging.Error(fmt.Errorf("could not get list of deployments in kubernetes namespace 'f2s-containers' because: %s", err.Error()))
	}
	logging.Debug(fmt.Sprintf("found %d deployments in kubernetes namespace f2s-containers", len(deployments.Items)))

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
		} else {
			logging.Debug(fmt.Sprintf("deployment for function '%s' is already existing in kubernetes namespace 'f2s-containers'", f.Name))
		}
	}
}

// check which deployments in k8s namespace f2s-containers have no corresponding f2sfunction
func removeSurplusItems() {
	logging.Debug("checking for k8s f2s-containers surplus deployments")

	// get list of defined functions in active configuration
	functions := configuration.ActiveConfiguration.Functions
	logging.Debug(fmt.Sprintf("found %d function definitions in current configuration", len(functions.Items)))

	// get deployments from kubernetes namespace 'f2s-containers'
	deployments, err := kubernetesservice.GetDeployments()
	if err != nil {
		logging.Error(fmt.Errorf("could not get list of deployments in kubernetes namespace 'f2s-containers' because: %s", err.Error()))
	}
	logging.Debug(fmt.Sprintf("found %d deployments in kubernetes namespace 'f2s-containers'", len(deployments.Items)))

	services, err := kubernetesservice.ListServices()
	if err != nil {
		log.Fatal(err)
	}

	// check for surplus deployments
	for _, d := range deployments.Items {
		// check if deployment can be found in functions
		functionExisting := stringArrayContains(d.Name, functions.GetNames())

		if !functionExisting {
			logging.Info(fmt.Sprintf("delete surplus deployment %s (%s)", d.Name, d.UID))
			kubernetesservice.DeleteDeployment(string(d.UID))
		} else {
			logging.Debug(fmt.Sprintf("deployment '%s' matches a function definition. => not surplus", d.Name))
		}
	}

	// check for surplus services
	logging.Debug("start checking for surplus kubernetes service definitions")
	for _, s := range services.Items {
		// check if deployment can be found in functions
		functionExisting := stringArrayContains(s.Name, functions.GetNames())

		if !functionExisting {
			logging.Info(fmt.Sprintf("delete surplus service %s (%s)", s.Name, s.UID))
			kubernetesservice.DeleteService(string(s.UID))
		} else {
			logging.Debug(fmt.Sprintf("service '%' matches a function definition. => not surplus", s.Name))
		}
	}
}

// convert a millisecond timestamp (string) to time.Time
func convertMillisToTime(timestampMillis string) (time.Time, error) {
	timestampMillisInt, err := strconv.ParseInt(timestampMillis, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	timestampSeconds := timestampMillisInt / 1000
	timestamp := time.Unix(timestampSeconds, 0)

	return timestamp, nil
}

// get timestamp of last scaling action
// the timestamp is saved in the deployment (namespace f2s-containers) as annotation 'f2s/last-scaling'
func getLastScalingTimestamp(deploymentName string) (time.Time, bool) {
	logging.Debug(fmt.Sprintf("get timestamp of last scaling action (annotation 'f2s/last-scaling') for deployment '%s'", deploymentName))

	// get annotations of k8s deployment
	logging.Debug(fmt.Sprintf("get annotations of kubernetes deployment '%s'", deploymentName))
	annotations, err := kubernetesservice.GetDeploymentAnnotations(deploymentName)
	if err != nil {
		logging.Error(fmt.Errorf("could not get annotations of kubernetes deployment '%s': %s", deploymentName, err.Error()))
		return time.Time{}, false
	}

	// find annotation 'f2s/last-scaling'
	timestamp, exists := annotations["f2s/last-scaling"]

	if exists {
		// convert to time.time
		tm, err := convertMillisToTime(timestamp)
		if err != nil {
			logging.Error(fmt.Errorf("could not convert timestamp %s to time.time: %s", timestamp, err.Error()))
			return time.Time{}, false
		}
		return tm, true
	} else {
		logging.Warn(fmt.Sprintf("could not get last scaling timestamp of deployment '%s'. it does not seem to have the annotation 'f2s/last-scaling'", deploymentName))
		return time.Time{}, false
	}
}

// scale each function deployment according to prometheus metric
// f2sscaling_function_scaling_decision_replicas_difference
func scaleDeployments() {
	logging.Debug("start scaling of deployments")

	// get function definitions from active configuration
	functions := configuration.ActiveConfiguration.Functions
	logging.Debug(fmt.Sprintf("found %d function definitions in active configuration", len(functions.Items)))

	for _, function := range functions.Items {
		logging.Debug(fmt.Sprintf("calculate scaling for function: %s", function.Name))

		var resultScale int

		// read current prometheus metrics necessary for scaling decision
		metricCurrentAvailableReplicas := fmt.Sprintf("kube_deployment_status_replicas_available{functionname=\"%s\"}", function.Name)
		metricRequiredContainers := fmt.Sprintf("job:function_containers_required:containers{functionname=\"%s\"} or vector(0)", function.Name)
		currentAvailableReplicas, availableReplicasErr := prometheus.ReadCurrentPrometheusMetricValue(&configuration.ActiveConfiguration, metricCurrentAvailableReplicas)
		requiredContainers, requiredContainersErr := prometheus.ReadCurrentPrometheusMetricValue(&configuration.ActiveConfiguration, metricRequiredContainers)

		if availableReplicasErr != nil {
			logging.Error(fmt.Errorf("there was an error when trying to read metric (%s). setting result-scale to 0: %s", metricCurrentAvailableReplicas, availableReplicasErr.Error()))
			resultScale = 0
		} else if requiredContainersErr != nil {
			logging.Error(fmt.Errorf("there was an error when trying to read metric (%s). setting result-scale to 0: %s", metricRequiredContainers, requiredContainersErr.Error()))
			resultScale = 0
		} else {
			resultScale = int(math.Ceil(requiredContainers))
		}
		logging.Info(fmt.Sprintf("function %s has currently %v replicas available", function.Name, currentAvailableReplicas))
		logging.Info(fmt.Sprintf("function %s has desired %v replicas. result scale: %d", function.Name, requiredContainers, resultScale))

		// get current inflight requests of function
		target, err := f2shub.F2SDispatcherHub.GetDispatcherFunctionByName(function.Name)
		if err != nil {
			logging.Error(fmt.Errorf("could not GetDispatcherFunctionByName: %s", err.Error()))
			logging.Error(fmt.Errorf("could not get function target for function-name: %s. skipping scaling of this function", function.Name))
			continue
		}
		numInflightRequests := target.GetTotalInflightRequests()
		logging.Debug(fmt.Sprintf("function '%s' has %d inflight requests", function.Name, numInflightRequests))
		if numInflightRequests > 0 && resultScale == 0 {
			logging.Info(fmt.Sprintf("dont scale function %s to zero because there are %d inflight requests. scale to 1", function.Name, numInflightRequests))
			resultScale = 1
		}

		// check if last scaling of the function was less than 15 seconds ago
		fifteenSecondsAgo := time.Now().Add(-15 * time.Second)
		lastScalingTimestamp, exists := getLastScalingTimestamp(function.Name)
		if exists && lastScalingTimestamp.After(fifteenSecondsAgo) {
			logging.Info(fmt.Sprintf("last scaling of function %s was less than 15 seconds ago (%s). skipping regular scaling for this iteration", function.Name, lastScalingTimestamp.String()))
			continue
		}

		// check minumums and maximums
		if resultScale > function.Target.MaxReplicas {
			logging.Debug(fmt.Sprintf("result scale for function '%s' was %d. but setting it insted to the functions defined MaxReplicas: %d", function.Name, resultScale, function.Target.MaxReplicas))
			resultScale = function.Target.MaxReplicas
		}
		if resultScale < function.Target.MinReplicas {
			logging.Debug(fmt.Sprintf("result scale for function '%s' was %d. but setting it instead to the functions defined MinReplicas: %d", function.Name, resultScale, function.Target.MinReplicas))
			resultScale = function.Target.MinReplicas
		}

		// do the scaling
		if currentAvailableReplicas != float64(resultScale) {
			logging.Info(fmt.Sprintf("scaling %s: before=>%v after=>%v. scaling it now on kubernetes", function.Name, currentAvailableReplicas, resultScale))
			kubernetesservice.ScaleDeployment(function.Name, int32(resultScale))

			// set last-scaling-time in dispatcher state
			logging.Debug(fmt.Sprintf("function '%s' was just scaled. saving last scaling timestamp annotation 'f2s/last-scaling' in kubernetes deployment", function.Name))
			dispatcherFunction, err := f2shub.F2SDispatcherHub.GetDispatcherFunctionByName(function.Name)
			if err != nil {
				logging.Warn(fmt.Sprintf("could not get functiontarget for function: %s. %s", function.Name, err.Error()))
			}
			dispatcherFunction.SetLastScaling()

			// annotate deployment with last-scaling timestamp
			timestampMillis := time.Now().UnixNano() / int64(time.Millisecond)
			kubernetesservice.AnnotateDeployment(function.Name, map[string]string{
				"f2s/last-scaling": fmt.Sprintf("%v", timestampMillis),
			})
			kubernetesservice.AnnotateFunction(function.Name, map[string]string{
				"f2s/last-scaling": fmt.Sprintf("%v", timestampMillis),
			})

			// send 'function scaled' event
			f2shub.F2SEventManager.Publish(eventmanager.Event{
				UID:         f2shub.F2SEventManager.GenerateUUID(),
				Data:        function.Prettify(),
				Type:        eventmanager.Event_FunctionScaled,
				Description: fmt.Sprintf("F2SFunction %s scaled from %v to %v", function.Name, currentAvailableReplicas, int(resultScale)),
			})
		} else {
			logging.Debug(fmt.Sprintf("scale of function '%s' already is set to %d. scaling not necessary", function.Name, resultScale))
		}
	}
}
