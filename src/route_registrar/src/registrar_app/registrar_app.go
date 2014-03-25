package main

import (
	. "registrar"
	"config"
	"os"
)

func main() {
	config := config.InitConfigFromFile("registrar_settings.yml")
	registrar := NewRegistrar(config)
	registrar.RegisterRoutes()
	os.Exit(1)
}
