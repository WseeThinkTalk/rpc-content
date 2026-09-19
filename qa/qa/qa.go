package qa

import (
	"context"

	"rpc-content/qa/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	PublishQuestionRequest  = pb.PublishQuestionRequest
	PublishQuestionResponse = pb.PublishQuestionResponse
	AnswerQuestionRequest   = pb.AnswerQuestionRequest
	AnswerQuestionResponse  = pb.AnswerQuestionResponse
	AcceptAnswerRequest     = pb.AcceptAnswerRequest
	AcceptAnswerResponse    = pb.AcceptAnswerResponse
	QuestionsRequest        = pb.QuestionsRequest
	QuestionsResponse       = pb.QuestionsResponse
	QuestionDetailRequest   = pb.QuestionDetailRequest
	QuestionDetailResponse  = pb.QuestionDetailResponse
	QuestionDeleteRequest   = pb.QuestionDeleteRequest
	QuestionDeleteResponse  = pb.QuestionDeleteResponse
	AnswerListRequest       = pb.AnswerListRequest
	AnswerListResponse      = pb.AnswerListResponse
	AnswerDeleteRequest     = pb.AnswerDeleteRequest
	AnswerDeleteResponse    = pb.AnswerDeleteResponse
	SearchQuestionsRequest  = pb.SearchQuestionsRequest
	SearchQuestionsResponse = pb.SearchQuestionsResponse
	QuestionItem            = pb.QuestionItem
	AnswerItem              = pb.AnswerItem
	SearchQuestionItem      = pb.SearchQuestionItem

	QA interface {
		PublishQuestion(ctx context.Context, in *PublishQuestionRequest, opts ...grpc.CallOption) (*PublishQuestionResponse, error)
		AnswerQuestion(ctx context.Context, in *AnswerQuestionRequest, opts ...grpc.CallOption) (*AnswerQuestionResponse, error)
		AcceptAnswer(ctx context.Context, in *AcceptAnswerRequest, opts ...grpc.CallOption) (*AcceptAnswerResponse, error)
		Questions(ctx context.Context, in *QuestionsRequest, opts ...grpc.CallOption) (*QuestionsResponse, error)
		QuestionDetail(ctx context.Context, in *QuestionDetailRequest, opts ...grpc.CallOption) (*QuestionDetailResponse, error)
		QuestionDelete(ctx context.Context, in *QuestionDeleteRequest, opts ...grpc.CallOption) (*QuestionDeleteResponse, error)
		AnswerList(ctx context.Context, in *AnswerListRequest, opts ...grpc.CallOption) (*AnswerListResponse, error)
		AnswerDelete(ctx context.Context, in *AnswerDeleteRequest, opts ...grpc.CallOption) (*AnswerDeleteResponse, error)
		SearchQuestions(ctx context.Context, in *SearchQuestionsRequest, opts ...grpc.CallOption) (*SearchQuestionsResponse, error)
	}

	defaultQA struct {
		cli zrpc.Client
	}
)

func NewQA(cli zrpc.Client) QA {
	return &defaultQA{cli: cli}
}

func (m *defaultQA) PublishQuestion(ctx context.Context, in *PublishQuestionRequest, opts ...grpc.CallOption) (*PublishQuestionResponse, error) {
	client := pb.NewQAClient(m.cli.Conn())
	return client.PublishQuestion(ctx, in, opts...)
}

func (m *defaultQA) AnswerQuestion(ctx context.Context, in *AnswerQuestionRequest, opts ...grpc.CallOption) (*AnswerQuestionResponse, error) {
	client := pb.NewQAClient(m.cli.Conn())
	return client.AnswerQuestion(ctx, in, opts...)
}

func (m *defaultQA) AcceptAnswer(ctx context.Context, in *AcceptAnswerRequest, opts ...grpc.CallOption) (*AcceptAnswerResponse, error) {
	client := pb.NewQAClient(m.cli.Conn())
	return client.AcceptAnswer(ctx, in, opts...)
}

func (m *defaultQA) Questions(ctx context.Context, in *QuestionsRequest, opts ...grpc.CallOption) (*QuestionsResponse, error) {
	client := pb.NewQAClient(m.cli.Conn())
	return client.Questions(ctx, in, opts...)
}

func (m *defaultQA) QuestionDetail(ctx context.Context, in *QuestionDetailRequest, opts ...grpc.CallOption) (*QuestionDetailResponse, error) {
	client := pb.NewQAClient(m.cli.Conn())
	return client.QuestionDetail(ctx, in, opts...)
}

func (m *defaultQA) QuestionDelete(ctx context.Context, in *QuestionDeleteRequest, opts ...grpc.CallOption) (*QuestionDeleteResponse, error) {
	client := pb.NewQAClient(m.cli.Conn())
	return client.QuestionDelete(ctx, in, opts...)
}

func (m *defaultQA) AnswerList(ctx context.Context, in *AnswerListRequest, opts ...grpc.CallOption) (*AnswerListResponse, error) {
	client := pb.NewQAClient(m.cli.Conn())
	return client.AnswerList(ctx, in, opts...)
}

func (m *defaultQA) AnswerDelete(ctx context.Context, in *AnswerDeleteRequest, opts ...grpc.CallOption) (*AnswerDeleteResponse, error) {
	client := pb.NewQAClient(m.cli.Conn())
	return client.AnswerDelete(ctx, in, opts...)
}

func (m *defaultQA) SearchQuestions(ctx context.Context, in *SearchQuestionsRequest, opts ...grpc.CallOption) (*SearchQuestionsResponse, error) {
	client := pb.NewQAClient(m.cli.Conn())
	return client.SearchQuestions(ctx, in, opts...)
}
