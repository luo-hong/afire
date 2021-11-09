package controller

var OperationType map[string]string

const (
	// Count 默认数值
	Count = 21

	OpUserAdd    = "user_add"
	OpUserAddStr = "创建用户"

	OpUserUpdate    = "user_update"
	OpUserUpdateStr = "更新用户"

	OpUserResetPwd    = "user_reset_pwd"
	OpUserResetPwdStr = "重置用户密码"
)
