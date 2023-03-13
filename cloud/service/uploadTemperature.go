package service

// 部分操作是像用户发布消息

import (
	"patient-edge/common"
	"patient-edge/entity"

	"github.com/pkg/errors"
)

type UploadTemperature struct {
	*context

	AbnormalPatientId string
}

func (p *UploadTemperature) createAbnormal(abnormal entity.Abnormal) error {
	return nil
}

func (p *UploadTemperature) UploadAbnormal(abnormal entity.Abnormal) error {
	if err := p.createAbnormal(abnormal); err != nil {
		return errors.Wrap(err, "createAbnormal error")
	}

	common.MqttPublish(p.context.mqttClient,
		"upload-temperature/notice-doctor",
		[]string{"patientId"},
		[]interface{}{p.AbnormalPatientId})
	return nil
}
