package catch

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sunreaver/antman/v2/redis"
)

var (
	cli *redis.Redis
	c   redis.Config
)

func InitRedis(cfg redis.Config) error {
	tmp, e := redis.MakeClient(cfg)
	if e != nil {
		return errors.Wrap(e, "make client")
	}
	cli = tmp
	c = cfg

	fmt.Println("redis cfg:", cfg)

	return nil
}

func KeyWithPrefix(key string) string {
	return fmt.Sprintf("%v%v", c.Prefix, key)
}

func Cli() *redis.Redis {
	return cli
}

func CloseRedis() {

	if cli != nil {
		cli.Close()
	}

}
