package business

import (
	"afire/configs"
	"encoding/xml"

	"github.com/pkg/errors"
)

var (
	resources = Resources{}
)

type resourceRangeFunc func(rid string, method, route []string) error

type Resources struct {
	XMLName xml.Name        `xml:"Resource"`
	Res     []ResourcesData `xml:"resource" json:"res,omitempty"`
}

func (r *Resources) Range(f resourceRangeFunc) error {
	for _, v := range r.Res {
		if e := v.Range(f); e != nil {
			return e
		}
	}
	return nil
}

// res 为空 返回 true , 不为空  返回false
func (r *Resources) isEmpty() bool {
	return len(r.Res) < 1
}

type ResourcesData struct {
	ID     string           `xml:"id" json:"id,omitempty"` // ID 一定没有为 空的
	Name   string           `xml:"name" json:"name,omitempty"`
	Method []string         `xml:"method" json:"method,omitempty"`
	Route  []string         `xml:"route" json:"route,omitempty"`
	Child  []*ResourcesData `xml:"resource" json:"child,omitempty"`
}

func (rd *ResourcesData) Range(f resourceRangeFunc) error {
	if e := f(rd.ID, rd.Method, rd.Route); e != nil {
		return e
	}
	for _, v := range rd.Child {
		if e := v.Range(f); e != nil {
			return e
		}
	}
	return nil
}

// InitResource 注意：该方法只能调用一次！
func InitResource() error {
	b, err := configs.Asset("configs/resources.xml")
	if err != nil {
		return err
	}

	err = xml.Unmarshal(b, &resources)
	if err != nil {
		return err
	}

	err = InitCharacterRules(&resources)
	if err != nil {
		return errors.Wrap(err, "init character")
	}
	return nil
}

// Resource2Json 调用该方法之前需要先 InitResource  需要手动调用！！！
func Resource2Json() (interface{}, error) {
	if resources.isEmpty() {
		return nil, errors.New("resource need init first")
	}
	return resources.Res, nil
}
