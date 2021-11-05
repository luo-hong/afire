package business

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/storyicon/grbac"
	"github.com/storyicon/grbac/pkg/loader"
)

var (
	rules loader.AdvancedRules
	rbac  *grbac.Controller
	once  sync.Once
)

func addRules2Grbac(method, path, authroles []string) error {
	once.Do(func() {
		rules = loader.AdvancedRules{}
	})
	rule := &loader.AdvancedRule{
		Host:   []string{"*"},
		Path:   path,
		Method: method,
	}
	if len(authroles) > 0 {
		rule.Permission = &grbac.Permission{
			AuthorizedRoles: authroles,
			AllowAnyone:     false,
		}
	} else {
		rule.Permission = &grbac.Permission{
			AllowAnyone: true,
		}
	}

	rules = append(rules, rule)
	return nil
}

func RBAC() *grbac.Controller {
	return rbac
}

func InitRules() (e error) {
	rules = append(rules,
		&loader.AdvancedRule{
			Host:   []string{"*"},
			Path:   []string{"**"},
			Method: []string{"GET"},
			Permission: &grbac.Permission{
				AllowAnyone: true,
			},
		}, // 默认所有get请求，所有都可以用
		&loader.AdvancedRule{
			Host:   []string{"*"},
			Path:   []string{"**/user/login/", "**/user/logout/", "**/user/update_pwd/"},
			Method: []string{"POST", "PUT"},
			Permission: &grbac.Permission{
				AllowAnyone: true,
			},
		}, // 登录，登出，修改密码
	)
	tmprbac, e := grbac.New(grbac.WithAdvancedRules(rules))
	if e != nil {
		return errors.Wrap(e, "grbac new")
	}
	rbac = tmprbac
	return nil
}
