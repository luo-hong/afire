package tool

import (
	"strings"
)

type TypeOfFuzzyFiled uint8

const (
	ContainsFuzzyFiled TypeOfFuzzyFiled = 0
	SuffixFuzzyFiled   TypeOfFuzzyFiled = 1
	PrefixFuzzyFiled   TypeOfFuzzyFiled = 2
)

// 生成匹配方式.
func MakeFuzzyFiled(str string, way TypeOfFuzzyFiled) string {
	switch way {
	case SuffixFuzzyFiled:
		index := strings.Index(str, "%")
		if index != 0 {
			str = strings.Join([]string{"%", str}, "")
		}
	case PrefixFuzzyFiled:
		index := strings.Index(str, "%")
		if index != len(str)-1 {
			str = strings.Join([]string{str, "%"}, "")
		}
	case ContainsFuzzyFiled:
		index := strings.Index(str, "%")
		if index != 0 {
			str = strings.Join([]string{"%", str}, "")
		}
		if index != len(str)-1 && str[len(str)-1:] != "%" {
			str = strings.Join([]string{str, "%"}, "")
		}
	default:
	}

	return str
}
