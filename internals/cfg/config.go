package cfg

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Cfg struct {
	Port   string
	DbName string
	DbUser string
	DbPass string
	DbHost string
	DbPort string
}

func LoadAndStoreConfig() Cfg {
	v := viper.New()
	v.SetEnvPrefix("SERV")
	v.SetDefault("PORT", "8080")
	v.SetDefault("DBNAME", "http_slurm")
	v.SetDefault("DBUSER", "root")
	v.SetDefault("DBPASS", "root")
	v.SetDefault("DBHOST", "localhost")
	v.SetDefault("DBPORT", "5432")

	var cfg Cfg

	err := v.Unmarshal(&cfg)
	if err != nil {
		logrus.Panic(err)
	}

	return cfg
}

func (cfg *Cfg) GetDBString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbPort, cfg.DbName)
}
