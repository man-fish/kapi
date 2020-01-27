package config

import (
	"encoding/json"
	"io/ioutil"
)

type AppConfig struct{
	MysqlDsn			string		`json:"mysql_dsn"`
	StaticPath			string		`json:"static_path"`
	Port				string		`json:"port"`
	SecurityKey			string		`json:"security_key"`
	SecurityExpiresIn	int			`json:"security_expiresIn"`
}

func GetConfig(configPath string) (appConfig *AppConfig, err error) {
	appConfig = new(AppConfig)
	/* 这里你不初始化的话AppConfig是nil值。 */
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, appConfig)
	return
}

