package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Question struct {
	ID         int64 `gorm:"primary_key"`
	Title      string
	Content    string
	AuthorID   int64
	Status     int
	AnswerNum  int
	ViewNum    int
	TagIds     string
	CreateTime time.Time
	UpdateTime time.Time
}

func (m *Question) TableName() string {
	return "question"
}

type QuestionModel struct {
	db *gorm.DB
}

func NewQuestionModel(db *gorm.DB) *QuestionModel {
	return &QuestionModel{db: db}
}

func (m *QuestionModel) Insert(ctx context.Context, data *Question) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *QuestionModel) FindOne(ctx context.Context, id int64) (*Question, error) {
	var result Question
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &result, err
}

func (m *QuestionModel) UpdateFields(ctx context.Context, id int64, values map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&Question{}).Where("id = ?", id).Updates(values).Error
}

func (m *QuestionModel) QuestionsByUserId(ctx context.Context, userId int64, sortType int, cursor, limit int64) ([]*Question, error) {
	var result []*Question
	query := m.db.WithContext(ctx).
		Where("status = ?", 0)
	if userId > 0 {
		query = query.Where("author_id = ?", userId)
	}

	if sortType == 1 {
		if cursor > 0 {
			query = query.Where("(answer_num < ? OR (answer_num = ? AND id < ?))", cursor, cursor, cursor)
		}
		query = query.Order("answer_num desc, id desc")
	} else {
		if cursor > 0 {
			query = query.Where("id < ?", cursor)
		}
		query = query.Order("id desc")
	}

	err := query.Limit(int(limit)).Find(&result).Error
	return result, err
}

func (m *QuestionModel) IncrAnswerNum(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Model(&Question{}).
		Where("id = ?", id).Update("answer_num", gorm.Expr("answer_num + 1")).Error
}

func (m *QuestionModel) DecrAnswerNum(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Model(&Question{}).
		Where("id = ? AND answer_num > 0", id).Update("answer_num", gorm.Expr("answer_num - 1")).Error
}
