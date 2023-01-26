package model

import (
	_ "database/sql"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"os"
	"strconv"
	"time"
)

var db *gorm.DB

func init() {
	mysqllifetime, err := strconv.Atoi(os.Getenv("MYSQL_LIFE_TIME"))
	if err != nil {
		panic(err.Error())
	}
	mysqlidletime, err := strconv.Atoi(os.Getenv("MYSQL_IDLE_TIME"))
	if err != nil {
		panic(err.Error())
	}
	mysqlmaxconn, err := strconv.Atoi(os.Getenv("MYSQL_MAX_CONN"))
	if err != nil {
		panic(err.Error())
	}

	timeout := "10s" //连接超时，10秒

	//拼接下dsn参数, dsn格式可以参考上面的语法，这里使用Sprintf动态拼接dsn参数，因为一般数据库连接参数，我们都是保存在配置文件里面，需要从配置文件加载参数，然后拼接dsn。
	db1Dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s&readTimeout=30s&writeTimeout=60s", os.Getenv("MYSQL_NAME_WR"), os.Getenv("MYSQL_PASSWORD_WR"), os.Getenv("MYSQL_HOST_WR"), os.Getenv("MYSQL_DB_NAME"), timeout)
	db2Dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s&readTimeout=30s&writeTimeout=60s", os.Getenv("MYSQL_NAME_RD"), os.Getenv("MYSQL_PASSWORD_RD"), os.Getenv("MYSQL_HOST_RD"), os.Getenv("MYSQL_DB_NAME"), timeout)

	fmt.Println(db1Dsn)
	fmt.Println(db2Dsn)

	db, err = gorm.Open(mysql.Open(db1Dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	err = db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(db1Dsn)}, // `db2` 作为 sources
		Replicas: []gorm.Dialector{mysql.Open(db2Dsn)},
		Policy:   dbresolver.RandomPolicy{}, // sources/replicas 负载均衡策略
	}).SetConnMaxIdleTime(time.Duration(mysqlidletime) * time.Second).SetConnMaxLifetime(time.Duration(mysqllifetime) * time.Second).SetMaxIdleConns(20).SetMaxOpenConns(100).SetMaxOpenConns(mysqlmaxconn))

	if err != nil {
		panic(err.Error())
	}
}

func GetGormDb() *gorm.DB {
	return db
}
