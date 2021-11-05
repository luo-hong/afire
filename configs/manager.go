package configs

import (
	"github.com/BurntSushi/toml"
	"github.com/sunreaver/antman/v2/db"
	"github.com/sunreaver/antman/v2/redis"
	"github.com/sunreaver/tomlanalysis/bytesize"
	"github.com/sunreaver/tomlanalysis/timesize"
)

type ManagerConfig struct {
	HTTP  ServerConfig         `toml:"server"`
	RMG   RMGConfig            `toml:"rmg"`
	DB    map[string]db.Config `toml:"db"`
	Redis redis.Config         `toml:"redis"`
	Log   LogConfig            `toml:"log"`
	Debug DebugConfig          `toml:"debug"`
}

type ServerConfig struct {
	Listen           string            `toml:"listen"`
	Mode             string            `toml:"mode"`
	WriteDir         string            `toml:"write_dir"`
	CookieTimeout    timesize.Duration `toml:"cookie_timeout"`
	CardIntervalTime timesize.Duration `toml:"card_interval_time"`
}

type RMGConfig struct {
	CornTime string `toml:"corn_time"`
}

type LogConfig struct {
	Path           string         `toml:"path"`
	Level          string         `toml:"level"`
	MaxSizeOneFile bytesize.Int64 `toml:"max_size_one_file"`
}

type DebugConfig struct {
	Port string `toml:"pprof_port"`
}

func NewManagerConfig(file string) (*ManagerConfig, error) {
	var cfg ManagerConfig
	_, e := toml.DecodeFile(file, &cfg)

	return &cfg, e
}
