package cacheguard

import "encoding/json"

// NullArticleID 空对象哨兵 ID，用于标识数据库中不存在的记录
const NullArticleID int64 = -1

// NullArticleTTL 防穿透空缓存过期时间（秒）
const NullArticleTTL = 60

// NormalArticleTTL 正常文章详情缓存过期时间（秒，2小时）
const NormalArticleTTL = 7200

// IsNullArticle 判断文章 ID 是否为防穿透空对象哨兵
func IsNullArticle(id int64) bool {
	return id == NullArticleID
}

// BuildNullArticlePayload 生成防穿透空对象 JSON 载荷
func BuildNullArticlePayload() string {
	payload, _ := json.Marshal(map[string]interface{}{
		"id": NullArticleID,
	})
	return string(payload)
}
