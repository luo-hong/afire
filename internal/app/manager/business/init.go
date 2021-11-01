package business

import (
	"afire/configs"
	"afire/internal/app/manager/business/models"
	"afire/internal/app/manager/service/exit"
	"afire/utils"
	"github.com/pkg/errors"
	"github.com/sunreaver/logger"
	"runtime"
	"sync"
	"time"
)

var (
	log    logger.Logger
	doneWG *sync.WaitGroup
	exiter <-chan struct{}
	cfg    configs.ManagerConfig
)

func init() {
	log = logger.Empty
	doneWG, exiter = exit.RegisterExiter()
	go func() {
		<-exiter
		defer func() {
			e := recover()
			if e != nil {
				stack := make([]byte, 1024)
				length := runtime.Stack(stack, false)
				err := errors.Errorf("panic: %v\nstatic: %v", e, string(stack[:length]))
				log.Errorw("business_do_exit",
					"err", err,
				)
			}
		}()

		start := time.Now()
		log.Infow("will exit")
		//关闭 httpClient
		utils.CloseIdleHttpCli()
		log.Infow("exited", "used", time.Since(start).String())
		doneWG.Done()
	}()
	models.InitCronTimer()
}

func SetLogger(l logger.Logger) {
	log = l
}

func SetCfg(c configs.ManagerConfig) {
	cfg = c
}

type UserInfo interface {
	GetUID() string
	IsAdmin() bool
	GetName() string
}
