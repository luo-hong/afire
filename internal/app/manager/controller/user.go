package controller

import (
	"afire/internal/app/manager/business"
	"afire/internal/pkg/catch"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		c.JSON(http.StatusBadRequest, responseWithStatus(0, "用户名/密码错误(1)"))
		return
	}

	//检查afire用户
	u, characters, resources, e := business.CheckUser(req.UID, req.Pwd)
	if e != nil {
		log.Errorw("login",
			"password_err", e.Error(),
			"uid", req.UID,
			"pwd", req.Pwd,
		)
		c.JSON(http.StatusBadRequest,
			redirect{
				UniversalResp: UniversalResp{
					Message: "用户名/密码错误（2）",
				},
			})
		return
	}

	// 生成session
	session := uuid.New().String()
	ui := UserInfoInCatch{
		ChangePWD: u.ChangePWD,
		Name:      u.Name,
		UID:       u.UID,
		Phone:     u.Phone,
		Email:     u.Email,
		Character: characters,
		Resources: resources,
		IP:        c.ClientIP(),
	}
	data, _ := json.Marshal(&ui)
	cli := catch.Cli()
	e1 := cli.Set(catch.KeyWithPrefix(catchSessionKey+session), u.UID, cfg.HTTP.CookieTimeout.Duration()).Err()
	e2 := cli.Set(catch.KeyWithPrefix(catchUIDKey+u.UID), data, cfg.HTTP.CookieTimeout.Duration()).Err()
	if e1 != nil || e2 != nil {
		log.Errorw("login", "e1", e1, "e2", e2)
		c.JSON(http.StatusUnauthorized,
			redirect{
				UniversalResp: UniversalResp{
					Message: "登录失败（6）",
				},
			})
		return
	}

	log.Infow("login", "uid", u.UID, "session", session)

	// 设置 cookie
	c.SetCookie(cookieName, session, 0, "/afire", "", false, true)
	c.JSON(http.StatusOK, UniversalResp{
		Message: "登录成功",
	})

	return
}
