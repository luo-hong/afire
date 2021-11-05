package business

import (
	"afire/internal/pkg/database"
	"afire/pkg/models"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tjfoc/gmsm/sm3"
	"sort"
)

func IsAdmin(character []string) bool {
	return len(character) > 0 && character[0] == "1"
}

// CheckUser 检查用户是否匹配
func CheckUser(uid,pwd string) (*models.User,[]string,[]string,error) {
	us := models.UserSelector{
		UID: []string{uid},
	}
	u, e := us.Find(database.AFIRESlave())
	if e != nil {
		return nil, nil, nil, errors.Wrap(e, "find user")
	}
	if len(u) == 0 {
		return nil, nil, nil, errors.New("find user empty")
	}

	// 验证用户名和密码匹配
	encryptedPwd := sm3.Sm3Sum([]byte(pwd + UserDefaultPWDSalt))
	log.Infow("checkout_user_afire", "pwd", pwd, "encryptedPwd", hex.EncodeToString(encryptedPwd), "pwd_db", u[0].Pwd)
	if u[0].Pwd != hex.EncodeToString(encryptedPwd) {
		return nil, nil, nil, errors.New("username do not match password")
	}

	// UID 为主键不可能有多条记录
	ucs := models.UserCharacterSelector{
		UID: []string{u[0].UID},
	}

	out, e := ucs.CIDs(database.AFIRESlave())
	if e != nil {
		return nil, nil, nil, errors.Wrap(e, "check user roles")
	}

	sort.Ints(out)
	charaStr := make([]string, len(out))
	for index, v := range out {
		charaStr[index] = fmt.Sprintf("%v", v)
	}

	var resources []string
	if !IsAdmin(charaStr) {
		characterResourceSelector := models.NewCharacterResourceSelector(0, 0)
		// 如果是超管，则返回所有资源id
		characterResourceSelector.CID = out
		var err error
		resources, err = characterResourceSelector.Resources(database.AFIRESlave())
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "check resource")
		}
	} else {
		// resList 为程序启动时，根据 resource.xml 初始化所得
		resources = resList
	}
	sort.Strings(resources)

	return &u[0], charaStr, resources, nil
}