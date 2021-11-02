package route

import (
	"afire/internal/app/manager/controller"
	"github.com/gin-gonic/gin"
)

func addUserRoute(r *gin.RouterGroup) {
	user := r.Group("/user")

	s := user.Group("/self")

	s.GET("/login", controller.UserInfo)
	s.POST("/logout")
}
