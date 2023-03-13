package cloudclient

// 提供给其他节点的rpc客户端

import (
	"net/rpc"
	"patient-edge/cloud/rpc/types"
	"patient-edge/entity"
)

type UploadTemperature struct {
	client *rpc.Client
}

func (p *UploadTemperature) UploadAbnormal(abnormal *entity.Abnormal) (bool, error) {
	arg := types.UploadAbnormalArg{
		Abnormal: *abnormal,
	}
	var reply bool
	err := p.client.Call("UploadTemperature.UploadAbnormal", &arg, &reply)
	return reply, err

}
