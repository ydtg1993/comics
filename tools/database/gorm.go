package database

import (
	"comics/global/orm"
	"comics/tools/config"
	_ "database/sql"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"time"
)

type MysqlManage struct {
}

func (m MysqlManage) Setup() (err error) {
	//拼接下dsn参数, dsn格式可以参考上面的语法，这里使用Sprintf动态拼接dsn参数，因为一般数据库连接参数，我们都是保存在配置文件里面，需要从配置文件加载参数，然后拼接dsn。
	db1Dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=%s&readTimeout=%s&writeTimeout=%s",
		config.Spe.MysqlNameWr,
		config.Spe.MysqlPasswordWr,
		config.Spe.MysqlHostWr,
		config.Spe.MysqlDbName,
		config.Spe.MysqlCharset,
		config.Spe.Timeout,
		config.Spe.ReadTimeout,
		config.Spe.WriteTimeout,
	)
	db2Dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=%s&readTimeout=%s&writeTimeout=%s",
		config.Spe.MysqlNameRd,
		config.Spe.MysqlPasswordRd,
		config.Spe.MysqlHostRd,
		config.Spe.MysqlDbName,
		config.Spe.MysqlCharset,
		config.Spe.Timeout,
		config.Spe.ReadTimeout,
		config.Spe.WriteTimeout,
	)

	newLogger := logger.New(Writer{}, logger.Config{
		SlowThreshold:             1 * time.Second,
		Colorful:                  false,
		IgnoreRecordNotFoundError: false,
		LogLevel:                  logger.Warn,
	})

	orm.Eloquent, err = gorm.Open(mysql.Open(db1Dsn), &gorm.Config{
		SkipDefaultTransaction: false,
		Logger:                 newLogger,
		PrepareStmt:            true,
	})

	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	err = orm.Eloquent.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(db1Dsn)}, // `db2` 作为 sources
		Replicas: []gorm.Dialector{mysql.Open(db2Dsn)},
		Policy:   dbresolver.RandomPolicy{}, // sources/replicas 负载均衡策略
	}).
		SetConnMaxIdleTime(time.Duration(config.Spe.MysqlIdleTime) * time.Second).
		SetConnMaxLifetime(time.Duration(config.Spe.MysqlLifeTime) * time.Second).
		SetMaxIdleConns(20).
		SetMaxOpenConns(config.Spe.MysqlMaxConn))
	return
}

type Writer struct{}

func (lw Writer) Printf(format string, v ...interface{}) {
	logs.Error(format, v...)
}
