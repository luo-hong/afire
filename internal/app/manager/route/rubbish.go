package route

import (
	"afire/internal/app/manager/controller"
	"github.com/gin-gonic/gin"
)

func addRubbishRoute(r *gin.RouterGroup) {
	waste := r.Group("/rubbish")
	waste.GET("", controller.PageChecker()) //TODO 获取xx垃圾信息
}
