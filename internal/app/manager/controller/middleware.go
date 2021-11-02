package controller

import (
	"afire/internal/pkg/gid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sunreaver/logger"
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

// SetRequestID 设置请求X-Request-ID.
func SetRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var id string
		if reqid := c.Request.Header.Get(XRequestID); len(reqid) == 0 {
			id = uuid.New().String()
			c.Request.Header.Add(XRequestID, id)
		} else {
			id = reqid
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
