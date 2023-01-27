package model

import (
	"comics/global/orm"
)

type SourceComicExt struct {
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
func (SourceComicExt) TableName() string {
	return "source_comic_ext"
}

func (d *SourceComicExt) Create() (err error) {
	err = orm.Eloquent.Create(&d).Error
	return
}
