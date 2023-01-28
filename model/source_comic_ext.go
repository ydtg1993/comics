package model

import (
	"comics/global/orm"
	"database/sql/driver"
	"encoding/json"
)

type Category []string

type SourceComicExt struct {
	Id           int      `json:"id" gorm:"primarykey"`
	ComicId      int      `json:"comic_id"`
	Author       string   `json:"Author"`
	Category     Category `json:"category" gorm:"type:json"`
	ChapterCount int      `json:"chapter_count"`
	LikeCount    int      `json:"like_count"`
	Popularity   int      `json:"popularity"`
	IsFree       int      `json:"is_free"`
	Description  string   `json:"description"`
	SourceData   string   `json:"source_data"`
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

func (t *Category) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, t)
}

func (t Category) Value() (driver.Value, error) {
	return json.Marshal(t)
}
