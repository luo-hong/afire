package business

import (
	"afire/internal/pkg/catch"
	"afire/internal/pkg/database"
	"afire/pkg/models"
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
)

var catchOperationTypes singleflight.Group

func NewOperation(reqid string, ui UserInfo, operation string, detail interface{}, result interface{}, rErr error) error {
	go func() {
		defer func() {
			e := recover()
			if e != nil {
				stack := make([]byte, 1024)
				length := runtime.Stack(stack, false)
				err := errors.Errorf("panic: %v\nstatic: %v", e, string(stack[:length]))
				log.Errorw("new_operation",
					"err", err.Error(),
				)
			}
		}()
		e := newOperation(reqid, ui, operation, detail, result, rErr)
		if e != nil {
			log.Errorw("new_operation",
				"err", e.Error(),
			)
		}
	}()
	return nil
}

func newOperation(reqid string, ui UserInfo, operation string, detail interface{}, result interface{}, rErr error) error {
	dStr := ""
	switch x := detail.(type) {
	case string:
		dStr = x
	case int, float64, bool:
		dStr = fmt.Sprintf("%v", x)
	case nil:
	default:
		d, e := json.Marshal(detail)
		if e != nil {
			log.Warnw("new_operation",
				"err", e.Error(),
			)
		} else {
			dStr = string(d)
		}
	}

	var rStr, eStr string
	if rErr != nil {
		eStr = rErr.Error()
	}
	r := struct {
		R interface{} `json:"r"`
		E string      `json:"e"`
	}{
		R: result,
		E: eStr,
	}
	byt, e := json.Marshal(r)
	if e != nil {
		log.Warnw("new_operation",
			"err", e.Error(),
		)
	} else {
		rStr = string(byt)
	}

	tmp := models.Operation{
		Operator:   ui.GetName(),
		Operation:  operation,
		OperatorID: ui.GetUID(),
		RequestID:  reqid,
		Details:    dStr,
		UtilInfo: models.UtilInfo{
			Editor: ui.GetName(),
		},
		Result: rStr,
	}
	e = tmp.Insert(database.AFIREMaster())
	if e != nil {
		return errors.Wrap(e, "insert")
	}

	return nil
}

func OperationList(ctx context.Context, offset, limit int, operator, operatorID, operation, reqID string, start, end int) (list []models.Operation, count int, err error) {
	operationSelect := models.OperationSelector{
		PageSelector: models.MakePageSelector(offset, limit),
	}

	operationSelect.StartAt = &start
	operationSelect.EndAt = &end
	if reqID != "" {
		operationSelect.RequestIDs = []string{reqID}
	}
	if operator != "" {
		operationSelect.Operator = []string{operator}
	}
	if operation != "" {
		operationSelect.Operation = []string{operation}
	}
	if operatorID != "" {
		operationSelect.OperatorID = []string{operatorID}
	}

	list, err = operationSelect.Find(database.AFIRESlave(), "ID", "Operator", "OperatorID", "Operation", "Details", "RequestID", "CreatedAt", "Editor", "Result")
	if err != nil {
		return nil, 0, errors.Wrap(err, "db find")
	}
	count, err = operationSelect.Count(database.AFIRESlave())
	if err != nil {
		log.Errorw("operation_list_count",
			"err", err.Error(),
		)
	}

	return list, count, nil
}

func OperationTypes() (allTypes []string, err error) {
	mySelect := models.OperationSelector{}
	catchKey := catch.KeyWithPrefix(":op_types")

	value, err, _ := catchOperationTypes.Do("get_operation_types", func() (interface{}, error) {
		cli := catch.Cli()
		body, e := cli.Get(catchKey).Result()
		if e != nil {
			// redis没有，则从db获取
			list, e := mySelect.GroupWithoutCondition(database.AFIRESlave(), "Operation")
			if e != nil {
				// db 也获取失败，则直接返回错误
				return nil, errors.Wrap(e, "from db")
			}

			listStr := strings.Join(list, ",")
			cli.Set(catchKey, listStr, time.Hour) // FIXME 失效实现配置化
			body = listStr
		}
		return strings.Split(body, ","), nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "get operation types")
	}

	return value.([]string), nil
}
