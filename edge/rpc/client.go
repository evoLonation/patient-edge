package rpc

import (
	"log"
	"net/rpc"
	server "patient-edge/cloud/rpc"
	"patient-edge/config"

	"github.com/pkg/errors"
)

var Client *rpc.Client

func init() {
	var err error
	address := config.Config.Cloud.Address + ":" + config.Config.Cloud.RpcPort
	Client, err = rpc.Dial("tcp", address)
	if err != nil {
		log.Fatal(errors.Wrap(err, "dialing rpc server error"))
	}
}

func AddAbnormal(value float64, patientId string) bool {
	arg := server.ReceiveAbnormalArg{
		Value:     value,
		PatientId: patientId,
	}
	var reply bool
	if err := Client.Call("Service.ReceiveAbnormal", &arg, &reply); err != nil {
		log.Fatal(errors.Wrap(err, "service call error"))
	}
	return reply

}
