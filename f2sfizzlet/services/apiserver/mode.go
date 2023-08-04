package apiserver

import "net"

// check if localhost is listening on port 80
func isPortListening() FizzletMode {
	conn, err := net.Dial("tcp", ":80")
	if err != nil {
		return ModeCommandLine
	}
	defer conn.Close()

	return ModeReverseProxy
}
