package controller

import (
	"afire/internal/app/manager/business"
	"afire/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OperationForm struct {
	Operator   string `form:"operator"`
	OperatorID string `form:"operator_id"`
	Operation  string `form:"operation"`
	RequestID  string `form:"req_id"`
	Start      string `form:"start"`
	End        string `form:"end"`
}

func OperationList(c *gin.Context) {
	var form OperationForm
	if e := c.Bind(&form); e != nil {
		c.JSON(http.StatusBadRequest, responseWithStatus(http.StatusBadRequest, e.Error()))
		return
	}
	ctx := utils.ContextWithID(c.GetHeader(XRequestID))
	log.Debugw("get_operation_list",
		"form", form,
	)

	start, _ := strconv.Atoi(form.Start)
	end, _ := strconv.Atoi(form.End)
	if start == 0 || end == 0 {
		c.JSON(http.StatusBadRequest, responseWithStatus(http.StatusBadRequest, "start or end time is 0"))
		return
	}
	out, count, err := business.OperationList(ctx, c.GetInt(offset), c.GetInt(size), form.Operator, form.OperatorID, form.Operation, form.RequestID, start, end)
	if err != nil {
		log.Errorw("get_operation_list",
			"form", form,
			"err", err.Error(),
		)
		c.JSON(http.StatusOK, UniversalResp{
			Status:  1,
			Message: fmt.Sprintf("查询操作历史失败：%v", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, UniversalRespByData{
		Count:  count,
		Offset: c.GetInt(offset),
		Data:   out,
		Size:   c.GetInt(size),
	})
}

func OperationTypesList(c *gin.Context) {
	if len(OperationType) != Count {
		OperationType = map[string]string{

			OpUserAdd:             OpUserAddStr,
			OpUserUpdate:          OpUserUpdateStr,
			OpUserResetPwd:        OpUserResetPwdStr,
			OpUserDelete:          OpUserDeleteStr,
			OpCharacterAdd:        OpCharacterAddStr,
			OpCharacterUpdate:     OpCharacterUpdateStr,
			OpCharacterDelete:     OpCharacterDeleteStr,
			OpCharacterUpdateUser: OpCharacterUpdateUser,
		}

	}
	c.JSON(http.StatusOK, UniversalRespByData{
		Data: OperationType,
	})
}
