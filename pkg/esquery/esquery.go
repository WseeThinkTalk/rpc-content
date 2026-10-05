package esquery

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// SearchParams 搜索请求参数
type SearchParams struct {
	Keyword   string
	AuthorIds []int64
	Status    int
	From      int
	Size      int
}

// BuildArticleSearchIDOnlyDSL 构造只返回 ID 的轻量 ES DSL
func BuildArticleSearchIDOnlyDSL(params SearchParams) (string, error) {
	var mustClauses []map[string]interface{}
	var filterClauses []map[string]interface{}

	// 状态过滤（走 filter 享受 BitSet 缓存）
	if params.Status > 0 {
		filterClauses = append(filterClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"status": params.Status,
			},
		})
	}

	// 作者 ID 过滤
	if len(params.AuthorIds) > 0 {
		filterClauses = append(filterClauses, map[string]interface{}{
			"terms": map[string]interface{}{
				"author_id": params.AuthorIds,
			},
		})
	}

	// 关键词分词与权重倾斜
	if params.Keyword != "" {
		mustClauses = append(mustClauses, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  params.Keyword,
				"fields": []string{"title^3", "description^1"},
				"type":   "best_fields",
			},
		})
	} else {
		mustClauses = append(mustClauses, map[string]interface{}{
			"match_all": map[string]interface{}{},
		})
	}

	query := map[string]interface{}{
		"bool": map[string]interface{}{
			"must":   mustClauses,
			"filter": filterClauses,
		},
	}

	if params.Size <= 0 {
		params.Size = 20
	}

	searchBody := map[string]interface{}{
		"from":    params.From,
		"size":    params.Size,
		"query":   query,
		"_source": []string{"id"}, // 两阶段检索核心：仅返回 ID
		"sort": []map[string]interface{}{
			{"_score": map[string]interface{}{"order": "desc"}},
			{"publish_time": map[string]interface{}{"order": "desc"}},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchBody); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ParseSearchIDs 从 ES 搜索结果中提取有序的文章 ID 列表
func ParseSearchIDs(responseBody []byte) (ids []int64, total int64, err error) {
	var res struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				ID     string `json:"_id"`
				Source struct {
					ID int64 `json:"id"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.Unmarshal(responseBody, &res); err != nil {
		return nil, 0, err
	}

	total = res.Hits.Total.Value
	ids = make([]int64, 0, len(res.Hits.Hits))
	for _, hit := range res.Hits.Hits {
		if hit.Source.ID > 0 {
			ids = append(ids, hit.Source.ID)
		} else {
			if id, err := strconv.ParseInt(hit.ID, 10, 64); err == nil {
				ids = append(ids, id)
			}
		}
	}

	return ids, total, nil
}
