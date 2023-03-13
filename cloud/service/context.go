package service

import (
	"patient-edge/common"
	"patient-edge/config"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jmoiron/sqlx"
)

type context struct {
	// config is private
	cloudDB           *sqlx.DB
	edgeDBMap         map[string]*sqlx.DB
	mqttClient        mqtt.Client //todo
	noticeDoctorTopic string
}

func newContext(config config.CloudServiceConf) *context {
	edgeDBMap := make(map[string]*sqlx.DB)
	for _, e := range config.Edges {
		edgeDBMap[e.EdgeId] = common.NewMysqlDB(e.DataSource)
	}
	return &context{
		cloudDB:   common.NewMysqlDB(config.CloudDataSource),
		edgeDBMap: edgeDBMap,
	}
}

func NewServices(conf config.CloudServiceConf) (uploadTemperatureSvc *UploadTemperature, doctorSvc *Doctor, patientSvc *Patient) {
	context := newContext(conf)
	uploadTemperatureSvc = &UploadTemperature{
		context: context,
	}
	doctorSvc = &Doctor{
		context: context,
	}
	patientSvc = &Patient{
		context: context,
	}
	return
}
