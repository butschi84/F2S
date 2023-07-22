package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// F2S Spex (for use in F2S Function Model)
type FunctionSpec struct {
	Endpoint    string `json:"endpoint"`
	Method      string `json:"method"`
	Description string `json:"description"`
}

// F2S Target (for use in F2S Function Model)
type FunctionTarget struct {
	ContainerImage string `json:"containerImage"`
	Endpoint       string `json:"endpoint"`
	Port           int    `json:"port"`
	MaxReplicas    int    `json:"maxReplicas"`
	MinReplicas    int    `json:"minReplicas"`
}

// F2S Function Model
type Function struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FunctionSpec   `json:"spec"`
	Target FunctionTarget `json:"target"`
}

type FunctionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Function `json:"items"`
}

// F2S Pretty Function
// Data Model for REST API
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

// parse K8S F2SFunction object (crd) for REST API
func (functionlist *FunctionList) Prettify() []PrettyFunction {
	outputArray := make([]PrettyFunction, len(functionlist.Items))
	for i, f := range functionlist.Items {
		outputArray[i] = f.Prettify()
	}
	return outputArray
}

// parse K8S F2SFunction object (crd) for REST API
func (functionlist *FunctionList) GetNames() []string {
	// Create a string array to store the names
	names := make([]string, len(functionlist.Items))

	// Extract the name property from each object
	for i, obj := range functionlist.Items {
		names[i] = obj.Name
	}

	return names
}
