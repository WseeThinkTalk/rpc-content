package server

import (
	"context"

	"rpc-content/qa/internal/logic"
	"rpc-content/qa/internal/svc"
	"rpc-content/qa/pb"
)

type QAServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedQAServer
}

func NewQAServer(svcCtx *svc.ServiceContext) *QAServer {
	return &QAServer{svcCtx: svcCtx}
}

func (s *QAServer) PublishQuestion(ctx context.Context, in *pb.PublishQuestionRequest) (*pb.PublishQuestionResponse, error) {
	l := logic.NewPublishQuestionLogic(ctx, s.svcCtx)
	return l.PublishQuestion(in)
}

func (s *QAServer) AnswerQuestion(ctx context.Context, in *pb.AnswerQuestionRequest) (*pb.AnswerQuestionResponse, error) {
	l := logic.NewAnswerQuestionLogic(ctx, s.svcCtx)
	return l.AnswerQuestion(in)
}

func (s *QAServer) AcceptAnswer(ctx context.Context, in *pb.AcceptAnswerRequest) (*pb.AcceptAnswerResponse, error) {
	l := logic.NewAcceptAnswerLogic(ctx, s.svcCtx)
	return l.AcceptAnswer(in)
}

func (s *QAServer) Questions(ctx context.Context, in *pb.QuestionsRequest) (*pb.QuestionsResponse, error) {
	l := logic.NewQuestionsLogic(ctx, s.svcCtx)
	return l.Questions(in)
}

func (s *QAServer) QuestionDetail(ctx context.Context, in *pb.QuestionDetailRequest) (*pb.QuestionDetailResponse, error) {
	l := logic.NewQuestionDetailLogic(ctx, s.svcCtx)
	return l.QuestionDetail(in)
}

func (s *QAServer) QuestionDelete(ctx context.Context, in *pb.QuestionDeleteRequest) (*pb.QuestionDeleteResponse, error) {
	l := logic.NewQuestionDeleteLogic(ctx, s.svcCtx)
	return l.QuestionDelete(in)
}

func (s *QAServer) AnswerList(ctx context.Context, in *pb.AnswerListRequest) (*pb.AnswerListResponse, error) {
	l := logic.NewAnswerListLogic(ctx, s.svcCtx)
	return l.AnswerList(in)
}

func (s *QAServer) AnswerDelete(ctx context.Context, in *pb.AnswerDeleteRequest) (*pb.AnswerDeleteResponse, error) {
	l := logic.NewAnswerDeleteLogic(ctx, s.svcCtx)
	return l.AnswerDelete(in)
}

func (s *QAServer) SearchQuestions(ctx context.Context, in *pb.SearchQuestionsRequest) (*pb.SearchQuestionsResponse, error) {
	l := logic.NewSearchQuestionsLogic(ctx, s.svcCtx)
	return l.SearchQuestions(in)
}
