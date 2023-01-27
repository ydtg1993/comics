package model

import (
	"gorm.io/gorm"
	"time"
)

type ComicExt struct {
	Id           int    `json:"id" gorm:"primarykey"`
	ComicId      string `json:"comic_id"`
	Author       string `json:"Author"`
	Category     string `json:"category"`
	ChapterCount int    `json:"chapter_count"`
	LikeCount    int    `json:"like_count"`
	Popularity   int    `json:"popularity"`
	IsFreeD      int    `json:"is_free"`
	Description  string `json:"description"`
	SourceData   string `json:"source_data"`
}

/**
指定表名
*/
func (ComicExt) TableName() string {
	return "source_comic_ext"
}

func (d *ComicExt) Create() (err error) {
	err = GetGormDb().Create(&d).Error
	return
}

func (ma *ComicExt) BeforeCreate(tx *gorm.DB) (err error) {
	ma.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}

func (ma *ComicExt) BeforeUpdate(tx *gorm.DB) (err error) {
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}
