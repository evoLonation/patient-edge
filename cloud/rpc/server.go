package rpc

import (
	"log"
	"net"
	"net/rpc"
	"patient-edge/cloud/operation"
	"patient-edge/config"

	"github.com/pkg/errors"
)

type ReceiveAbnormalArg struct {
	Value     float64
	PatientId string
}

type Service int

func (t *Service) ReceiveAbnormal(req *ReceiveAbnormalArg, reply *bool) error {
	*reply = operation.ReceiveAbnormal(req.Value, req.PatientId)
	return nil
}

func Start() {
	cloudConfig := config.Config.Cloud
	service := new(Service)
	rpc.Register(service)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":"+cloudConfig.RpcPort)
	if e != nil {
		log.Fatal("listen error:", e)
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
