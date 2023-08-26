package apiserver

type RequestResult struct {
	Result string
	Data   map[string]interface{}
}

type Invocation struct {
	File       string
	Parameters string
	Result     string
}

type FizzletMode int

const (
	ModeReverseProxy FizzletMode = 0
	ModeCommandLine  FizzletMode = 1
)
