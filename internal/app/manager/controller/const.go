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

	OpUserDelete    = "user_delete"
	OpUserDeleteStr = "删除用户"

	OpCharacterAdd    = "character_add"
	OpCharacterAddStr = "新增角色"

	OpCharacterUpdate    = "character_update"
	OpCharacterUpdateStr = "更新角色"

	OpCharacterDelete    = "character_delete"
	OpCharacterDeleteStr = "删除角色"

	OpCharacterUpdateUser    = "character_update_user"
	OpCharacterUpdateUserStr = "修改角色用户"
)
