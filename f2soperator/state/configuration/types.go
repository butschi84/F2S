package configuration

import (
	"butschi84/f2s/state/configuration/api/types/v1alpha1"
	typesV1alpha1 "butschi84/f2s/state/configuration/api/types/v1alpha1"
	"fmt"
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
	Kafka    F2SConfigMapKafka       `yaml:"kafka"`
}

type F2SConfigMapF2STimeouts struct {
	RequestTimeout int `yaml:"request_timeout"`
	HttpTimeout    int `yaml:"http_timeout"`
	ScalingTimeout int `yaml:"scaling_timeout"`
}

// *********************************
// kafka section
// *********************************
type F2SConfigMapKafka struct {
	Enabled   bool                        `yaml:"enabled"`
	Brokers   []string                    `yaml:"brokers"`
	Listeners []F2SConfigMapKafkaListener `yaml:"listeners"`
}

type F2SConfigMapKafkaListener struct {
	Topic         string                            `yaml:"topic"`
	ConsumerGroup string                            `yaml:"consumergroup"`
	Actions       []F2SConfigMapKafkaListenerAction `yaml:"actions"`
}
type F2SConfigMapKafkaListenerAction struct {
	Triggers     []F2SConfigMapKafkaListenerActionTriggers `yaml:"triggers"`
	F2SFunctions []string                                  `yaml:"f2sfunctions"`
	Response     F2SConfigMapKafkaListenerActionResponse   `yaml:"response"`
}
type F2SConfigMapKafkaListenerActionResponse struct {
	Key string `yaml:"key"`
}

type F2SConfigMapKafkaListenerActionTriggers struct {
	Type   string `yaml:"type"`
	Filter string `yaml:"filter"`
	Value  string `yaml:"value"`
}

// *********************************
// auth section
// *********************************
type F2SConfigMapAuth struct {
	GlobalConfig  F2SConfigMapAuthGlobalConfig         `yaml:"global_config"`
	Basic         []F2SConfigMapAuthBasicUser          `yaml:"basic"`
	Token         F2SConfigMapAuthToken                `yaml:"token"`
	Authorization []F2SConfigMapAuthAuthorizationGroup `yaml:"authorization"`
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

type F2SConfigMapAuthAuthorizationGroup struct {
	Group      string   `yaml:"group"`
	Privileges []string `yaml:"privileges"`
}

// *********************************
// getters
// *********************************

// get a function from running config with specific uid
func (config *F2SConfiguration) GetFunctionByUID(uid string) (v1alpha1.PrettyFunction, error) {
	for _, item := range config.Functions.Items {
		if string(item.UID) == uid {
			return item.Prettify(), nil
		}
	}
	return v1alpha1.PrettyFunction{}, fmt.Errorf("function with id %s not found in running config", uid)
}
