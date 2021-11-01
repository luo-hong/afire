package service

import (
	"afire/configs"
	"afire/internal/app/manager/route"
	"afire/internal/app/manager/service/exit"
	"afire/pkg/tool"
	"flag"
	"github.com/sunreaver/logger"
	"log"
	"net/http"
	"os"
)

var (
	cfg configs.ManagerConfig
	l   logger.Logger
)

// Start 服务启动
func Start() error {
	go func() {
		//启动http服务
		e := route.SyncInitHTTP(cfg)
		if e != nil {
			log.Fatalln("start http:", e.Error())
		}
	}()

	go func() {
		//性能分析
		if len(cfg.Debug.Port) > 0 {
			go func() {
				e := http.ListenAndServe(":"+cfg.Debug.Port, nil)
				if e != nil {
					l.Warnw("pprof", "err", e)
				}
			}()
		}
	}()

	defer func() {
		l.Infow("stopped", "cfg", cfg)
	}()

	return exit.SyncWaitDone()
}

func init() {
	if e := start(); e != nil {
		log.Fatalln("start:", e.Error())
	}
}

func start() error {
	var configFile string
	flag.StringVar(&configFile, "c", "", "config file")
	//v := flag.Bool("v",false,"show version")
	flag.Parse()

	if e := tool.SetDBKey([]byte(os.Getenv("DBKEY"))); e != nil {
		log.Fatalln("get Env: \"DBKEY\" err:", e.Error())
	}

	//解析配置文件
}
