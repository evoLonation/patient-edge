package main

import (
	"patient-edge/common"
	"patient-edge/config"
	"patient-edge/edge/listen"
	"patient-edge/edge/service"
)

func main() {
	config := config.ParseConfig().Edge
	uploadTemperatureSvc := service.NewServices(config.Service)
	listen.Start(config.Listen, uploadTemperatureSvc)
	common.Block()
}
