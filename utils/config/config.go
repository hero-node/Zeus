package config

import (
	"log"

	conf "github.com/robfig/config"
)

var defaultConfig *conf.Config
var env map[string]string

func GetConfigString(key string) (v string, err error) {
	v, err = defaultConfig.String(env["section"], key)
	return
}

func GetConfigInt(key string) (v int, err error) {
	v, err = defaultConfig.Int(env["section"], key)
	return
}

func GetConfigInt64(key string) (v int64, err error) {
	v32, er := defaultConfig.Int(env["section"], key)
	v = int64(v32)
	err = er
	return
}

func GetConfigBool(key string) (v bool, err error) {
	v, err = defaultConfig.Bool(env["section"], key)
	return
}

func GetConfigFloat(key string) (v float64, err error) {
	v, err = defaultConfig.Float(env["section"], key)
	return
}

func LoadConfigFile() error {
	configpath, _ := env["config"]
	c, err := conf.ReadDefault(configpath)
	if err != nil {
		return err
	}
	defaultConfig.Merge(c)
	return nil
}

func InitConfigEnv(en map[string]string) {
	var err error
	env = en
	_, ok := env["config"]
	_, ok1 := env["section"]
	if !ok || !ok1 {
		log.Fatalf("config and section not found in Env\n")
	}
	configpath, _ := env["config"]
	defaultConfig, err = conf.ReadDefault(configpath)
	if err != nil {
		log.Fatalf("ReadDefault failed. err=%v\n", err)
	}
}
