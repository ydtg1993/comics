package model

import (
	"comics/global/orm"
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Images []string
type SourceImage struct {
	Id         int       `json:"id" gorm:"primarykey;->"`
	Source     int       `json:"source"`
	ComicId    int       `json:"comic_id"`
	State      int       `json:"state"`
	SourceData Images    `json:"source_data" gorm:"type:json"`
	Images     Images    `json:"images" gorm:"type:json"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ChapterId  int       `json:"chapter_id"`
}

/**
指定表名
*/
func (SourceImage) TableName() string {
	return "source_image"
}

func (d *SourceImage) Create() (err error) {
	err = orm.Eloquent.Create(&d).Error
	return
}

func (ma *SourceImage) BeforeCreate(tx *gorm.DB) (err error) {
	ma.CreatedAt = time.Now()
	ma.UpdatedAt = time.Now()
	return
}

func (ma *SourceImage) BeforeUpdate(tx *gorm.DB) (err error) {
	ma.UpdatedAt = time.Now()
	return
}

func (t *Images) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, t)
}

func (t Images) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (ma *SourceImage) Exists(chapterId int) bool {
	result := orm.Eloquent.Where("chapter_id = ?", chapterId).Limit(1).Find(&ma)
	if result.Error == nil && result.RowsAffected == 1 {
		return true
	}
	return false
}
