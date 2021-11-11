package route

import (
	"afire/internal/app/manager/controller"
	"github.com/gin-gonic/gin"
)

func addCharacterRoute(r *gin.RouterGroup) {
	character := r.Group("/character")
	character.PUT("", controller.AddCharacter) // 新增角色
	character.PUT("/:cid")                     // 更新角色
	character.GET("/list")                     // 获取所有角色列表，如果有name参数就查和这个name相关的角色列表
	character.GET("/user/:cid")                // 根据角色id查询用户，如果数据多后期要做筛选
	character.DELETE("/:cid")                  // 删除角色
	character.DELETE("/:cid/:uid")             // 删除角色下某个用户
}
