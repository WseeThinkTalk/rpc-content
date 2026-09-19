package code

import "rpc-content/pkg/xcode"

var (
	TagNameEmpty       = xcode.New(50001, "标签名不能为空")
	TagNameTooLong     = xcode.New(50002, "标签名过长")
	TagNameExists      = xcode.New(50003, "标签名已存在")
	TagNotFound        = xcode.New(50004, "标签不存在")
	TagResourceExists  = xcode.New(50005, "该资源已关联此标签")
	BizIdEmpty         = xcode.New(50006, "业务ID不能为空")
	TargetIdEmpty      = xcode.New(50007, "资源ID不能为空")
	UserIdEmpty        = xcode.New(50008, "用户ID不能为空")
	TagIdEmpty         = xcode.New(50009, "标签ID不能为空")
)
