package exit

import (
	"github.com/sunreaver/exit"
	"sync"
)

var (
	done chan struct{}
	wg   *sync.WaitGroup
)

func init() {
	done = make(chan struct{})
	wg = &sync.WaitGroup{}
}

func RegisterExiter() (doneWG *sync.WaitGroup, exiterChan <-chan struct{}) {
	wg.Add(1)

	return wg, done
}

func SyncWaitDone() error {
	//驻车监听kill信号
	exiter, delay := exit.RegistExiter()
	<-exiter
	//服务被kill
	close(done)
	wg.Wait()
	close(delay)

	return nil
}
