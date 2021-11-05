package route

import (
	"afire/internal/app/manager/controller"
	"github.com/gin-gonic/gin"
)

func addUserRoute(r *gin.RouterGroup) {
	user := r.Group("/user")

	s := user.Group("/self")
	m := user.Group("/manager")

	s.POST("/login", controller.UserInfo)
	s.POST("/logout")

	m.POST("/login")
}
