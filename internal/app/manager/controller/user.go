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
	Name      string   `json:"name"`
	UID       string   `json:"uid"`
	Phone     string   `json:"phone"`
	Email     string   `json:"email"`
	ResIDs    []string `json:"resources"`
	ChangePWD bool     `json:"cpwd"`
}

// UserInfo 用户详情
func UserInfo(c *gin.Context) {
	ui := c.MustGet(userinfo).(*UserInfoInCatch)
	out := resUserInfo{
		Name:      ui.Name,
		UID:       ui.UID,
		Email:     ui.Email,
		Phone:     ui.Phone,
		ResIDs:    ui.Resources,
		ChangePWD: ui.ChangePWD,
	}

	c.JSON(http.StatusOK, UniversalRespByData{
		Data: out,
	})
}

// LoginReq 登录参数
type LoginReq struct {
	UID string `json:"uid" binding:"required"`
	Pwd string `json:"pwd" binding:"required"`
}

// Login 用户登录
func Login(c *gin.Context) {
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

// Logout 用户退出
func Logout(c *gin.Context) {
	ui := c.MustGet(userinfo).(*UserInfoInCatch)
	session, _ := c.Cookie(cookieName)
	cli := catch.Cli()
	e := cli.Del(catch.KeyWithPrefix(catchSessionKey + session)).Err()
	log.Infow("user_logout",
		"user", ui.Name,
		"session", session,
		"err", e,
	)

	c.JSON(http.StatusOK,
		redirect{
			URL: "",
		},
	)
}

type UpdateUserInfoReq struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

// UpdateUserInfo 用户更新信息
func UpdateUserInfo(c *gin.Context) {
	req := UpdateUserInfoReq{}
	if err := c.BindJSON(&req); err != nil{
		log.Errorw("update_user",
			"bind_err",
			err.Error())
		c.JSON(http.StatusBadRequest,responseWithStatus(0,"参数错误"+err.Error()))
		return
	}
	log.Infow("update_user","form",req)
	ui := c.MustGet(userinfo).(*UserInfoInCatch)
	var err error
	defer func() {
		if err == nil{
			// 修改缓存
			if len(req.Name) > 0 {
				ui.Name = req.Name
			}
			if len(req.Phone) >0 {
				ui.Phone = req.Phone
			}
			if len(req.Email) >0 {
				ui.Email = req.Email
			}
			data,_ := json.Marshal(ui)
			cli := catch.Cli()
			e := cli.Set(catch.KeyWithPrefix(catchUIDKey+ui.UID),data,cfg.HTTP.CookieTimeout.Duration()).Err()
			log.Debugw("update_user","new_catch",string(data),"err",e)
		}
	}()

	err = business.UpdateUser(ui.GetUID(),req.Name,req.Phone,req.Email,nil)
	if err != nil {
		log.Errorw("update_user",
			"err", err.Error(),
			"req", req)
		c.JSON(http.StatusOK, responseWithStatus(1, err.Error()))
		_ = business.NewOperation(c.GetHeader(XRequestID), ui,
			OpUserUpdate, req, false, err)
		return
	}
	_ = business.NewOperation(c.GetHeader(XRequestID), ui,
		OpUserUpdate, req, true, nil)
	c.JSON(http.StatusOK, responseWithStatus(1, "更新成功"))
}

// UserUpdatePwd 用户更新密码
func UserUpdatePwd(c *gin.Context) {

}
