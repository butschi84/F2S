package configuration

type F2SFunctionMetadata struct {
	Id       string
	Name     string
	Endpoint string
	Method   string
}

type F2SFunctionTarget struct {
	ContainerImage string
	Endpoint       string
	Port           int
	MaxReplicas    int
	MinReplicas    int
}

type F2SFunction struct {
	Metadata F2SFunctionMetadata
	Target   F2SFunctionTarget
}

type F2SConfiguration struct {
	Functions []F2SFunction
}

var F2SFunctions []F2SFunction

func GetConfiguration() F2SConfiguration {
	Configuration := F2SConfiguration{
		Functions: []F2SFunction{
			F2SFunction{
				Metadata: F2SFunctionMetadata{
					Id:       "1",
					Name:     "Test Function",
					Endpoint: "/testfunction",
					Method:   "GET",
				},
				Target: F2SFunctionTarget{
					ContainerImage: "nginx",
					Endpoint:       "/",
					Port:           80,
					MaxReplicas:    1,
					MinReplicas:    1,
				},
			},
		},
	}

	return Configuration
}
