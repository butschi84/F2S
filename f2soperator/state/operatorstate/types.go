package operatorstate

import "time"

type F2SKnownOperator struct {
	PodName     string    `json:"pod_name"`
	PodUID      string    `json:"pod_uid"`
	IsMaster    bool      `json:"is_master"`
	LastContact time.Time `json:"last_contact"`
}

type F2SOperatorState struct {
	IsMaster       bool               `json:"-"`
	KnownOperators []F2SKnownOperator `json:"known_operators"`
}
