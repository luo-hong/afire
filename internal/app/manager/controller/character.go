package controller

import (
	"afire/internal/app/manager/business"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

func UpdateCharacter(c *gin.Context) {
	chID := c.Param("cid")
	chIDInt, err := strconv.Atoi(chID)
	if err != nil {
		log.Errorw("update_character", "err", err.Error())
		c.JSON(http.StatusBadRequest, responseWithStatus(0, "类型转换失败"+err.Error()))
		return
	}
	var form business.CharacterAddReq
	if err := c.Bind(&form); err != nil {
		log.Errorw("update_character", "err", err.Error())
		c.JSON(http.StatusBadRequest, responseWithStatus(0, "提交表单失败"+err.Error()))
		return
	}
	if err := form.Verify(); err != nil {
		log.Errorw("update_character", "warn", err.Error())
		c.JSON(http.StatusBadRequest, responseWithStatus(0, err.Error()))
		return
	}
	err = business.UpdateChar(form, chIDInt)
	if err != nil {
		log.Errorw("update_character", "err", err.Error())
		c.JSON(http.StatusBadRequest, responseWithStatus(0, "更新角色失败"+err.Error()))
		_ = business.NewOperation(c.GetHeader(XRequestID), c.MustGet(userinfo).(*UserInfoInCatch),
			OpCharacterUpdate, form, false, err)
		return
	}
	c.JSON(http.StatusOK, responseWithStatus(1, "更新角色成功"))
	_ = business.NewOperation(c.GetHeader(XRequestID), c.MustGet(userinfo).(*UserInfoInCatch),
		OpCharacterUpdate, form, true, nil)
}
