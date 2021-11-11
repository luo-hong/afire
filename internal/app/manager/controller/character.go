package controller

import (
	"afire/internal/app/manager/business"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddCharacter(c *gin.Context) {
	var form business.CharacterAddReq
	if err := c.Bind(&form); err != nil {
		log.Errorw("add_character", "err", err.Error())
		c.JSON(http.StatusBadRequest, responseWithStatus(0, "提交表单失败"+err.Error()))

		return
	}
	if err := form.Verify(); err != nil {
		log.Errorw("add_character", "err", err.Error())
		c.JSON(http.StatusBadRequest, responseWithStatus(0, err.Error()))

		return
	}
	err := business.AddChar(form)
	if err != nil {
		log.Errorw("add_character", "err", err.Error())
		c.JSON(http.StatusBadRequest, responseWithStatus(0, "新增角色失败"+err.Error()))
		_ = business.NewOperation(c.GetHeader(XRequestID), c.MustGet(userinfo).(*UserInfoInCatch),
			OpCharacterAdd, form, false, err)
		return
	}
	c.JSON(http.StatusOK, responseWithStatus(1, "新增角色成功"))
	_ = business.NewOperation(c.GetHeader(XRequestID), c.MustGet(userinfo).(*UserInfoInCatch),
		OpCharacterAdd, form, true, nil)
}
