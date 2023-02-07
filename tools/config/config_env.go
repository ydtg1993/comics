package config

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Specification struct {
	AppName  string `envconfig:"APP_NAME"`
	AppEnv   string `envconfig:"APP_ENV"`
	AppDebug bool   `envconfig:"APP_DEBUG"`

	MysqlHostWr     string `envconfig:"MYSQL_HOST_WR"`
	MysqlNameWr     string `envconfig:"MYSQL_NAME_WR"`
	MysqlPasswordWr string `envconfig:"MYSQL_PASSWORD_WR"`

	MysqlHostRd     string `envconfig:"MYSQL_HOST_RD"`
	MysqlNameRd     string `envconfig:"MYSQL_NAME_RD"`
	MysqlPasswordRd string `envconfig:"MYSQL_PASSWORD_RD"`

	MysqlCharset  string `envconfig:"MYSQL_CHARSET"`
	MysqlDbName   string `envconfig:"MYSQL_DB_NAME"`
	MysqlLifeTime int    `envconfig:"MYSQL_LIFE_TIME"`
	MysqlIdleTime int    `envconfig:"MYSQL_IDLE_TIME"`
	MysqlMaxConn  int    `envconfig:"MYSQL_MAX_CONN"`
	ReadTimeout   string `envconfig:"READ_TIMEOUT" default:"30s"`
	WriteTimeout  string `envconfig:"WRITE_TIMEOUT" default:"60s"`
	Timeout       string `envconfig:"MYSQL_TIMEOUT" default:"10000ms"`

	RedisHost string `envconfig:"REDIS_HOST"`
	RedisPort int    `envconfig:"REDIS_PORT"`
	RedisPass string `envconfig:"REDIS_PASS"`
	RedisDb   int    `envconfig:"REDIS_DB"`

	UserAgent    string `envconfig:"USER_AGENT"`
	SourceId     int    `envconfig:"SOURCE_ID"`
	SourceUrl    string `envconfig:"SOURCE_URL"`
	Maxthreads   int    `envconfig:"MAX_THREADS"`
	DownloadPath string `envconfig:"DOWNLOAD_PATH"`
	SeleniumPath string `envconfig:"SELENIUM_PATH"`

	Logday   int      `envconfig:"LOG_DAY"`
	Logpath  string   `envconfig:"LOG_PATH"`
	Loglevel []string `envconfig:"LOG_LEVEL" default:"error"`

	MongoHost     string        `envconfig:"MONGO_HOST"`
	MongoUser     string        `envconfig:"MONGO_USER"`
	MongoPassword string        `envconfig:"MONGO_PASSWORD"`
	MongoPort     string        `envconfig:"MONGO_PORT"`
	MongoDatabase string        `envconfig:"MONGO_DATABASE"`
	MongoMaxNum   uint64        `envconfig:"MONGO_MAX_NUM"`
	MongoTimeout  time.Duration `envconfig:"MONGO_TIMEOUT" default:"2s"`
}

var Spe Specification

func (s Specification) SetUp() (err error) {
	err = envconfig.Process("", &Spe)
	return
}
