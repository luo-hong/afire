package controller

import (
	"afire/internal/app/manager/business"
	"afire/internal/pkg/catch"
	"afire/internal/pkg/gid"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/storyicon/grbac"
	"github.com/sunreaver/logger"
	"net/http"
	"strconv"
	"strings"
)

const (
	userinfo        = "_user_info_"
	offset          = "_offset_"
	size            = "_size_"
	serverid        = "_serverid_"
	dbsource        = "_dbsource_"
	dbusersource    = "_dbusersource_"
	XRequestID      = "X-Request-ID"
	cookieName      = "AFIRE_JSESSIONID"
	catchSessionKey = "user:"
	catchUIDKey     = "user:uid:"
)

type redirect struct {
	UniversalResp
	URL string `json:"redirect"`
}

type UserInfoInCatch struct {
	ChangePWD bool     `json:"pwd_c"`
	Name      string   `json:"name"`
	UID       string   `json:"uid"`
	Phone     string   `json:"phone"`
	Email     string   `json:"email"`
	Character []string `json:"chara"`
	Resources []string `json:"res"`
	IP        string   `json:"ip"`
}

func (uiic *UserInfoInCatch) GetUID() string {
	return uiic.UID
}

func (uiic *UserInfoInCatch) IsAdmin() bool {
	return business.IsAdmin(uiic.Character)
}

func (uiic *UserInfoInCatch) GetName() string {
	return uiic.Name
}

func (uiic *UserInfoInCatch) HadChangedPWD() bool {
	return uiic.ChangePWD
}

func fetchUserinfoFromCatch(session string) (*UserInfoInCatch, error) {
	cli := catch.Cli()
	// 从afire redis中获取用户其它详情
	uid, _ := cli.Get(catch.KeyWithPrefix(catchSessionKey + session)).Result()
	if len(uid) == 0 {
		return nil, errors.New("未登录（101）")
	}

	uiAFIREStr, e := cli.Get(catch.KeyWithPrefix(catchUIDKey + uid)).Bytes()
	if e != nil && len(uiAFIREStr) == 0 {
		return nil, errors.New("未登录（102）")
	}

	var ui UserInfoInCatch
	e = json.Unmarshal(uiAFIREStr, &ui)
	if e != nil {
		return nil, errors.Wrap(e, "未登录（103）")
	}

	tx := catch.Cli()
	tx.Expire(catch.KeyWithPrefix(catchSessionKey+session), cfg.HTTP.CookieTimeout.Duration())
	tx.Expire(catch.KeyWithPrefix(catchUIDKey+ui.GetUID()), cfg.HTTP.CookieTimeout.Duration())
	return &ui, nil
}

// CheckLogin 检测用户.
// 将用户信息写入gin.Context(userinfo)中.
func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasSuffix(c.FullPath(), "/login") {
			c.Next()
			return
		}
		session, _ := c.Cookie(cookieName)
		if len(session) == 0 {
			c.JSON(http.StatusUnauthorized, redirect{UniversalResp: UniversalResp{Message: "未登录(1)"}})
			c.Abort()

			return
		}
		uInfo, e := fetchUserinfoFromCatch(session)
		if e != nil {
			log.Warnw("check_login",
				"session", session,
				"msg", "未登录",
				"err", e.Error(),
			)
			c.JSON(http.StatusUnauthorized, redirect{UniversalResp: UniversalResp{Message: e.Error()}})
			c.Abort()
			return
		} else if uInfo.IP != c.ClientIP() {
			// 非GET请求，必须先修改密码后才能使用
			log.Warnw("check_login",
				"login_ip", uInfo.IP,
				"now", c.ClientIP(),
			)
			c.JSON(http.StatusUnauthorized, redirect{UniversalResp: UniversalResp{Message: "账号已经在其它地方登录，请重新登录"}})
			c.Abort()
			return
		}

		if !uInfo.HadChangedPWD() && c.Request.Method != http.MethodGet && !strings.HasSuffix(c.Request.URL.Path, "/user/self/update_pwd") {
			// 非GET请求，必须先修改密码后才能使用
			log.Warnw("check_login",
				"session", session,
				"msg", "未修改密码",
			)
			c.JSON(http.StatusOK,
				redirect{
					UniversalResp: UniversalResp{
						Status:  999,
						Message: "请您修改您的密码后重新登录使用",
					},
				},
			)
			c.Abort()
			return
		}

		c.Set(userinfo, uInfo)
		tmpLogger := log.Debugw
		if c.Request.Method == http.MethodPut ||
			c.Request.Method == http.MethodPost ||
			c.Request.Method == http.MethodDelete {
			tmpLogger = log.Infow
		}
		tmpLogger("user",
			"uid", uInfo.UID,
			"chara", uInfo.Character,
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
		)
		c.Next()
	}
}

// SetRequestID 设置请求X-Request-ID.
func SetRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var id string
		if reqID := c.Request.Header.Get(XRequestID); len(reqID) == 0 {
			id = uuid.New().String()
			c.Request.Header.Add(XRequestID, id)
		} else {
			id = reqID
		}
		c.Header(XRequestID, id)
		goroutineID := logger.GetGID()
		if goroutineID > 0 {
			gid.GetGidMap().Store(goroutineID, id)
			defer gid.GetGidMap().Delete(goroutineID)
		}
		c.Next()
	}
}

// CheckUserRole 检测用户是否有权限
func CheckUserRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasSuffix(c.FullPath(), "/login") {
			c.Next()
			return
		}
		ui := c.MustGet(userinfo).(*UserInfoInCatch)
		if ui.IsAdmin() {
			// 超管逻辑
			c.Next()
			return
		}

		tmpPath := c.Request.URL.Path
		if !strings.HasSuffix(tmpPath, "/") {
			tmpPath += "/"
		}

		p, e := business.RBAC().IsQueryGranted(&grbac.Query{
			Path:   tmpPath,
			Host:   c.Request.Host,
			Method: c.Request.Method,
		}, ui.Character)
		if e != nil {
			c.String(http.StatusOK, "校验权限失败: %v", e.Error())
			c.Abort()
			return
		} else if !p.IsGranted() {
			log.Warnw("permission",
				"res", ui.Resources,
				"err", p.String(),
			)
			c.JSON(http.StatusOK, responseWithStatus(2, business.ErrPrimitted.Error()))
			c.Abort()
			return
		}
		c.Next()
	}
}

// PageChecker 分页检查器.
// 例如使用offset，则c.GetInt(offset) .
func PageChecker() gin.HandlerFunc {
	return PageCheckerWithSize(100)
}

// PageCheckerWithSize 分页检查器.
// 例如使用offset，则c.GetInt(offset) .
func PageCheckerWithSize(maxSize int) gin.HandlerFunc {
	return func(c *gin.Context) {
		off, e1 := strconv.Atoi(c.DefaultQuery("offset", "0"))
		s, e2 := strconv.Atoi(c.DefaultQuery("size", "20"))
		if e1 != nil || e2 != nil {
			c.String(http.StatusBadRequest, "分页参数错误: size: %v, offset: %v", c.Query("size"), c.Query("offset"))
			c.Abort()

			return
		}
		if s > maxSize {
			s = maxSize
		}
		c.Set(offset, off)
		c.Set(size, s)
		c.Next()
	}
}
