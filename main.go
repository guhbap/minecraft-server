package main

import (
	"github.com/fatih/color"

	"github.com/golangmc/minecraft-server/impl"
	"github.com/golangmc/minecraft-server/impl/conf"
)

func main() {
	color.NoColor = false
	config := readConfig("serverConfig.yml")

	server := impl.NewServer(config)
	server.Load()
}

func readConfig(filename string) (config conf.ServerInfo) {
	// read from yaml file

	return *conf.NewServerInfo(filename)
}
