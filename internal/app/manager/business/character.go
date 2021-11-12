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

type GetCharListWithName struct {
	Name string `form:"name" json:"name"`
}

type CharListRes struct {
	ID         int      `json:"id,omitempty"`
	Name       string   `json:"name,omitempty"`
	Introduce  string   `json:"introduce,omitempty"`
	ResourceID []string `json:"resources"`
}

// ListChar 获取角色列表
func ListChar(offset, size int, name string) ([]CharListRes, int, error) {
	// 查出ID、Name、Introduce
	selector := models.NewCharacterSelector(offset, size)

	selector.NameLike = name

	count, err := selector.Count(database.AFIRESlave())
	if err != nil {
		return nil, 0, err
	}
	// 获取角色id，角色name，角色介绍
	out, err := selector.Find(database.AFIRESlave(), "ID", "Name", "Introduce")
	if err != nil {
		return nil, 0, err
	} else if len(out) == 0 {
		log.Warnw("list_char", "warn", "list_char_is_nil")
		return []CharListRes{}, int(count), nil
	}
	// 查出CID和ResourceID
	res := models.NewCharacterResourceSelector(0, 0)
	for _, v := range out {
		res.CID = append(res.CID, v.ID)
	}

	outRid, err := res.Find(database.AFIRESlave(), "CID", "ResourceID")
	if err != nil {
		return nil, 0, err
	}
	// 组装数据
	cid2rids := map[int][]string{} // 用map关联角色id和资源id列表
	for _, v := range outRid {
		_, ok := cid2rids[v.CID] // 查看角色id是否在map中
		if !ok {
			cid2rids[v.CID] = []string{v.ResourceID} // 不在map中，新增
		} else {
			cid2rids[v.CID] = append(cid2rids[v.CID], v.ResourceID) // 在map中，将对应的资源append数组
		}
	}

	resList := make([]CharListRes, len(out))
	for index, v := range out {
		resList[index] = CharListRes{
			ID:         v.ID,
			Name:       v.Name,
			Introduce:  v.Introduce,
			ResourceID: []string{},
		}
		if _, ok := cid2rids[v.ID]; ok {
			resList[index].ResourceID = cid2rids[v.ID]
		}
	}

	return resList, int(count), err
}

// CidGetUserInfo 角色查用户
func CidGetUserInfo(cid, offset, size int) ([]models.User, int64, error) {
	selector := models.NewUserCharacterSelector(offset, size)
	selector.CID = []int{cid}
	out, err := selector.UIDs(database.AFIRESlave())
	if err != nil {
		return nil, 0, err
	}
	if len(out) == 0 {
		log.Warnw("cid_get_user_info", "warn", "uid_count_is_zero")
		return []models.User{}, 0, err
	}
	sUser := models.NewUserSelector(offset, size)
	sUser.UID = out
	count, err := sUser.Count(database.AFIRESlave())
	if err != nil {
		return nil, 0, err
	}
	if count == 0 {
		log.Warnw("cid_get_user_info", "warn", "user_count_is_zero")
		return nil, 0, err
	}
	users, err := sUser.Find(database.AFIRESlave(), "UID", "Name","Phone","Email")
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}
