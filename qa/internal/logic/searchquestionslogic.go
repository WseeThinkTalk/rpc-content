package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"rpc-content/qa/internal/svc"
	"rpc-content/qa/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchQuestionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchQuestionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchQuestionsLogic {
	return &SearchQuestionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchQuestionsLogic) SearchQuestions(in *pb.SearchQuestionsRequest) (*pb.SearchQuestionsResponse, error) {
	if in.PageSize == 0 {
		in.PageSize = 20
	}

	query := buildQuestionSearchQuery(in.Keyword, in.PageSize+1, in.Cursor)

	res, err := l.svcCtx.Es.Search(
		l.svcCtx.Es.Search.WithContext(l.ctx),
		l.svcCtx.Es.Search.WithIndex("question-index"),
		l.svcCtx.Es.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		l.Errorf("[SearchQuestions] ES search err: %v keyword: %s", err, in.Keyword)
		return nil, err
	}
	defer res.Body.Close()

	var result esQuestionResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		l.Errorf("[SearchQuestions] decode err: %v", err)
		return nil, err
	}

	hits := result.Hits.Hits
	var isEnd bool
	if len(hits) > int(in.PageSize) {
		hits = hits[:in.PageSize]
	} else {
		isEnd = true
	}

	items := make([]*pb.SearchQuestionItem, 0, len(hits))
	for _, hit := range hits {
		src := hit.Source
		items = append(items, &pb.SearchQuestionItem{
			Id:         src.ID,
			Title:      src.Title,
			Content:    src.Content,
			AuthorId:   src.AuthorID,
			AnswerNum:  src.AnswerNum,
			TagIds:     src.TagIds,
			CreateTime: src.CreateUnix,
		})
	}

	var cursor int64
	if len(hits) > 0 && !isEnd {
		cursor = hits[len(hits)-1].Sort[0]
	}

	return &pb.SearchQuestionsResponse{
		Items:  items,
		Cursor: cursor,
		IsEnd:  isEnd,
	}, nil
}

func buildQuestionSearchQuery(keyword string, size, cursor int64) string {
	q := fmt.Sprintf(
		`{"query":{"bool":{"must":[{"multi_match":{"query":"%s","fields":["title^3","content"]}},{"term":{"status":0}}]}},"size":%d`,
		keyword, size,
	)
	if cursor > 0 {
		q += fmt.Sprintf(`,"search_after":[%d]`, cursor)
	}
	q += `,"sort":[{"_score":"desc"},{"id":"asc"}]}`
	return q
}

type esQuestionResult struct {
	Hits struct {
		Hits []struct {
			Source struct {
				ID         int64  `json:"id"`
				Title      string `json:"title"`
				Content    string `json:"content"`
				AuthorID   int64  `json:"author_id"`
				AnswerNum  int64  `json:"answer_num"`
				TagIds     string `json:"tag_ids"`
				CreateUnix int64  `json:"create_time"`
			} `json:"_source"`
			Sort []int64 `json:"sort"`
		} `json:"hits"`
	} `json:"hits"`
}
