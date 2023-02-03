package model

import (
	"comics/global/orm"
	"comics/tools/config"
	"github.com/beego/beego/v2/core/logs"
	"gorm.io/gorm"
	"time"
)

type FailInfo struct {
	Id        int       `json:"id" gorm:"primarykey"`
	Source    int       `json:"source"`
	Type      int       `json:"type"`
	Err       string    `json:"err"`
	Url       string    `json:"url"`
	Info      string    `json:"info"`
	CreatedAt time.Time `json:"created_at"`
}

/**
指定表名
*/
func (FailInfo) TableName() string {
	return "fail_info"
}

func (ma *FailInfo) BeforeCreate(tx *gorm.DB) (err error) {
	ma.CreatedAt = time.Now()
	return
}

func RecordFail(url, info, err string, ty int) {
	logs.Info(info)
	var ma FailInfo
	ma.Source = config.Spe.SourceId
	ma.Url = url
	ma.Info = info
	ma.Err = err
	ma.Type = ty
	orm.Eloquent.Create(&ma)
}
