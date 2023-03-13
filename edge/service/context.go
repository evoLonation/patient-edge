package service

import (
	"patient-edge/common"
	"patient-edge/config"

	"github.com/jmoiron/sqlx"
)

type context struct {
	// config is private
	edgeDB  *sqlx.DB
	cloudDB *sqlx.DB
	redis   string // todo
}

func newContext(config config.EdgeServiceConf) *context {
	return &context{
		edgeDB:  common.NewMysqlDB(config.EdgeDataSource),
		cloudDB: common.NewMysqlDB(config.CloudDataSource),
		//todo redis
	}
}

func NewServices(conf config.EdgeServiceConf) (uploadTemperatureService *UploadTemperature) {
	context := newContext(conf)

	uploadTemperatureService = &UploadTemperature{
		context: context,
	}
	return
}
