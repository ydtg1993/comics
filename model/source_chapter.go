package model

import (
	"comics/global/orm"
	"comics/tools/config"
	"gorm.io/gorm"
	"time"
)

type SourceChapter struct {
	Id              int       `json:"id" gorm:"primarykey;->"`
	ComicId         int       `json:"comic_id"`
	Source          int       `json:"source"`
	SourceChapterId int       `json:"source_chapter_id"`
	Sort            int       `json:"sort"`
	IsFree          int       `json:"is_free"`
	SourceUrl       string    `json:"source_url"`
	Cover           string    `json:"cover"`
	Title           string    `json:"title"`
	SourceData      string    `json:"source_data"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

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
	ma.CreatedAt = time.Now()
	ma.UpdatedAt = time.Now()
	return
}

func (ma *SourceChapter) BeforeUpdate(tx *gorm.DB) (err error) {
	ma.UpdatedAt = time.Now()
	return
}

func (ma *SourceChapter) Exists(comicId int, sourceUrl string) bool {
	result := orm.Eloquent.Where("source = ? and comic_id = ? and source_url = ?",
		config.Spe.SourceId, comicId, sourceUrl).Limit(1).Find(&ma)
	if result.Error == nil && result.RowsAffected == 1 {
		return true
	}
	return false
}
