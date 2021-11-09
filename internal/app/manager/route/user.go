package route

import (
	"afire/internal/app/manager/controller"
	"github.com/gin-gonic/gin"
)

func addUserRoute(r *gin.RouterGroup) {
	user := r.Group("/user")

	s := user.Group("/self")    // 用户
	m := user.Group("/manager") // 管理员

	s.GET("", controller.UserInfo)                 // 用户详情
	s.POST("/login", controller.Login)             // 用户登录
	s.POST("/logout", controller.Logout)           // 用户退出
	s.PUT("/update", controller.UpdateUserInfo)    // 用户更新信息
	s.PUT("/update_pwd", controller.UserUpdatePwd) // 用户更新密码

	m.GET("/info/:uid", controller.UserInfoV2)                       // 管理查看用户详情
	m.GET("list", controller.PageChecker(), controller.CheckoutUser) // 获取用户列表
	m.GET("find", controller.PageChecker(), controller.UserFind)     // 查询用户列表
	m.POST("", controller.UserCreate)                                // 增加新用户
}
