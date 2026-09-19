package code

import "rpc-content/pkg/xcode"

var (
	UserIdInvalid        = xcode.New(90001, "用户ID无效")
	TitleEmpty           = xcode.New(90002, "问题标题不能为空")
	ContentEmpty         = xcode.New(90003, "内容不能为空")
	QuestionNotFound     = xcode.New(90004, "问题不存在")
	AnswerNotFound       = xcode.New(90005, "回答不存在")
	NotQuestionAuthor    = xcode.New(90006, "仅提问者可操作")
	NotAnswerAuthor      = xcode.New(90007, "仅回答者可操作")
	AlreadyAccepted      = xcode.New(90008, "已有采纳回答")
	QuestionIdEmpty      = xcode.New(90009, "问题ID不能为空")
)
