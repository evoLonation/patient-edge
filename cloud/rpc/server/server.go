package server

// 提供给其他节点的rpc服务

import (
	"patient-edge/cloud/service"
	"patient-edge/common"
	"patient-edge/config"
)

func Start(uploadTemperatureSvc *service.UploadTemperature, config config.RpcServerConf) {
	uploadTemperatureRpc := &UploadTemperature{
		uploadTemperatureSvc: uploadTemperatureSvc,
	}
	common.StartRpcServer(config, uploadTemperatureRpc)
}
