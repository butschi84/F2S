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
	Auth     F2SConfigMapAuth        `yaml:"auth"`
}

type F2SConfigMapF2STimeouts struct {
	RequestTimeout int `yaml:"request_timeout"`
	HttpTimeout    int `yaml:"http_timeout"`
	ScalingTimeout int `yaml:"scaling_timeout"`
}

// auth
type F2SConfigMapAuth struct {
	GlobalConfig F2SConfigMapAuthGlobalConfig `yaml:"global_config"`
	Basic        []F2SConfigMapAuthBasicUser  `yaml:"basic"`
	Token        F2SConfigMapAuthToken        `yaml:"token"`
}
type F2SConfigMapAuthGlobalConfig struct {
	Type string `yaml:"type"`
}

type F2SConfigMapAuthBasicUser struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Group    string `yaml:"group"`
}
type F2SConfigMapAuthToken struct {
	Tokens    []F2SConfigMapAuthTokenToken `yaml:"tokens"`
	JwtSecret string                       `yaml:"jwt_secret"`
}

type F2SConfigMapAuthTokenToken struct {
	Token string `yaml:"token"`
}
