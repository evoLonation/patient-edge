package main

import (
	"patient-edge/cloud/listen"
	"patient-edge/cloud/rpc/server"
	"patient-edge/cloud/service"
	"patient-edge/common"
	"patient-edge/config"
)

func main() {
	config := config.ParseConfig().Cloud
	uploadTemperatureSvc, doctorSvc, patientSvc := service.NewServices(config.Service)
	listen.Start(config.Listen, patientSvc, doctorSvc)
	server.Start(uploadTemperatureSvc, config.RpcServer)
	common.Block()
}
