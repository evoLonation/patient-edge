package common

import (
	"log"
	"net"
	"net/rpc"
	"patient-edge/config"

	"github.com/pkg/errors"
)

func StartRpcServer[PT any](config config.RpcServerConf, rpcServers ...PT) {
	for _, rpcServer := range rpcServers {
		rpc.Register(rpcServer)
	}
	l, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		log.Fatal("listen error", err)
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
