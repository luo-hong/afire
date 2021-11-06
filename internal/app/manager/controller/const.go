package controller

// OperationType
var OperationType map[string]string

const (
	OpUserAdd    = "user_add"
	OpUserAddStr = "创建用户"

	OpUserUpdate    = "user_update"
	OpUserUpdateStr = "更新用户"

	OpUserResetPwd    = "user_reset_pwd"
	OpUserResetPwdStr = "重置用户密码"
)
