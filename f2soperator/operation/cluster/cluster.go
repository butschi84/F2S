package cluster

import (
	"butschi84/f2s/hub"
	"butschi84/f2s/services/logger"
	clusterstate "butschi84/f2s/state/cluster"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/memberlist"
)

var logging *logger.F2SLogger
var f2shub *hub.F2SHub

type customLogger struct {
	logger *logger.F2SLogger
}

func (cl customLogger) Output(calldepth int, s string) error {
	logging.Info(s)
	return nil
}

func (cl customLogger) Write(p []byte) (n int, err error) {
	logging.Info(string(p))
	return len(p), nil
}

func Initialize(h *hub.F2SHub) {

	// consume variables
	f2shub = h

	// initialize logging
	logging = logger.Initialize("dispatcher")

	// subscribe to configuration changes
	logging.Info("subscribing to config package events")
	f2shub.F2SEventManager.Subscribe(handleEvents)

	// join cluster memberlist
	joinMemberList()

	for {
		time.Sleep(time.Second)
	}
}

func joinMemberList() {
	// Get the hostname of the current machine
	logging.Info("resolving hostname in order to join cluster memberlist")
	hostname, err := os.Hostname()
	if err != nil {
		logging.Error(fmt.Errorf("failed to resolve hostname because: %s", err.Error()))
	}

	// Configure memberlist
	config := memberlist.DefaultLANConfig()
	config.BindPort = f2shub.F2SConfiguration.Config.F2S.Memberlist.BindPort
	config.Name = hostname
	config.LogOutput = &customLogger{logger: logging}
	list, err := memberlist.Create(config)
	if err != nil {
		logging.Error(fmt.Errorf("failed to join memberlist because: %s", err.Error()))
	}

	// Join an existing cluster
	f2shub.F2SClusterState.MemberlistAddress = fmt.Sprintf("%s:%v",
		f2shub.F2SConfiguration.Config.F2S.Memberlist.Cluster,
		f2shub.F2SConfiguration.Config.F2S.Memberlist.BindPort)

	// use environment var to join memberlist when it is set
	envMemberListAddress := os.Getenv("MEMBERLIST_ADDRESS")
	if envMemberListAddress != "" {
		f2shub.F2SClusterState.MemberlistAddress = envMemberListAddress
		logging.Info(fmt.Sprintf("environment variable MEMBERLIST_ADDRESS is set. using memberlist address: %s", f2shub.F2SClusterState.MemberlistAddress))
	}

	logging.Info(fmt.Sprintf("joining memberlist at: %s", f2shub.F2SClusterState.MemberlistAddress))
	_, err = list.Join([]string{f2shub.F2SClusterState.MemberlistAddress})
	if err != nil {
		log.Fatal("Failed to join cluster: ", err)
	}

	// Handle signals to gracefully leave the cluster
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		list.Leave(5 * time.Second)
		os.Exit(0)
	}()

	// Perform application logic
	// For example, loop to periodically print memberlist members
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		members := list.Members()
		f2shub.F2SClusterState.ClusterMembers = make([]clusterstate.F2SClusterStateClusterMember, len(members))
		for idx, member := range members {
			f2shub.F2SClusterState.ClusterMembers[idx] = clusterstate.F2SClusterStateClusterMember{
				Name:    member.Name,
				Address: member.Address(),
			}
		}
	}
}
