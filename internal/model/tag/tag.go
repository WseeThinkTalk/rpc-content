package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID         int64 `gorm:"primary_key"`
	TagName    string
	TagDesc    string
	CreateTime time.Time
	UpdateTime time.Time
}

func (m *Tag) TableName() string {
	return "tag"
}

type TagModel struct {
	db *gorm.DB
}

func NewTagModel(db *gorm.DB) *TagModel {
	return &TagModel{db: db}
}

func (m *TagModel) Insert(ctx context.Context, data *Tag) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *TagModel) FindOne(ctx context.Context, id int64) (*Tag, error) {
	var result Tag
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &result, err
}

func (m *TagModel) FindByName(ctx context.Context, name string) (*Tag, error) {
	var result Tag
	err := m.db.WithContext(ctx).Where("tag_name = ?", name).First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &result, err
}

func (m *TagModel) UpdateFields(ctx context.Context, id int64, values map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&Tag{}).Where("id = ?", id).Updates(values).Error
}

func (m *TagModel) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Where("id = ?", id).Delete(&Tag{}).Error
}

func (m *TagModel) FindByCursor(ctx context.Context, cursor, limit int64) ([]*Tag, error) {
	var result []*Tag
	err := m.db.WithContext(ctx).
		Where("id < ?", cursor).
		Order("id desc").
		Limit(int(limit)).
		Find(&result).Error
	return result, err
}

func (m *TagModel) FindByIds(ctx context.Context, ids []int64) ([]*Tag, error) {
	var result []*Tag
	err := m.db.WithContext(ctx).Where("id IN ?", ids).Find(&result).Error
	return result, err
}
