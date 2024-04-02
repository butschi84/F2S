package clusterstate

func Initialize() *F2SClusterState {
	return &F2SClusterState{
		MemberlistAddress: "",
	}
}
