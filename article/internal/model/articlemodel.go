package model

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

var _ ArticleModel = (*customArticleModel)(nil)

type (
	// ArticleModel is the interface for article operations
	ArticleModel interface {
		Insert(ctx context.Context, data *Article) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Article, error)
		Update(ctx context.Context, data *Article) error
		Delete(ctx context.Context, id int64) error
		ArticlesByUserId(ctx context.Context, userId int64, status int, likeNum int64, pubTime, sortField string, limit int) ([]*Article, error)
		UpdateArticleStatus(ctx context.Context, id int64, status int) error
		ArticlesByUserIdWithoutCursor(ctx context.Context, userId int64, status int) ([]*Article, error)
		ArticlesAllVisible(ctx context.Context, status int) ([]*Article, error)
		SearchArticles(ctx context.Context, keyword string, authorIds []int64, status int, limit int, cursor int64) ([]*Article, error)
		ArticlesPending(ctx context.Context, limit int, cursor int64) ([]*Article, error)
		FindByIds(ctx context.Context, ids []int64) ([]*Article, error)
	}

	customArticleModel struct {
		db *gorm.DB
	}

	Article struct {
		Id          int64     `gorm:"primaryKey;column:id"` // 主键ID
		Title       string    `gorm:"column:title"`         // 标题
		Content     string    `gorm:"column:content"`       // 内容
		Cover       string    `gorm:"column:cover"`         // 封面
		Description string    `gorm:"column:description"`   // 描述
		AuthorId    int64     `gorm:"column:author_id"`     // 作者ID
		Status      int64     `gorm:"column:status"`        // 状态 0:待审核 1:审核不通过 2:可见 3:用户删除
		CommentNum  int64     `gorm:"column:comment_num"`   // 评论数
		LikeNum     int64     `gorm:"column:like_num"`      // 点赞数
		CollectNum  int64     `gorm:"column:collect_num"`   // 收藏数
		ViewNum     int64     `gorm:"column:view_num"`      // 浏览数
		ShareNum    int64     `gorm:"column:share_num"`     // 分享数
		TagIds      string    `gorm:"column:tag_ids"`       // 标签ID
		PublishTime time.Time `gorm:"column:publish_time"`  // 发布时间
		CreateTime  time.Time `gorm:"column:create_time;autoCreateTime"` // 创建时间
		UpdateTime  time.Time `gorm:"column:update_time;autoUpdateTime"` // 最后修改时间
	}
)

func (Article) TableName() string {
	return "article"
}

// sqlResult implements sql.Result
type sqlResult struct {
	id       int64
	affected int64
}

func (r sqlResult) LastInsertId() (int64, error) { return r.id, nil }
func (r sqlResult) RowsAffected() (int64, error) { return r.affected, nil }

// NewArticleModel returns a model for the database table.
func NewArticleModel(db *gorm.DB) ArticleModel {
	return &customArticleModel{
		db: db,
	}
}

func (m *customArticleModel) Insert(ctx context.Context, data *Article) (sql.Result, error) {
	err := m.db.WithContext(ctx).Create(data).Error
	if err != nil {
		return nil, err
	}
	return sqlResult{id: data.Id, affected: 1}, nil
}

func (m *customArticleModel) FindOne(ctx context.Context, id int64) (*Article, error) {
	var resp Article
	err := m.db.WithContext(ctx).First(&resp, id).Error
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *customArticleModel) Update(ctx context.Context, data *Article) error {
	return m.db.WithContext(ctx).Save(data).Error
}

func (m *customArticleModel) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Delete(&Article{}, id).Error
}

func (m *customArticleModel) ArticlesByUserId(ctx context.Context, userId int64, status int, likeNum int64, pubTime, sortField string, limit int) ([]*Article, error) {
	var articles []*Article
	query := m.db.WithContext(ctx).Where("author_id = ? AND status = ?", userId, status)
	if sortField == "like_num" {
		query = query.Where("like_num < ?", likeNum)
	} else {
		// 解析发布时间字符串以匹配 PostgreSQL 数据类型比较
		if pubTime != "" {
			t, err := time.Parse("2006-01-02 15:04:05", pubTime)
			if err == nil {
				query = query.Where("publish_time < ?", t)
			} else {
				query = query.Where("publish_time < ?", pubTime)
			}
		}
	}
	err := query.Order(sortField + " desc").Limit(limit).Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (m *customArticleModel) UpdateArticleStatus(ctx context.Context, id int64, status int) error {
	return m.db.WithContext(ctx).Model(&Article{}).Where("id = ?", id).Update("status", status).Error
}

func (m *customArticleModel) ArticlesByUserIdWithoutCursor(ctx context.Context, userId int64, status int) ([]*Article, error) {
	var articles []*Article
	var err error
	if status == -1 {
		err = m.db.WithContext(ctx).Where("author_id = ? AND status != ?", userId, 3).Order("publish_time desc").Find(&articles).Error
	} else {
		err = m.db.WithContext(ctx).Where("author_id = ? AND status = ?", userId, status).Order("publish_time desc").Find(&articles).Error
	}
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (m *customArticleModel) ArticlesAllVisible(ctx context.Context, status int) ([]*Article, error) {
	var articles []*Article
	err := m.db.WithContext(ctx).Where("status = ?", status).Order("publish_time desc").Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (m *customArticleModel) SearchArticles(ctx context.Context, keyword string, authorIds []int64, status int, limit int, cursor int64) ([]*Article, error) {
	var articles []*Article
	query := m.db.WithContext(ctx).Where("status = ?", status)
	if keyword != "" {
		likeKeyword := "%" + keyword + "%"
		if len(authorIds) > 0 {
			query = query.Where("title ILIKE ? OR content ILIKE ? OR description ILIKE ? OR author_id IN ?", likeKeyword, likeKeyword, likeKeyword, authorIds)
		} else {
			query = query.Where("title ILIKE ? OR content ILIKE ? OR description ILIKE ?", likeKeyword, likeKeyword, likeKeyword)
		}
	}
	if cursor > 0 {
		query = query.Where("publish_time < ?", time.Unix(cursor, 0))
	}
	err := query.Order("publish_time desc").Limit(limit).Find(&articles).Error
	return articles, err
}

func (m *customArticleModel) FindByIds(ctx context.Context, ids []int64) ([]*Article, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var articles []*Article
	err := m.db.WithContext(ctx).Where("id IN ?", ids).Order("publish_time desc").Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (m *customArticleModel) ArticlesPending(ctx context.Context, limit int, cursor int64) ([]*Article, error) {
	var articles []*Article
	query := m.db.WithContext(ctx).Where("status = 0")
	if cursor > 0 {
		query = query.Where("publish_time < ?", time.Unix(cursor, 0))
	}
	err := query.Order("publish_time desc").Limit(limit).Find(&articles).Error
	return articles, err
}
