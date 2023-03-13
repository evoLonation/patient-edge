package listen

import (
	"patient-edge/common"
	"patient-edge/config"
	"patient-edge/edge/service"
)

type uploadTemperatureReq struct {
	Temperature float64 `json:"temperature"`
	PatientId   string  `json:"patient_id"`
}

func Start(config config.ListenConf, uploadTemperatureService *service.UploadTemperature) {
	client := common.NewMqttClient(config.Mqtt)

	common.MqttSubscribe(client, "upload-temperature/upload-temperature", func(req uploadTemperatureReq) {
		uploadTemperatureService.UploadTemperature(req.Temperature, req.PatientId)
	})

}
