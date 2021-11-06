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

// IsAdmin 是否是管理员
func IsAdmin(character []string) bool {
	return len(character) > 0 && character[0] == "1"
}

// CheckUser 检查用户是否匹配
func CheckUser(uid, pwd string) (*models.User, []string, []string, error) {
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

// UpdateUser 更新用户信息  characters 如果为nil，则不修改characters
func UpdateUser(uid, name, phone, email string, characters []int) (err error) {
	us := models.UserSelector{
		UID: []string{uid},
	}
	userList, e := us.Find(database.AFIRESlave())
	if e != nil {
		return errors.Wrap(e, "find user")
	}
	if len(userList) == 0 {
		return errors.New("find user empty")
	}

	if characters != nil {
		// 批量删除用户角色关系（硬删）
		e = models.UserCharacterBatchDelete(database.AFIREMaster(), uid, []int{})
		if e != nil {
			return errors.Wrap(e, "BatchDelete UserCharacter")
		}

		ucList := make([]models.UserCharacter, 0)
		for _, v := range characters {
			ucList = append(ucList, models.UserCharacter{UID: uid, CID: v})
		}

		// 批量插入用户角色关系
		e = models.UserCharacterBatchInsert(database.AFIREMaster(), ucList)
		if e != nil {
			return errors.Wrap(e, "BatchInsert UserCharacter")
		}
	}

	// 修改用户信息
	userList[0].Name = name
	userList[0].Phone = phone
	userList[0].Email = email
	e = userList[0].Update(database.AFIREMaster())
	if e != nil {
		return errors.Wrap(e, "user_update")
	}

	return nil
}

// UserUpdatePwd 用户自己更新密码
func UserUpdatePwd(uid, oldPwd, newPwd string) (err error) {
	funcName := "user_update_pwd"
	us := models.UserSelector{
		UID: []string{uid},
	}
	userList, e := us.Find(database.AFIRESlave())
	if e != nil {
		return errors.Wrap(e, "find user")
	}
	if len(userList) == 0 {
		return errors.New("find user empty")
	}
	// 验证用户名和密码匹配
	oldEncryptedPwd := sm3.Sm3Sum([]byte(oldPwd + UserDefaultPWDSalt))
	log.Infow(funcName, "uid", uid, "oldPwd", oldPwd, "oldEncryptedPwd", hex.EncodeToString(oldEncryptedPwd))
	if userList[0].Pwd != hex.EncodeToString(oldEncryptedPwd) {
		return errors.New("username do not match password")
	}
	// 更新密码
	newEncryptedPwd := sm3.Sm3Sum([]byte(newPwd + UserDefaultPWDSalt))
	log.Infow(funcName, "uid", uid, "newPwd", newPwd, "newEncryptedPwd", hex.EncodeToString(newEncryptedPwd))

	user := models.User{
		UID:       uid,
		ChangePWD: UserChangePWDYes,
		Pwd:       hex.EncodeToString(newEncryptedPwd),
	}

	e = user.UpdatePwd(database.AFIREMaster())
	if e != nil {
		log.Errorw(funcName, "update_err", e.Error())
		return errors.Wrap(e, funcName)
	}
	return nil
}

type CheckoutUsersCharactersData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CheckUsers return: User, characters, resources, error
func CheckUsers(uid string) (*models.User, []string, []CheckoutUsersCharactersData, error) {
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

	charMap := map[int]string{}
	// 查找角色详情
	cs := models.CharacterSelector{}
	chars, e := cs.Find(database.AFIRESlave(), "ID", "Name")
	if e == nil {
		for _, v := range chars {
			charMap[v.ID] = v.Name
		}
	}

	//resources := []CheckoutUsersCharactersData{}
	var resources []CheckoutUsersCharactersData
	if len(out) != 0 {
		characterResourceSelector := models.NewCharacterResourceSelector(0, 0)
		// 如果是超管，则返回所有资源id
		if !IsAdmin(charaStr) {
			characterResourceSelector.CID = out
		}
		var err error
		resourcesList, err := characterResourceSelector.ResourcesID(database.AFIRESlave())
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "check resource")
		}
		for _, v := range resourcesList {
			resources = append(resources, CheckoutUsersCharactersData{
				ID:   v,
				Name: charMap[v],
			})

		}
	}

	return &u[0], charaStr, resources, nil
}
