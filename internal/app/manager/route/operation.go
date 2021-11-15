package route

import (
	"afire/internal/app/manager/controller"
	"github.com/gin-gonic/gin"
)

func addOperationRoute(r *gin.RouterGroup) {
	operator := r.Group("/operation")
	operator.GET("", controller.PageChecker(), controller.OperationList) // 获取操作记录列表
	operator.GET("types", controller.OperationTypesList)                 // 获取操作记录类型
}
