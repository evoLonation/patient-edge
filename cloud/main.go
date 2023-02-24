package main

import (
	"patient-edge/cloud/http"
	"patient-edge/cloud/rpc/server"
	"sync"
)

func main() {
	server.Start()
	http.Start()
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
