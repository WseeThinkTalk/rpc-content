package server

import (
	"context"

	"rpc-content/tag/internal/logic"
	"rpc-content/tag/internal/svc"
	"rpc-content/tag/pb"
)

type TagServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedTagServer
}

func NewTagServer(svcCtx *svc.ServiceContext) *TagServer {
	return &TagServer{
		svcCtx: svcCtx,
	}
}

// 创建标签
func (s *TagServer) CreateTag(ctx context.Context, in *pb.CreateTagRequest) (*pb.CreateTagResponse, error) {
	l := logic.NewCreateTagLogic(ctx, s.svcCtx)
	return l.CreateTag(in)
}

// 更新标签
func (s *TagServer) UpdateTag(ctx context.Context, in *pb.UpdateTagRequest) (*pb.UpdateTagResponse, error) {
	l := logic.NewUpdateTagLogic(ctx, s.svcCtx)
	return l.UpdateTag(in)
}

// 删除标签
func (s *TagServer) DeleteTag(ctx context.Context, in *pb.DeleteTagRequest) (*pb.DeleteTagResponse, error) {
	l := logic.NewDeleteTagLogic(ctx, s.svcCtx)
	return l.DeleteTag(in)
}

// 标签详情
func (s *TagServer) TagDetail(ctx context.Context, in *pb.TagDetailRequest) (*pb.TagDetailResponse, error) {
	l := logic.NewTagDetailLogic(ctx, s.svcCtx)
	return l.TagDetail(in)
}

// 标签列表
func (s *TagServer) TagList(ctx context.Context, in *pb.TagListRequest) (*pb.TagListResponse, error) {
	l := logic.NewTagListLogic(ctx, s.svcCtx)
	return l.TagList(in)
}

// 热门标签
func (s *TagServer) HotTags(ctx context.Context, in *pb.HotTagsRequest) (*pb.HotTagsResponse, error) {
	l := logic.NewHotTagsLogic(ctx, s.svcCtx)
	return l.HotTags(in)
}

// 资源打标签
func (s *TagServer) TagResource(ctx context.Context, in *pb.TagResourceRequest) (*pb.TagResourceResponse, error) {
	l := logic.NewTagResourceLogic(ctx, s.svcCtx)
	return l.TagResource(in)
}

// 资源去标签
func (s *TagServer) UntagResource(ctx context.Context, in *pb.UntagResourceRequest) (*pb.UntagResourceResponse, error) {
	l := logic.NewUntagResourceLogic(ctx, s.svcCtx)
	return l.UntagResource(in)
}

// 标签下的资源列表
func (s *TagServer) ResourcesByTag(ctx context.Context, in *pb.ResourcesByTagRequest) (*pb.ResourcesByTagResponse, error) {
	l := logic.NewResourcesByTagLogic(ctx, s.svcCtx)
	return l.ResourcesByTag(in)
}

// 资源关联的标签列表
func (s *TagServer) TagsByResource(ctx context.Context, in *pb.TagsByResourceRequest) (*pb.TagsByResourceResponse, error) {
	l := logic.NewTagsByResourceLogic(ctx, s.svcCtx)
	return l.TagsByResource(in)
}
