package main

import (
	"patient-edge/edge/mqtt"

	"sync"
)

func main() {
	mqtt.Start()
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
