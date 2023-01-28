package model

import (
	"comics/global/orm"
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Category []string

type SourceComic struct {
	Id           int      `json:"id" gorm:"primarykey"`
	Source       int      `json:"source"`
	SourceId     int      `json:"source_id"`
	SourceUri    string   `json:"source_uri"`
	Cover        string   `json:"cover"`
	Title        string   `json:"title"`
	Author       string   `json:"Author"`
	Category     Category `json:"category" gorm:"type:json"`
	ChapterCount int      `json:"chapter_count"`
	LikeCount    string   `json:"like_count"`
	Popularity   string   `json:"popularity"`
	IsFree       int      `json:"is_free"`
	Description  string   `json:"description"`
	SourceData   string   `json:"source_data"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
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

func (t *Category) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, t)
}

func (t Category) Value() (driver.Value, error) {
	return json.Marshal(t)
}
