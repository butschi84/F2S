package queue

import (
	"butschi84/f2s/state/configuration/api/types/v1alpha1"
)

type F2SAuthUser struct {
	Username string
	Group    string
}

// F2SRequest to invoke a function
type F2SRequest struct {
	UID           string
	Path          string
	Method        string
	Payload       string
	Function      v1alpha1.PrettyFunction
	ResultChannel chan F2SRequestResult

	F2SUser F2SAuthUser
}

type F2SRequestResult struct {
	UID         string `json:"uid"`
	Success     bool   `json:"success"`
	Result      []byte `json:"result"`
	Details     string `json:"details"`
	ContentType string `json:"contentType"`

	Request F2SRequest `json:"-"`

	Duration               float64 `json:"-"`
	F2SDispatcherTargetUID string  `json:"-"`
}
