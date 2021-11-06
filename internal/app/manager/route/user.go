package route

import (
	"afire/internal/app/manager/controller"
	"github.com/gin-gonic/gin"
)

func addUserRoute(r *gin.RouterGroup) {
	user := r.Group("/user")

	s := user.Group("/self")
	m := user.Group("/manager")

	s.GET("", controller.UserInfo)
	s.POST("/login", controller.Login) // 用户登录
	s.POST("/logout", controller.Logout) // 用户退出
	s.PUT("/update", controller.UpdateUserInfo) // 用户更新信息
	s.PUT("/update_pwd", controller.UserUpdatePwd) // 用户更新密码

	m.POST("/login")
}
