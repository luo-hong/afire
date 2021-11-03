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

// LoginReq 登录参数
type LoginReq struct {
	UID string `json:"uid" binding:"required"`
	Pwd string `json:"pwd" binding:"required"`
}

// UserInfo 用户详情
func UserInfo(c *gin.Context) {
	req := LoginReq{}
	if err := c.BindJSON(&req); err != nil {
		log.Errorw("login",
			"bind_err",
			err.Error())
		c.JSON(http.StatusBadRequest, responseWithStatus(1, "用户名/密码错误"))
		return
	}

	return
}
