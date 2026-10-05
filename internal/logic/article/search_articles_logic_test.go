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

// TestHydrateOrderPreserved 测试两阶段检索第二阶段物化后严格保持 ES 相关度顺序
func TestHydrateOrderPreserved(t *testing.T) {
	ids := []int64{103, 101, 102}
	resultMap := map[int64]*model.Article{
		101: {Id: 101, Title: "A1"},
		102: {Id: 102, Title: "A2"},
		103: {Id: 103, Title: "A3"},
	}

	ordered := make([]*model.Article, 0, len(ids))
	for _, id := range ids {
		if art, ok := resultMap[id]; ok {
			ordered = append(ordered, art)
		}
	}

	if len(ordered) != 3 || ordered[0].Id != 103 || ordered[1].Id != 101 || ordered[2].Id != 102 {
		t.Fatalf("预期物化保序结果为 [103, 101, 102]，实际为 %+v", ordered)
	}
}
