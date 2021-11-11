package business

import (
	"afire/internal/pkg/database"
	"afire/pkg/models"
	"github.com/pkg/errors"
)

/*
什么时候用form标签: {PREFIX}/v1/character/list?name=""&size=20&offset=0
这时的size和offset就需要form
*/

type CharacterAddReq struct {
	Name      string   `json:"name"`
	Introduce string   `json:"introduce"`
	Resources []string `json:"resources"`
}

func (c *CharacterAddReq) Verify() error {
	if len(c.Name) == 0 {
		return errors.New("角色名称为空")
	}
	if len(c.Introduce) == 0 {
		return errors.New("角色介绍为空")
	}
	if len(c.Resources) == 0 {
		return errors.New("角色资源为空")
	}

	return nil
}

// AddChar 增加角色
func AddChar(req CharacterAddReq) (e error) {
	// 插入角色name和介绍
	char := models.Character{
		Name:      req.Name,
		Introduce: req.Introduce,
	}

	defer func() {
		if e == nil {
			err := InitCharacterRules(&resources)
			if err != nil {
				log.Warnw("reload_chara_rules",
					"err", err)
			}
		}
	}()

	err := char.InsertCharacter(database.AFIREMaster())
	if err != nil {
		return err
	}

	// 插入角色资源
	for _, res := range req.Resources {
		characterRes := models.CharacterResource{
			ResourceID: res,
			CID:        char.ID,
		}
		err := characterRes.InsertCharacterResource(database.AFIREMaster())
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateChar 更新角色
func UpdateChar(req CharacterAddReq, cid int) (e error) {
	// 插入角色name和介绍
	char := models.Character{
		ID:        cid,
		Name:      req.Name,
		Introduce: req.Introduce,
	}

	defer func() {
		if e == nil {
			err := InitCharacterRules(&resources)
			if err != nil {
				log.Warnw("reload_chara_rules",
					"err", err)
			}
		}
	}()

	// 写全部用户Master,读用Slave
	err := char.Update(database.AFIREMaster())
	if err != nil {
		return err
	}
	// 删除所有旧的资源
	s := models.CharacterResource{
		CID: cid,
	}
	err = s.DeleteWithCid(database.AFIREMaster())
	if err != nil {
		return err
	}

	// 插入角色资源
	for _, res := range req.Resources {
		characterRes := models.CharacterResource{
			ResourceID: res,
			CID:        char.ID,
		}
		err := characterRes.InsertCharacterResource(database.AFIREMaster())
		if err != nil {
			return err
		}
	}

	return nil
}
