package cloudclient

import (
	"log"
	"net/rpc"
	"patient-edge/config"

	"github.com/pkg/errors"
)

func NewRpcClients(config config.RpcClientConf) (uploadTemperature *UploadTemperature) {
	address := config.Address + ":" + config.Port
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		log.Fatal(errors.Wrap(err, "dialing rpc server error"))
	}
	uploadTemperature = &UploadTemperature{
		client: client,
	}
	return
}
