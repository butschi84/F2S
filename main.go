package main

import (
	"butschi84/f2s/configuration"
	"butschi84/f2s/routes"
)

var F2SConfiguration configuration.F2SConfiguration

func main() {
	F2SConfiguration = configuration.ActiveConfiguration
	routes.HandleRequests(F2SConfiguration)
}
