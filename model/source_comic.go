package model

import (
	"gorm.io/gorm"
	"time"
)

type Comic struct {
	Id        int    `json:"id" gorm:"primarykey"`
	Source    string `json:"source"`
	SourceId  int    `json:"source_id"`
	SourceUri string `json:"source_uri"`
	Cover     int    `json:"cover"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

/**
指定表名
*/
func (Comic) TableName() string {
	return "source_comic"
}

func (d *Comic) Create() (err error) {
	err = GetGormDb().Create(&d).Error
	return
}

func (ma *Comic) BeforeCreate(tx *gorm.DB) (err error) {
	ma.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}

func (ma *Comic) BeforeUpdate(tx *gorm.DB) (err error) {
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}
