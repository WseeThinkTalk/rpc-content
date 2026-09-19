package code

import (
	"rpc-content/pkg/xcode"
)

var (
	SortTypeInvalid         = xcode.New(60001, "排序类型无效")     // 排序类型无效
	UserIdInvalid           = xcode.New(60002, "用户ID无效")     // 用户ID无效
	ArticleTitleCantEmpty   = xcode.New(60003, "文章标题不能为空")   // 文章标题不能为空
	ArticleContentCantEmpty = xcode.New(60004, "文章内容不能为空")   // 文章内容不能为空
	ArticleIdInvalid        = xcode.New(60005, "文章ID无效")     // 文章ID无效
	ArticleTitleTooLong     = xcode.New(60006, "文章标题过长")     // 文章标题过长
	ArticleContentTooLong   = xcode.New(60007, "文章内容过长")     // 文章内容过长
	ArticleDescTooLong      = xcode.New(60008, "文章描述过长")     // 文章描述过长
)
