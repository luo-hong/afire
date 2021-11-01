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
	v1 := r.Group("/v1")

	addUserRoute(v1)

	return r.Run(cfg.HTTP.Listen)
}
