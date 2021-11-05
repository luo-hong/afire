package business

import (
	"afire/internal/pkg/database"
	"afire/pkg/models"
	"fmt"

	"github.com/pkg/errors"
)

var resList []string

func InitCharacterRules(resources *Resources) (e error) {
	// 检索出所有角色、资源对照表
	s := models.CharacterResourceSelector{}
	// 此处使用master，因为每次更新都需要重载，slave可能存在延迟
	crs, e := s.Find(database.AFIREMaster(), "CID", "ResourceID")
	if e != nil {
		return errors.Wrap(e, "select character_resource")
	}

	// resourceID 对应的所有角色列表
	r2cids := map[string][]string{}
	for _, cr := range crs {
		v, ok := r2cids[cr.ResourceID]
		if !ok {
			v = make([]string, 0)
		}
		v = append(v, fmt.Sprintf("%v", cr.CID))
		r2cids[cr.ResourceID] = v
	}
	defer func() {
		if e == nil {
			e = InitRules()
		}
	}()

	resList = []string{}

	return resources.Range(func(resourceID string, method, route []string) error {
		resList = append(resList, resourceID)
		return addRules2Grbac(method, route, r2cids[resourceID])
	})
}
