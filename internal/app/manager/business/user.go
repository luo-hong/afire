package business

import (
	"afire/internal/pkg/database"
	"afire/pkg/models"
	"afire/pkg/tool"
	"crypto/md5"
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

type CheckoutUsersForm struct {
	UID       string `form:"uid"`
	Name      string `form:"name"`
	Character int    `form:"character"`
	Offset    int    `form:"-"`
	Size      int    `form:"-"`
}

type CheckoutUsersData struct {
	UID        string                        `json:"uid"`
	Name       string                        `json:"name"`
	Phone      string                        `json:"phone"`
	Characters []CheckoutUsersCharactersData `json:"characters"`
}

// CheckoutUsers return: 用户信息、个数、错误
func CheckoutUsers(form CheckoutUsersForm) ([]CheckoutUsersData, int, error) {
	var us models.UserSelector
	if len(form.UID) > 0 {
		us.UID = []string{form.UID}
	} else if form.Character > 0 {
		//先根据角色查用户
		ucs := models.UserCharacterSelector{
			CID: []int{form.Character},
		}
		uids, e := ucs.UIDs(database.AFIRESlave())
		if e != nil {
			return nil, 0, errors.WithMessagef(e, "find uids with character: %v", form.Character)
		}
		us.UID = uids
	}
	if len(form.Name) > 0 {
		us.NameLike = tool.MakeFuzzyFiled(form.Name, tool.ContainsFuzzyFiled)
	}
	us.PageSelector = models.MakePageSelector(form.Offset, form.Size)
	count, e := us.Count(database.AFIRESlave())
	if e != nil {
		return nil, 0, errors.Wrap(e, "count")
	}

	userList, e := us.Find(database.AFIRESlave(), "UID", "Name", "Phone")
	if e != nil {
		return nil, 0, errors.WithMessagef(e, "find users with form: %v", us)
	} else if len(userList) == 0 {
		return []CheckoutUsersData{}, 0, nil
	}

	//补全用户的角色
	uids := make([]string, len(userList))
	for index, v := range userList {
		uids[index] = v.UID
	}
	//查找关联角色
	ucs := models.UserCharacterSelector{
		UID: uids,
	}
	uidAndCidList, e := ucs.Find(database.AFIRESlave(), "UID", "CID")
	if e != nil {
		return nil, 0, errors.Wrap(e, "find cids")
	}
	uid2CidMap := map[string][]int{}
	cidList := make([]int, len(uidAndCidList))
	for index, v := range uidAndCidList {
		if _, ok := uid2CidMap[v.UID]; !ok {
			uid2CidMap[v.UID] = []int{}
		}
		uid2CidMap[v.UID] = append(uid2CidMap[v.UID], v.CID)
		cidList[index] = v.CID
	}
	charMap := map[int]string{} // 角色中文名表
	if len(cidList) > 0 {
		// 查找角色的中文名
		cs := models.CharacterSelector{
			ID: cidList,
		}
		chars, e := cs.Find(database.AFIRESlave(), "ID", "Name")
		if e == nil {
			for _, v := range chars {
				charMap[v.ID] = v.Name
			}
		}
	}

	// 组装结果集合
	outList := make([]CheckoutUsersData, len(userList))
	for index, u := range userList {
		outList[index] = CheckoutUsersData{
			UID:   u.UID,
			Name:  u.Name,
			Phone: u.Phone,
		}
		uc := uid2CidMap[u.UID]
		chars := make([]CheckoutUsersCharactersData, len(uc))
		for index, v := range uc {
			c := CheckoutUsersCharactersData{
				ID:   v,
				Name: charMap[v],
			}
			chars[index] = c
		}
		outList[index].Characters = chars
	}

	return outList, int(count), nil
}

type UserFindRes struct {
	UID   string `json:"uid"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// UserList 查询用户列表
func UserList(name string, offset, size *int) (result []UserFindRes, c int64, err error) {
	us := models.UserSelector{
		UIDOrName: name,
		PageSelector: models.PageSelector{
			Offset: offset,
			Limit:  size,
		},
	}
	c, e := us.Count(database.AFIRESlave())
	if e != nil {
		return nil, 0, errors.Wrap(e, "find user")
	}
	if c == 0 {
		return []UserFindRes{}, c, nil
	}
	userList, e := us.Find(database.AFIRESlave())
	if e != nil {
		return nil, 0, errors.Wrap(e, "find user")
	}

	result = make([]UserFindRes, 0)
	for _, v := range userList {
		result = append(result, UserFindRes{
			UID:   v.UID,
			Name:  v.Name,
			Phone: v.Phone,
		})
	}

	return
}

type UserRes struct {
	UID        string   `json:"uid"`
	Name       string   `json:"name"`
	Phone      string   `json:"phone"`
	Resources  []string `json:"resources"`
	Characters []int    `json:"characters"`
}

// UserCreate 新增创建新用户
func UserCreate(uid, name, phone, email string, characters []int) (result *UserRes, err error) {
	us := models.UserSelector{
		UID: []string{uid},
	}
	userList, e := us.Find(database.AFIRESlave())
	if e != nil {
		return nil, errors.Wrap(e, "find user")
	}
	if len(userList) != 0 {
		return nil, errors.New("UID重复")
	}

	byt := []byte(UserDefaultPWD + uid)
	md5byte := md5.Sum(byt)
	pwd := hex.EncodeToString(md5byte[:])
	// 以上模拟login时传参md5(密码 + uid)
	h64byte := []byte(pwd + UserDefaultPWDSalt)
	encryptedPwd := sm3.Sm3Sum(h64byte)
	log.Infow("create_user", "uid", uid,
		"md5", hex.EncodeToString(md5byte[:]),
		"pwd", hex.EncodeToString(encryptedPwd),
	)

	user := models.User{
		UID:       uid,
		Name:      name,
		Phone:     phone,
		Email:     email,
		ChangePWD: UserChangePWDNO,
		Pwd:       hex.EncodeToString(encryptedPwd),
	}

	e = user.Insert(database.AFIREMaster())
	if e != nil {
		return nil, errors.Wrap(e, "insert user")
	}

	ucList := make([]models.UserCharacter, 0)
	for _, v := range characters {
		ucList = append(ucList, models.UserCharacter{
			UID: uid,
			CID: v,
		})
	}

	// 批量插入用户角色关系
	e = models.UserCharacterBatchInsert(database.AFIREMaster(), ucList)
	if e != nil {
		return nil, errors.Wrap(e, "BatchInsert UserCharacter")
	}

	return &UserRes{
		UID:        uid,
		Name:       name,
		Phone:      phone,
		Resources:  []string{},
		Characters: characters,
	}, nil
}

// ResetUserPwd 重置用户密码
func ResetUserPwd(uid string) (err error) {
	//更新密码
	pwd := md5.Sum([]byte(UserDefaultPWD + uid))
	newEncryptedPwd := sm3.Sm3Sum([]byte(hex.EncodeToString(pwd[:]) + UserDefaultPWDSalt))
	user := models.User{
		UID:       uid,
		ChangePWD: UserChangePWDNO,
		Pwd:       hex.EncodeToString(newEncryptedPwd),
	}
	e := user.UpdatePwd(database.AFIREMaster())
	if e != nil {
		log.Errorw("user_reset_pwd", "err", e.Error())
		return errors.Wrap(e, "user reset pwd")
	}

	return nil
}

// DeleteUser 删除用户
func DeleteUser(uid string) (err error) {
	tx := database.AFIREMaster().Begin()
	defer tx.Rollback()
	// 删除user表中的数据
	user := models.User{
		UID: uid,
	}
	userErr := user.Delete(tx)
	if userErr != nil {
		return userErr
	}

	// 删除user_character表中的数据
	char := models.UserCharacter{
		UID: uid,
	}
	charErr := char.DeleteUserCharacter(tx)
	if charErr != nil {
		return charErr
	}

	return tx.Commit().Error
}
