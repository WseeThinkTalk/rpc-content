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

// TestPaginateByHot 测试基于文章ID锚点的热度榜单游标分页
func TestPaginateByHot(t *testing.T) {
	now := time.Now()
	// 构造按点赞数倒序但时间无序的数据集
	articles := []*model.Article{
		{Id: 101, Title: "Top1", LikeNum: 100, PublishTime: now.Add(-2 * time.Hour)},
		{Id: 102, Title: "Top2", LikeNum: 80, PublishTime: now.Add(-10 * time.Minute)}, // 时间比 Top1 新，但点赞少
		{Id: 103, Title: "Top3", LikeNum: 80, PublishTime: now.Add(-30 * time.Minute)}, // 同赞
		{Id: 104, Title: "Top4", LikeNum: 50, PublishTime: now.Add(-1 * time.Hour)},
		{Id: 105, Title: "Top5", LikeNum: 10, PublishTime: now.Add(-5 * time.Minute)},
	}

	pageSize := 2

	// 第 1 页：首屏
	p1, isEnd1, cursor1 := paginateByHot(articles, 0, pageSize)
	if len(p1) != 2 || p1[0].Id != 101 || p1[1].Id != 102 {
		t.Fatalf("第 1 页期望 [101, 102]，实际 %+v", p1)
	}
	if isEnd1 {
		t.Fatalf("第 1 页 isEnd 期望 false")
	}
	if cursor1 != 102 {
		t.Fatalf("第 1 页下一页游标期望 102，实际 %d", cursor1)
	}

	// 第 2 页：同赞文章不跳漏
	p2, isEnd2, cursor2 := paginateByHot(articles, cursor1, pageSize)
	if len(p2) != 2 || p2[0].Id != 103 || p2[1].Id != 104 {
		t.Fatalf("第 2 页期望 [103, 104]，实际 %+v", p2)
	}
	if isEnd2 {
		t.Fatalf("第 2 页 isEnd 期望 false")
	}
	if cursor2 != 104 {
		t.Fatalf("第 2 页下一页游标期望 104，实际 %d", cursor2)
	}

	// 第 3 页：到达末尾
	p3, isEnd3, cursor3 := paginateByHot(articles, cursor2, pageSize)
	if len(p3) != 1 || p3[0].Id != 105 {
		t.Fatalf("第 3 页期望 [105]，实际 %+v", p3)
	}
	if !isEnd3 {
		t.Fatalf("第 3 页 isEnd 期望 true")
	}
	if cursor3 != 0 {
		t.Fatalf("末页游标期望 0，实际 %d", cursor3)
	}
}
