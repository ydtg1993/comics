package model

import (
	"comics/global/orm"
	"gorm.io/gorm"
	"time"
)

type SourceChapter struct {
	Id              int    `json:"id" gorm:"primarykey"`
	ComicId         int    `json:"comic_id"`
	Source          int    `json:"source"`
	SourceId        int    `json:"source_id"`
	SourceChapterId int    `json:"source_chapter_id"`
	Sort            int    `json:"sort"`
	IsFree          int    `json:"is_free"`
	SourceUri       string `json:"source_uri"`
	Cover           string `json:"cover"`
	Title           string `json:"title"`
	SourceData      string `json:"source_data"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

const SourceChapterTASK = "source:comic:chapter"

/**
指定表名
*/
func (SourceChapter) TableName() string {
	return "source_chapter"
}

func (d *SourceChapter) Create() (err error) {
	err = orm.Eloquent.Create(&d).Error
	return
}

func (ma *SourceChapter) BeforeCreate(tx *gorm.DB) (err error) {
	ma.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}

func (ma *SourceChapter) BeforeUpdate(tx *gorm.DB) (err error) {
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}
