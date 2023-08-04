package apiserver

type test struct {
	Result string
	Data   string
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
