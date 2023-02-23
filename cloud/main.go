package main

import (
	"patient-edge/cloud/rpc"
	"sync"
)

func main() {
	rpc.Start()
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
