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
	F2S        F2SConfigMapF2S        `yaml:"f2s"`
}

type F2SConfigMapPrometheus struct {
	URL string `yaml:"url"`
}
type F2SConfigMapF2S struct {
	Timeouts F2SConfigMapF2STimeouts `yaml:"timeouts"`
}

type F2SConfigMapF2STimeouts struct {
	RequestTimeout int `yaml:"request_timeout"`
	HttpTimeout    int `yaml:"http_timeout"`
	ScalingTimeout int `yaml:"scaling_timeout"`
}
