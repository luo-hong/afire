package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type resUserInfo struct {
	Name  string `json:"name"`
	UID   string `json:"uid"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// UserInfo 用户详情
func UserInfo(c *gin.Context) {
	c.String(http.StatusOK, "hello word")
	return
}
