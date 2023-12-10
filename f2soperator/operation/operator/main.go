package operator

import (
	"butschi84/f2s/hub"
	"butschi84/f2s/services/logger"
	"butschi84/f2s/services/prometheus"
	"butschi84/f2s/state/configuration"
	"butschi84/f2s/state/operatorstate"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	kubernetesservice "butschi84/f2s/services/kubernetes"
)

var logging *logger.F2SLogger
var f2shub *hub.F2SHub

func init() {
	// initialize logging
	logging = logger.Initialize("operator")
}

func RunOperator(hub *hub.F2SHub) {
	f2shub = hub

	// initialize default (no master)
	f2shub.F2SOperatorState.IsMaster = false

	// subscribe to configuration changes
	logging.Info("subscribing to events")
	go f2shub.F2SEventManager.Subscribe(handleEvent)

	// watch change events of f2sfunction crd in k8s
	logging.Info("starting to watch f2sfunctions in k8s")
	go kubernetesservice.WatchKubernetesResource("functions.v1alpha1.f2s.opensight.ch", "f2s", OnF2SFunctionChanged)

	// // watch change events of endpoints in k8s (namespace f2s-containers)
	logging.Info("starting to watch endpoints in k8s")
	go kubernetesservice.WatchKubernetesResource("endpoints.v1.", "f2s-containers", OnF2SEndpointsChanged)

	go RebalancerLoop()

	for {
		time.Sleep(time.Second)
	}
}

func RebalancerLoop() {
	for {
		// check if this f2s replica is the master
		masterDecision, _ := CheckMaster()
		if masterDecision != f2shub.F2SOperatorState.IsMaster && f2shub.F2SOperatorState.IsMaster == true {
			logging.Info("this f2s pod is now master")
		}
		f2shub.F2SOperatorState.IsMaster = masterDecision

		// rebalance
		if f2shub.F2SOperatorState.IsMaster {
			// Perform the desired task
			logging.Info("rebalancing...")
			Rebalance()
		}

		// Sleep for 15 seconds
		time.Sleep(15 * time.Second)
	}
}

func CheckMaster() (bool, error) {
	logging.Debug("[check master] reading prometheus metric 'f2s_master_election_ready_pods'")
	result, err := prometheus.ReadPrometheusMetric(&configuration.ActiveConfiguration, "f2s_master_election_ready_pods")
	if err != nil {
		logging.Error(fmt.Errorf("[check master] prometheus seems not to be reachable. prometheus URL can also specified by 'export Prometheus_URL=localhost:9090'"))
	}

	// get all f2s replica uid's
	fs2Pods := make(map[string]string, len(result.Data.Result))

	// update state
	f2shub.F2SOperatorState.KnownOperators = make([]operatorstate.F2SKnownOperator, len(result.Data.Result))
	for i, _ := range result.Data.Result {
		fs2Pods[result.Data.Result[i].Metric["uid"]] = result.Data.Result[i].Metric["pod"]
		f2shub.F2SOperatorState.KnownOperators[i].PodName = result.Data.Result[i].Metric["pod"]
		f2shub.F2SOperatorState.KnownOperators[i].PodUID = result.Data.Result[i].Metric["uid"]

		// Extract the timestamp from the first result
		timestamp, _ := strconv.ParseFloat(result.Data.Result[i].Values[len(result.Data.Result[i].Values)-1][1].(string), 64)
		// timestamp, err := strconv.ParseFloat(timestampStr, 64)
		if err == nil {
			f2shub.F2SOperatorState.KnownOperators[i].LastContact = time.Unix(int64(timestamp), 0)
		}
	}

	// Extract the keys from the map
	uids := make([]string, 0, len(fs2Pods))
	for k := range fs2Pods {
		uids = append(uids, k)
	}

	// order alphabetically
	sort.Strings(uids)

	masterPod := fs2Pods[uids[0]]
	logging.Debug(fmt.Sprintf("[check master] master is: %s", masterPod))

	for i, knownOperator := range f2shub.F2SOperatorState.KnownOperators {
		if masterPod == knownOperator.PodName {
			f2shub.F2SOperatorState.KnownOperators[i].IsMaster = true
		}
	}

	hostname, err := os.Hostname()
	if hostname == masterPod {
		return true, nil
	}

	return false, nil
}
