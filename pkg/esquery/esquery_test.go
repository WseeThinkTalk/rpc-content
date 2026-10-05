package esquery

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBuildArticleSearchIDOnlyDSL(t *testing.T) {
	params := SearchParams{
		Keyword:   "露营穿搭",
		AuthorIds: []int64{101, 102},
		Status:    2,
		From:      0,
		Size:      20,
	}

	dsl, err := BuildArticleSearchIDOnlyDSL(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var dslMap map[string]interface{}
	if err := json.Unmarshal([]byte(dsl), &dslMap); err != nil {
		t.Fatalf("invalid json generated: %v", err)
	}

	// 验证只返回 id 字段（两阶段检索核心）
	source, ok := dslMap["_source"].([]interface{})
	if !ok || len(source) != 1 || source[0] != "id" {
		t.Fatalf("expected _source to be [id], got: %v", dslMap["_source"])
	}

	if !strings.Contains(dsl, "露营穿搭") {
		t.Fatalf("dsl missing keyword")
	}
	if !strings.Contains(dsl, "multi_match") {
		t.Fatalf("dsl missing multi_match")
	}
}

func TestParseSearchIDs(t *testing.T) {
	mockESResponse := `{
		"hits": {
			"total": {"value": 2},
			"hits": [
				{"_id": "1001", "_source": {"id": 1001}},
				{"_id": "1002", "_source": {"id": 1002}}
			]
		}
	}`

	ids, total, err := ParseSearchIDs([]byte(mockESResponse))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}
	if len(ids) != 2 || ids[0] != 1001 || ids[1] != 1002 {
		t.Fatalf("expected ids [1001, 1002], got: %v", ids)
	}
}
