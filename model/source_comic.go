package model

import (
	"comics/global/orm"
	"gorm.io/gorm"
	"time"
)

type SourceComic struct {
	Id        int    `json:"id" gorm:"primarykey"`
	Source    int    `json:"source"`
	SourceId  int    `json:"source_id"`
	SourceUri string `json:"source_uri"`
	Cover     string `json:"cover"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

const SourceComicTASK = "source:comic:task"

/**
指定表名
*/
func (SourceComic) TableName() string {
	return "source_comic"
}

func (d *SourceComic) Create() (err error) {
	err = orm.Eloquent.Create(&d).Error
	return
}

func (ma *SourceComic) BeforeCreate(tx *gorm.DB) (err error) {
	ma.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}

func (ma *SourceComic) BeforeUpdate(tx *gorm.DB) (err error) {
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}
