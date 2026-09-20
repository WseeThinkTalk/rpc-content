package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Answer struct {
	ID         int64 `gorm:"primary_key"`
	QuestionID int64
	AuthorID   int64
	Content    string
	IsAccepted int
	LikeNum    int
	ReplyNum   int
	Status     int
	CreateTime time.Time
	UpdateTime time.Time
}

func (m *Answer) TableName() string {
	return "answer"
}

type AnswerModel struct {
	db *gorm.DB
}

func NewAnswerModel(db *gorm.DB) *AnswerModel {
	return &AnswerModel{db: db}
}

func (m *AnswerModel) Insert(ctx context.Context, data *Answer) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *AnswerModel) FindOne(ctx context.Context, id int64) (*Answer, error) {
	var result Answer
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &result, err
}

func (m *AnswerModel) UpdateFields(ctx context.Context, id int64, values map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&Answer{}).Where("id = ?", id).Updates(values).Error
}

func (m *AnswerModel) FindByQuestionId(ctx context.Context, questionId int64, cursor, limit int64) ([]*Answer, error) {
	var result []*Answer
	query := m.db.WithContext(ctx).
		Where("question_id = ? AND status = ?", questionId, 0).
		Order("is_accepted desc, like_num desc, id asc")
	if cursor > 0 {
		query = query.Where("id > ?", cursor)
	}
	err := query.Limit(int(limit)).Find(&result).Error
	return result, err
}

func (m *AnswerModel) FindAcceptedByQuestionId(ctx context.Context, questionId int64) (*Answer, error) {
	var result Answer
	err := m.db.WithContext(ctx).
		Where("question_id = ? AND is_accepted = 1 AND status = 0", questionId).
		First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &result, err
}

func (m *AnswerModel) CountByQuestionId(ctx context.Context, questionId int64) (int64, error) {
	var count int64
	err := m.db.WithContext(ctx).Model(&Answer{}).
		Where("question_id = ? AND status = 0", questionId).
		Count(&count).Error
	return count, err
}
