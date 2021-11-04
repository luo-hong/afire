package business

import (
	"afire/internal/pkg/database"
	"afire/pkg/models"
	"github.com/pkg/errors"
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
}