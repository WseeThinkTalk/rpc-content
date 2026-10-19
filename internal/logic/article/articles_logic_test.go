package articlelogic

import (
	"testing"
)

func TestSlicePageIds(t *testing.T) {
	// 场景 1：15 篇文章点赞数全部相同（都为 100 分），pageSize=10
	articleIds := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	scores := []int64{100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}

	pageSize := 10

	// 第 1 页：首屏拉取
	p1, isEnd1 := slicePageIds(articleIds, scores, 0, 0, pageSize, true)
	if len(p1) != 10 {
		t.Fatalf("第 1 页期望 10 条，实际为 %d 条", len(p1))
	}
	if isEnd1 {
		t.Fatalf("第 1 页 isEnd 期望 false")
	}
	if p1[0] != 1 || p1[9] != 10 {
		t.Fatalf("第 1 页数据不匹配: %+v", p1)
	}

	// 第 2 页：游标指向第 10 篇 (Score=100, ArticleId=10)，翻页必须返回剩余 5 篇且不能为空
	p2, isEnd2 := slicePageIds(articleIds, scores, 100, 10, pageSize, true)
	if len(p2) != 5 {
		t.Fatalf("第 2 页期望返回剩余 5 条，实际返回 %d 条", len(p2))
	}
	if !isEnd2 {
		t.Fatalf("第 2 页命中哨兵且已到末尾，isEnd 期望 true")
	}
	if p2[0] != 11 || p2[4] != 15 {
		t.Fatalf("第 2 页数据不匹配: %+v", p2)
	}

	// 场景 2：数据量充足时，第 2 页返回条数必须精准为 pageSize (不能因为切除上一页末尾项而缩水为 pageSize-1)
	moreIds := make([]int64, 30)
	moreScores := make([]int64, 30)
	for i := 0; i < 30; i++ {
		moreIds[i] = int64(i + 1)
		moreScores[i] = int64(30 - i)
	}

	// 上一页最后一条是 id=10, score=21
	pNext, isEndNext := slicePageIds(moreIds, moreScores, 21, 10, pageSize, false)
	if len(pNext) != 10 {
		t.Fatalf("第 2 页期望精准返回 10 条，实际返回 %d 条 (缩水缺陷)", len(pNext))
	}
	if isEndNext {
		t.Fatalf("数据充足时 isEnd 期望 false")
	}
	if pNext[0] != 11 || pNext[9] != 20 {
		t.Fatalf("第 2 页内容期望 [11..20]，实际: %+v", pNext)
	}
}
