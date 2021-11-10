package service

import (
	"afire/configs"
	"afire/internal/app/manager/business"
	"afire/internal/app/manager/route"
	"afire/internal/app/manager/service/exit"
	"afire/internal/pkg/catch"
	"afire/internal/pkg/database"
	"afire/internal/pkg/gid"
	"afire/pkg/models"
	"afire/pkg/tool"
	"afire/version"
	"flag"
	"github.com/pkg/errors"
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
	v := flag.Bool("v", false, "show version")
	flag.Parse()

	if *v {
		log.Println(":\n%v", version.Show())
		os.Exit(0)
	}

	if e := tool.SetDBKey([]byte(os.Getenv("DBKEY"))); e != nil {
		log.Fatalln("get Env: \"DBKEY\" err:", e.Error())
	}

	//解析配置文件
	c, e := configs.NewManagerConfig(configFile)
	if e != nil {
		return errors.Wrap(e, "start")
	}

	//连接db
	log.Println("====连接db====")
	if e := database.InitDateBase(c.DB); e != nil {
		return errors.Wrap(e, "init database")
	}

	//连接redis
	log.Println("====连接redis====")
	if e := catch.InitRedis(c.Redis); e != nil {
		return errors.Wrap(e, "init redis")
	}

	//连接logger
	log.Println("====连接logger====")
	//创建logger
	log.Println("++++创建logger++++")
	// 创建logger
	if e := logger.InitLoggerWithConfig(logger.Config{
		Path:     c.Log.Path,
		Loglevel: logger.LevelString(c.Log.Level),
		MaxSize:  c.Log.MaxSizeOneFile.AsInt() / 1024 / 1024, // 单位由byte转为MB
	}, nil, gid.GetGidMap()); e != nil {
		return errors.Wrap(e, "init log")
	}

	//初始化表
	log.Println("models")
	// 初始化表
	e = models.InitModels(database.AFIREMaster())
	if e != nil {
		return errors.Wrap(e, "init models")
	}

	log.Println("====logger正常====")
	cfg = *c
	l = logger.GetSugarLogger("manager.log")

	l.Infow("started", "cfg", cfg)

	log.Println("started")
	business.SetLogger(l)

	return nil
}
