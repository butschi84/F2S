package clusterstate

type F2SClusterState struct {
	MemberlistAddress string                         `json:"memberlist_address"`
	ClusterMembers    []F2SClusterStateClusterMember `json:"cluster_members`
}

type F2SClusterStateClusterMember struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}
