package apiserver

import "net"

// check if localhost is listening on port 8080
func isPortListening() FizzletMode {
	logging.Info("checking if port 8080 is listening...")
	conn, err := net.Dial("tcp", ":8080")
	defer conn.Close()
	if err != nil {
		logging.Info("port is not listening. use 'command line' mode")
		return ModeCommandLine
	}

	logging.Info("port is listening. use 'reverse proxy' mode")
	return ModeReverseProxy
}
