package operator

import (
	"butschi84/f2s/hub"
	"butschi84/f2s/services/logger"
	"butschi84/f2s/services/prometheus"
	"butschi84/f2s/state/configuration"
	"os"
	"sort"
	"sync"
	"time"

	kubernetesservice "butschi84/f2s/services/kubernetes"
)

var logging logger.F2SLogger
var master bool
var f2shub *hub.F2SHub

func init() {
	// initialize logging
	logging = logger.Initialize("operator")
	master = false
}

func RunOperator(hub *hub.F2SHub, wg *sync.WaitGroup) {
	defer wg.Done()

	f2shub = hub

	// subscribe to configuration changes
	logging.Info("subscribing to events")
	hub.F2SEventManager.Subscribe(handleEvent)

	// watch change events of f2sfunction crd in k8s
	logging.Info("starting to watch f2sfunctions in k8s")
	go kubernetesservice.WatchKubernetesResource("functions.v1alpha1.f2s.opensight.ch", "f2s", OnF2SFunctionChanged)

	// // watch change events of endpoints in k8s (namespace f2s-containers)
	logging.Info("starting to watch endpoints in k8s")
	go kubernetesservice.WatchKubernetesResource("endpoints.v1.", "f2s-containers", OnF2SEndpointsChanged)

	for {
		// check if this f2s replica is the master
		masterDecision, _ := CheckMaster()
		if masterDecision != master {
			logging.Info("this f2s pod is now master")
			master = true
		}

		// rebalance
		if master {
			// Perform the desired task
			logging.Info("rebalancing...")
			Rebalance()
		}

		// Sleep for 15 seconds
		time.Sleep(15 * time.Second)
	}
}

func CheckMaster() (bool, error) {
	result, err := prometheus.ReadPrometheusMetric(&configuration.ActiveConfiguration, "f2s_master_election_ready_pods", map[string]string{})
	if err != nil {
		logging.Error(err)
	}

	// jsonBytes, err := json.MarshalIndent(result, "", "  ")
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return false, nil
	// }

	// get all f2s replica uid's
	fs2Pods := make(map[string]string, len(result.Data.Result))

	for i, _ := range result.Data.Result {
		fs2Pods[result.Data.Result[i].Metric["uid"]] = result.Data.Result[i].Metric["pod"]
	}

	// Extract the keys from the map
	uids := make([]string, 0, len(fs2Pods))
	for k := range fs2Pods {
		uids = append(uids, k)
	}

	// order alphabetically
	sort.Strings(uids)

	masterPod := fs2Pods[uids[0]]

	// Print the JSON string to stdout
	// fmt.Println(fmt.Sprintf("master pod name: %s", masterPod))

	hostname, err := os.Hostname()

	if hostname == masterPod {
		return true, nil
	}

	return false, nil
}
