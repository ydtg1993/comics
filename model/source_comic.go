package model

import (
	"comics/global/orm"
	"comics/tools/config"
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Label []string

type SourceComic struct {
	Id                  int       `json:"id" gorm:"primarykey;->"`
	Source              int       `json:"source"`
	SourceId            int       `json:"source_id"`
	SourceUrl           string    `json:"source_url"`
	Cover               string    `json:"cover"`
	Title               string    `json:"title"`
	Author              string    `json:"Author"`
	Label               Label     `json:"label" gorm:"type:json"`
	Category            string    `json:"category"`
	Region              string    `json:"region"`
	ChapterCount        int       `json:"chapter_count"`
	ChapterPick         int       `json:"chapter_pick"`
	Like                string    `json:"like"`
	Popularity          string    `json:"popularity"`
	IsFree              int       `json:"is_free"`
	IsFinish            int       `json:"is_finish"`
	Retry               int       `json:"retry"`
	Description         string    `json:"description"`
	SourceData          string    `json:"source_data"`
	LastChapterUpdateAt time.Time `json:"last_chapter_update_at"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

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
	ma.CreatedAt = time.Now()
	ma.UpdatedAt = time.Now()
	return
}

func (ma *SourceComic) BeforeUpdate(tx *gorm.DB) (err error) {
	ma.UpdatedAt = time.Now()
	return
}

func (t *Label) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, t)
}

func (t Label) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (ma *SourceComic) Exists(sourceId int) bool {
	result := orm.Eloquent.Where("source = ? and source_id = ?", config.Spe.SourceId, sourceId).Limit(1).Find(&ma)
	if result.Error == nil && result.RowsAffected == 1 {
		return true
	}
	return false
}
