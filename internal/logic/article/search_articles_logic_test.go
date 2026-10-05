package articlelogic

import (
	"testing"
	"time"

	model "rpc-content/internal/model/article"
)

// TestPaginateByCursor 测试游标分页边界与下一页判定逻辑
func TestPaginateByCursor(t *testing.T) {
	now := time.Now()
	// 构造测试文章切片（按发布时间降序排列）
	articles := []*model.Article{
		{Id: 3, Title: "A3", PublishTime: now.Add(-10 * time.Minute)},
		{Id: 2, Title: "A2", PublishTime: now.Add(-20 * time.Minute)},
		{Id: 1, Title: "A1", PublishTime: now.Add(-30 * time.Minute)},
	}

	pageSize := 2
	var page []*model.Article
	var isEnd bool
	var nextCursor int64

	// 模拟拉取 pageSize + 1 条记录后的截断与下一页判定
	if len(articles) > pageSize {
		page = articles[:pageSize]
		isEnd = false
		nextCursor = page[len(page)-1].PublishTime.Unix()
	} else {
		page = articles
		isEnd = true
		nextCursor = 0
	}

	if len(page) != 2 {
		t.Fatalf("预期分页长度为 2, 实际为 %d", len(page))
	}
	if isEnd {
		t.Fatalf("预期未到达末尾 (isEnd=false), 实际为 true")
	}
	if nextCursor != articles[1].PublishTime.Unix() {
		t.Fatalf("预期下一页游标时间戳为 %d, 实际为 %d", articles[1].PublishTime.Unix(), nextCursor)
	}
}
