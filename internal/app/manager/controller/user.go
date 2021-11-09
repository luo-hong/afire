package controller

import (
	"afire/internal/app/manager/business"
	"afire/internal/pkg/catch"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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
	if err := c.BindJSON(&req); err != nil {
		log.Errorw("update_user",
			"bind_err",
			err.Error())
		c.JSON(http.StatusBadRequest, responseWithStatus(0, "参数错误"+err.Error()))
		return
	}
	log.Infow("update_user", "form", req)
	ui := c.MustGet(userinfo).(*UserInfoInCatch)
	var err error
	defer func() {
		if err == nil {
			// 修改缓存
			if len(req.Name) > 0 {
				ui.Name = req.Name
			}
			if len(req.Phone) > 0 {
				ui.Phone = req.Phone
			}
			if len(req.Email) > 0 {
				ui.Email = req.Email
			}
			data, _ := json.Marshal(ui)
			cli := catch.Cli()
			e := cli.Set(catch.KeyWithPrefix(catchUIDKey+ui.UID), data, cfg.HTTP.CookieTimeout.Duration()).Err()
			log.Debugw("update_user", "new_catch", string(data), "err", e)
		}
	}()

	err = business.UpdateUser(ui.GetUID(), req.Name, req.Phone, req.Email, nil)
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

type UserUpdatePwdReq struct {
	OldPwd string `json:"old_pwd" binding:"required"`
	NewPwd string `json:"new_pwd" binding:"required"`
}

// UserUpdatePwd 用户更新密码
func UserUpdatePwd(c *gin.Context) {
	req := UserUpdatePwdReq{}
	if err := c.BindJSON(&req); err != nil {
		log.Errorw("user_update_pwd",
			"bind_err",
			err.Error())
		c.JSON(http.StatusBadRequest, responseWithStatus(0, "参数错误:"+err.Error()))
		return
	}
	ui := c.MustGet(userinfo).(*UserInfoInCatch)
	err := business.UserUpdatePwd(ui.UID, req.OldPwd, req.NewPwd)
	if err != nil {
		log.Errorw("user_update_pwd",
			"err", err.Error(),
			"uid", ui.UID,
			"req", req)
		c.JSON(http.StatusBadRequest, responseWithStatus(-1, err.Error()))
		_ = business.NewOperation(c.GetHeader(XRequestID), ui,
			OpUserUpdate, req, false, err)
		return
	}

	c.JSON(http.StatusOK, responseWithStatus(1, "更新成功"))
	_ = business.NewOperation(c.GetHeader(XRequestID), ui,
		OpUserUpdate, req, true, nil)
}

type resUserInfoV2 struct {
	Name      string                                 `json:"name"`
	UID       string                                 `json:"uid"`
	Email     string                                 `json:"email"`
	Phone     string                                 `json:"phone"`
	Resources []business.CheckoutUsersCharactersData `json:"resources"`
}

// UserInfoV2 管理查看用户信息
func UserInfoV2(c *gin.Context) {
	uid := c.Param("uid")
	u, _, resources, err := business.CheckUsers(uid)
	if err != nil {
		log.Errorw("user_info_v2",
			"uid", uid,
			"err", err.Error())
		c.JSON(http.StatusOK, responseWithStatus(-1, err.Error()))
		return
	}

	res := resUserInfoV2{
		Name:      u.Name,
		UID:       u.UID,
		Email:     u.Email,
		Phone:     u.Phone,
		Resources: resources,
	}

	c.JSON(http.StatusOK, UniversalRespByData{
		Data: res,
	})
}

// CheckoutUser 获取用户列表
func CheckoutUser(c *gin.Context) {
	var form business.CheckoutUsersForm
	if e := c.Bind(&form); e != nil {
		log.Errorw("checkout_users",
			"err", e.Error())
		c.JSON(http.StatusBadRequest, responseWithStatus(0, "参数错误"+e.Error()))
		return
	}

	form.Offset = c.GetInt(offset)
	form.Size = c.GetInt(size)
	out, count, e := business.CheckoutUsers(form)
	if e != nil {
		log.Errorw("checkout_users",
			"form", form,
			"err", e.Error())
		return
	}

	c.JSON(http.StatusOK, responseWithData(out, count, form.Size, form.Offset, ""))
}

type UserFindReq struct {
	Name string `form:"name" json:"name"`
}

// UserFind 模糊查询用户列表
func UserFind(c *gin.Context) {
	funcName := "user_find"
	reqID := c.GetHeader(XRequestID)
	req := UserFindReq{}
	if err := c.BindQuery(&req); err != nil {
		log.Errorw(funcName,
			"param_valid",
			err.Error())
		c.JSON(http.StatusBadRequest, responseWithStatus(0, "参数错误"+err.Error()))
	}
	log.Debugw(funcName, "req", req)
	offset := c.GetInt(offset)
	size := c.GetInt(size)
	list, count, err := business.UserList(req.Name, &offset, &size)
	if err != nil {
		log.Errorw(funcName,
			"err", err.Error(),
			"req_id", reqID,
			"req", req)
		c.JSON(http.StatusOK, responseWithStatus(-1, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responseWithData(list, int(count), size, offset, ""))
}

type UpdateUserManagerReq struct {
	UpdateUserInfoReq
	Character []int `json:"characters"` // 角色
}

type UserCreateReq struct {
	UID string `json:"uid" binding:"required"`
	UpdateUserManagerReq
}

func (r *UserCreateReq) Validate() error {
	if len(r.Character) == 0 {
		return errors.New("角色列表为空")
	}
	return nil
}

// UserCreate 新增创建用户
func UserCreate(c *gin.Context) {
	funcName := "user_create"
	reqID := c.GetHeader(XRequestID)
	req := UserCreateReq{}
	if err := c.BindJSON(&req); err != nil {
		log.Errorw(funcName, "bind_err", err.Error())
		c.JSON(http.StatusBadRequest, responseWithStatus(0, "参数错误"+err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		log.Errorw(funcName, "param_valid", err.Error())
		c.JSON(http.StatusBadRequest, responseWithStatus(0, err.Error()))
		return
	}

	u, err := business.UserCreate(req.UID, req.Name, req.Phone, req.Email, req.Character)
	if err != nil {
		log.Errorw(funcName, "err", err.Error(), "req_id", reqID, "req", req)
		c.JSON(http.StatusOK, responseWithStatus(-1, err.Error()))
		_ = business.NewOperation(c.GetHeader(XRequestID), c.MustGet(userinfo).(*UserInfoInCatch),
			OpUserAdd, req, false, err)
		return
	}

	_ = business.NewOperation(c.GetHeader(XRequestID), c.MustGet(userinfo).(*UserInfoInCatch),
		OpUserAdd, req, true, nil)
	c.JSON(http.StatusOK, responseWithData(u, 0, 0, 0, ""))
}

// ResetPwd 重置用户密码
func ResetPwd(c *gin.Context) {
	uid := c.Param("uid")
	err := business.ResetUserPwd(uid)
	if err != nil {
		log.Errorw("reset_pwd", "err", err.Error(), "uid", uid)
		c.JSON(http.StatusOK, responseWithStatus(-1, err.Error()))
		_ = business.NewOperation(c.GetHeader(XRequestID), c.MustGet(userinfo).(*UserInfoInCatch),
			OpUserResetPwd, uid, false, err)
		return
	}

	cli := catch.Cli()
	_ = cli.Del(catch.KeyWithPrefix(catchUIDKey + c.Param("uid"))).Err()
	_ = business.NewOperation(c.GetHeader(XRequestID), c.MustGet(userinfo).(*UserInfoInCatch),
		OpUserResetPwd, uid, true, nil)
	c.JSON(http.StatusOK, responseWithStatus(1, "重置密码成功"))
}

// UpdateUserManager 更新用户信息
func UpdateUserManager(c *gin.Context) {
	req := UpdateUserManagerReq{}
	if err := c.BindJSON(&req); err != nil {
		log.Errorw("update_user_manager", "bind_err", err.Error())
		c.JSON(http.StatusBadRequest, responseWithStatus(0, "参数错误"+err.Error()))
		return
	}
	log.Infow("update_user_manager", "form", req)

	ui := c.MustGet(userinfo).(*UserInfoInCatch)
	var err error
	defer func() {
		if err == nil {
			//踢更新用户下线
			cli := catch.Cli()
			e := cli.Del(catch.KeyWithPrefix(catchUIDKey + c.Param("uid"))).Err()
			log.Debugw("update_user_manager", "form", req,
				"editor", ui.UID, "err", e)
		}
	}()

	err = business.UpdateUser(c.Param("uid"), req.Name, req.Phone, req.Email, req.Character)
	if err != nil {
		log.Errorw("update_user", "err", err.Error(),
			"uid", c.Param("uid"), "req", req)
		c.JSON(http.StatusOK, responseWithStatus(0, err.Error()))
		_ = business.NewOperation(c.GetHeader(XRequestID), ui,
			OpUserUpdate, req, false, err)
		return
	}

	c.JSON(http.StatusOK, responseWithStatus(1, "更新成功"))
	_ = business.NewOperation(c.GetHeader(XRequestID), ui,
		OpUserUpdate, req, true, nil)
}
