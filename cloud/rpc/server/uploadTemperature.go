package server

import (
	"patient-edge/cloud/rpc/types"
	"patient-edge/cloud/service"
)

type UploadTemperature struct {
	uploadTemperatureSvc *service.UploadTemperature
}

func (t *UploadTemperature) UploadAbnormal(req *types.UploadAbnormalArg, reply *bool) error {
	return t.uploadTemperatureSvc.UploadAbnormal(req.Abnormal)
}
