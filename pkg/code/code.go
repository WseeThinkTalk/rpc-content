package code

import "rpc-content/pkg/xcode"

var (
	// Common
	ServerErr  = xcode.ServerErr
	NotFound   = xcode.NotFound
	RequestErr = xcode.RequestErr

	// Article (60000+)
	SortTypeInvalid         = xcode.New(60001, "排序类型无效")
	UserIdInvalid           = xcode.New(60002, "用户ID无效")
	ArticleTitleCantEmpty   = xcode.New(60003, "文章标题不能为空")
	ArticleContentCantEmpty = xcode.New(60004, "文章内容不能为空")
	ArticleIdInvalid        = xcode.New(60005, "文章ID无效")
	ArticleTitleTooLong     = xcode.New(60006, "文章标题过长")
	ArticleContentTooLong   = xcode.New(60007, "文章内容过长")
	ArticleDescTooLong           = xcode.New(60008, "文章描述过长")
	ContentContainsSensitiveWord = xcode.New(60009, "内容包含违规敏感词汇")

	// QA (90000+)
	QAUserIdInvalid  = xcode.New(90001, "用户ID无效")
	TitleEmpty       = xcode.New(90002, "问题标题不能为空")
	ContentEmpty     = xcode.New(90003, "内容不能为空")
	QuestionNotFound = xcode.New(90004, "问题不存在")
	AnswerNotFound   = xcode.New(90005, "回答不存在")
	NotQuestionAuthor = xcode.New(90006, "仅提问者可操作")
	NotAnswerAuthor   = xcode.New(90007, "仅回答者可操作")
	AlreadyAccepted  = xcode.New(90008, "已有采纳回答")
	QuestionIdEmpty  = xcode.New(90009, "问题ID不能为空")

	// Tag (50000+)
	TagNameEmpty      = xcode.New(50001, "标签名不能为空")
	TagNameTooLong    = xcode.New(50002, "标签名过长")
	TagNameExists     = xcode.New(50003, "标签名已存在")
	TagNotFound       = xcode.New(50004, "标签不存在")
	TagResourceExists = xcode.New(50005, "该资源已关联此标签")
	BizIdEmpty        = xcode.New(50006, "业务ID不能为空")
	TargetIdEmpty     = xcode.New(50007, "资源ID不能为空")
	TagUserIdEmpty    = xcode.New(50008, "用户ID不能为空")
	TagIdEmpty        = xcode.New(50009, "标签ID不能为空")
)
