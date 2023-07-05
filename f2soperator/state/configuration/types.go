package configuration

import (
	typesV1alpha1 "butschi84/f2s/state/configuration/api/types/v1alpha1"
)

// f2s configmap
type F2SConfiguration struct {
	Config    F2SConfigMap
	Functions *typesV1alpha1.FunctionList
}

type F2SConfigMap struct {
	Debug      bool                   `yaml:"debug"`
	Prometheus F2SConfigMapPrometheus `yaml:"prometheus"`
}

type F2SConfigMapPrometheus struct {
	URL string `yaml:"url"`
}
