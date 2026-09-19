package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type TagResource struct {
	ID         int64 `gorm:"primary_key"`
	BizID      string
	TargetID   int64
	TagID      int64
	UserID     int64
	CreateTime time.Time
	UpdateTime time.Time
}

func (m *TagResource) TableName() string {
	return "tag_resource"
}

type TagResourceModel struct {
	db *gorm.DB
}

func NewTagResourceModel(db *gorm.DB) *TagResourceModel {
	return &TagResourceModel{db: db}
}

func (m *TagResourceModel) Insert(ctx context.Context, data *TagResource) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *TagResourceModel) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Where("id = ?", id).Delete(&TagResource{}).Error
}

func (m *TagResourceModel) FindOne(ctx context.Context, id int64) (*TagResource, error) {
	var result TagResource
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &result, err
}

func (m *TagResourceModel) FindByTagIDAndBizIDAndTargetID(ctx context.Context, tagId int64, bizId string, targetId int64) (*TagResource, error) {
	var result TagResource
	err := m.db.WithContext(ctx).
		Where("tag_id = ? AND biz_id = ? AND target_id = ?", tagId, bizId, targetId).
		First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &result, err
}

func (m *TagResourceModel) FindResourcesByTagID(ctx context.Context, tagId int64, cursor, limit int64) ([]*TagResource, error) {
	var result []*TagResource
	err := m.db.WithContext(ctx).
		Where("tag_id = ? AND id < ?", tagId, cursor).
		Order("id desc").
		Limit(int(limit)).
		Find(&result).Error
	return result, err
}

func (m *TagResourceModel) FindResourcesByTagIDAndBizID(ctx context.Context, tagId int64, bizId string, cursor, limit int64) ([]*TagResource, error) {
	var result []*TagResource
	err := m.db.WithContext(ctx).
		Where("tag_id = ? AND biz_id = ? AND id < ?", tagId, bizId, cursor).
		Order("id desc").
		Limit(int(limit)).
		Find(&result).Error
	return result, err
}

func (m *TagResourceModel) FindTagsByBizIDAndTargetID(ctx context.Context, bizId string, targetId int64) ([]*TagResource, error) {
	var result []*TagResource
	err := m.db.WithContext(ctx).
		Where("biz_id = ? AND target_id = ?", bizId, targetId).
		Find(&result).Error
	return result, err
}

func (m *TagResourceModel) CountByTagID(ctx context.Context, tagId int64) (int64, error) {
	var count int64
	err := m.db.WithContext(ctx).Model(&TagResource{}).Where("tag_id = ?", tagId).Count(&count).Error
	return count, err
}

func (m *TagResourceModel) CountByTagIDs(ctx context.Context, tagIds []int64) (map[int64]int64, error) {
	type tagCount struct {
		TagID int64
		Count int64
	}
	var rows []tagCount
	err := m.db.WithContext(ctx).
		Model(&TagResource{}).
		Select("tag_id, count(*) as count").
		Where("tag_id IN ?", tagIds).
		Group("tag_id").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make(map[int64]int64)
	for _, row := range rows {
		result[row.TagID] = row.Count
	}
	return result, nil
}

func (m *TagResourceModel) FindHotTagIDs(ctx context.Context, limit int) ([]int64, error) {
	var rows []struct {
		TagID int64
		Count int64
	}
	err := m.db.WithContext(ctx).
		Model(&TagResource{}).
		Select("tag_id, count(*) as count").
		Group("tag_id").
		Order("count desc").
		Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.TagID)
	}
	return ids, nil
}
