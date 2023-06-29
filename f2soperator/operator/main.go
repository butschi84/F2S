package operator

import (
	"butschi84/f2s/configuration"
	"butschi84/f2s/services/logger"
	"butschi84/f2s/services/prometheus"
	"os"
	"sort"
	"sync"
	"time"
)

var logging logger.F2SLogger
var master bool

func init() {
	// initialize logging
	logging = logger.Initialize("operator")
	master = false
}

func RunOperator(config *configuration.F2SConfiguration, wg *sync.WaitGroup) {
	defer wg.Done()

	// subscribe to configuration changes
	logging.Info("subscribing to events")
	config.EventManager.Subscribe(handleEvent)

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
