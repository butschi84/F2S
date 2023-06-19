package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

//go:generate controller-gen object paths=$GOFILE

type FunctionSpec struct {
	Endpoint string `json:"endpoint"`
	Method   string `json:"method"`
}

type FunctionTarget struct {
	ContainerImage string `json:"containerImage"`
	Endpoint       string `json:"endpoint"`
	Port           int    `json:"port"`
	MaxReplicas    int    `json:"maxReplicas"`
	MinReplicas    int    `json:"minReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Function struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FunctionSpec   `json:"spec"`
	Target FunctionTarget `json:"target"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type FunctionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Function `json:"items"`
}

type PrettyFunction struct {
	Name   string         `json:"name"`
	UID    string         `json:"uid,omitempty"`
	Spec   FunctionSpec   `json:"spec"`
	Target FunctionTarget `json:"target"`
}

// parse K8S Function object in REST API Schema
func (function *Function) Prettify() PrettyFunction {
	prettifiedFunction := PrettyFunction{
		Name:   function.ObjectMeta.Name,
		UID:    string(function.ObjectMeta.UID),
		Spec:   function.Spec,
		Target: function.Target,
	}
	return prettifiedFunction
}

// parse K8S Function object in REST API Schema
func (functionlist *FunctionList) Prettify() []PrettyFunction {
	outputArray := make([]PrettyFunction, len(functionlist.Items))
	for i, f := range functionlist.Items {
		outputArray[i] = f.Prettify()
	}
	return outputArray
}
