package business

import (
	"errors"
)

var (
	errInstsEmpty             = errors.New("类型为指令集，但是没有子指令")
	errCodeNoFound            = errors.New("编码未存在")
	errInstsType              = errors.New("指令更改前后类型不一致")
	errParameterFail          = errors.New("参数格式错误")
	errCurrencySubProjectFail = errors.New("通用子项目无法启动任务")
	errInstructionNotStart    = errors.New("指令集未启用")
	errMachineNotInSubProject = errors.New("机器与服务不匹配")
	errSubprojectIsEmpty      = errors.New("该服务不存在")
	ErrPrimitted              = errors.New("无权限进行此操作")
	ErrProjectCodeNoFound     = errors.New("未查找到项目")
	errMachineCode            = errors.New("主机Code不属于该子项目")
	ErrMachineCodeNoFound     = errors.New("未找到主机")
	ErrCodeNil                = errors.New("编码不能为空")
	ErrBadRequset             = errors.New("请求失败")
	ErrNoRMGServerAddr        = errors.New("RMG Server Addr 为空")
)

const (
	CurrencySubProjectCode      = "1"
	CurrencyServerName          = "内置服务"
	CurrencyServerID            = -1
	DefaultSubProjectCode       = "0"
	instTypeCMD                 = 0
	instTypeSerial              = 1
	instTypeParallel            = 2
	OffState                    = 1
	OnState                     = 0
	NotRunning                  = 2
	ConvertFail                 = "数据转为字符串失败"
	InstructionIsNeedParameter  = 2
	InstructionNotNeedParameter = 1

	MachineOnLine              = 1 //Agent 在线
	MachineOffLine             = 2 //Agent 离线
	MissionTemplateNeedArgs    = 1
	MissionTemplateNotNeedArgs = 2

	CardPauseCome  = 1
	CardPauseReset = 2

	UserDefaultPWD     string = "123456"
	UserDefaultPWDSalt string = "mw1]XDA@b)Z"
	UserChangePWDNO    bool   = false
	UserChangePWDYes   bool   = true
)
