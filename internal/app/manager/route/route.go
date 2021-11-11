package route

import (
	"afire/configs"
	"afire/internal/app/manager/controller"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sunreaver/logger"
)

func SyncInitHTTP(cfg configs.ManagerConfig) error {
	controller.SetLogger(logger.GetSugarLogger("manager.log"))
	controller.SetConfig(cfg)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(gzip.Gzip(gzip.BestSpeed, gzip.WithExcludedPathsRegexs([]string{
		".*/processing/.*",
		".*/sql/exec/.*",
	})))
	gin.SetMode(cfg.HTTP.Mode)
	r.Use(controller.SetRequestID())
	v1 := r.Group("/v1", controller.CheckLogin(), controller.CheckUserRole())

	addUserRoute(v1)      // 用户相关的借口
	addCharacterRoute(v1) // 角色相关的接口 说明角色指的是权限资源

	return r.Run(cfg.HTTP.Listen)
}
