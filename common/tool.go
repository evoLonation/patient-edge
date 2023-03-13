package common

import "sync"

func Block() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
