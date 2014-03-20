package main

import (
	. "route_register"
	"config"
	"os"
)

func main() {
	config := config.InitConfigFromFile("register_settings.yml")
	registrar := NewRegistrar(config)
	registrar.RegisterRoutes()
	os.Exit(1)
}
