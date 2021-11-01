package controller

import (
	"afire/configs"
	"afire/internal/app/manager/business"
	"github.com/sunreaver/logger"
)

var (
	log logger.Logger
	cfg configs.ManagerConfig
)

func init() {
	log = logger.Empty
}
func SetLogger(l logger.Logger) {
	log = l
	business.SetLogger(l)
}

func SetConfig(c configs.ManagerConfig) {
	cfg = c
	business.SetCfg(c)
}

type UniversalResp struct {
	Status  int    `json:"status"`
	Message string `json:"msg,omitempty"`
}

type UniversalRespByData struct {
	UniversalResp
	Count  int         `json:"count"`
	Size   int         `json:"size"`
	Offset int         `json:"offset"`
	Data   interface{} `json:"data,omitempty"`
}

type ParamImpl interface {
	Verify() error
}
