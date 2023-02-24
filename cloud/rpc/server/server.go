package server

import (
	"log"
	"net"
	"net/rpc"
	"patient-edge/cloud/operation"
	"patient-edge/cloud/rpc/types"
	"patient-edge/config"

	"github.com/pkg/errors"
)

type Service int

func (t *Service) ReceiveAbnormal(req *types.ReceiveAbnormalArg, reply *bool) error {
	*reply = operation.ReceiveAbnormal(req.Value, req.PatientId)
	return nil
}

func Start() {
	cloudConfig := config.Config.Cloud
	rpc.Register(new(Service))
	l, err := net.Listen("tcp", ":"+cloudConfig.RpcPort)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Fatal(errors.Wrap(err, "accept error"))
			}
			rpc.ServeConn(conn)
		}
	}()
}
