package tool

import (
	"encoding/base64"
	"regexp"

	"github.com/pkg/errors"
	sm4 "github.com/tjfoc/gmsm/sm4"
)

const (
	levelD = iota
	LevelC
	LevelB
	LevelA
	LevelS
)

func Check(minLength, maxLength, minLevel int, pwd string) error {
	if len(pwd) < minLength {
		return errors.Errorf("密码长度必须大于: %v", minLength)
	}
	if len(pwd) > maxLength {
		return errors.Errorf("密码长度必须小于: %v", maxLength)
	}

	level := levelD
	patternList := []string{`[0-9]+`, `[a-z]+`, `[A-Z]+`, `[~!@#$%^&*?_-]+`}
	for _, pattern := range patternList {
		match, _ := regexp.MatchString(pattern, pwd)
		if match {
			level++
		}
	}

	if level < minLevel {
		return errors.New("密码必须包含数字、大写字母，小写字母和以下字符 ~!@#$%^&*?_-")
	}
	return nil
}

var dbkey []byte

func SetDBKey(key []byte) error {
	dbkey = key
	if len(dbkey) != sm4.BlockSize {
		return errors.New("DBKEY len wrong")
	}
	return nil
}

func MarshalDBPwd(pwd string) (string, error) {
	out, err := sm4.Sm4CFB(dbkey, []byte(pwd), true)
	if err != nil {
		return "", nil
	}
	return base64.RawStdEncoding.EncodeToString(out), nil
}

func UnmarshalDBPwd(pwd string) (out string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("panic: %v", e)
		}
	}()
	data, e := base64.RawStdEncoding.DecodeString(pwd)
	if e != nil {
		return "", e
	}
	outBytes, e := sm4.Sm4CFB(dbkey, data, false)
	if e != nil {
		return "", e
	}
	return string(outBytes), nil
}
