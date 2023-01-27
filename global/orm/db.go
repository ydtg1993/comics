package orm

import (
	goredis "github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// mysql 连接池
var Eloquent *gorm.DB

// redis 连接池
var Client *goredis.Client
