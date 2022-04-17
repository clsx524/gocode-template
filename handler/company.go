package handler

import (
	"context"
	"github.com/Rippling/gocode-template/client"
	"github.com/Rippling/gocode-template/model"
	pb "github.com/Rippling/gocode-template/rpc/company"
	"github.com/Rippling/gocode-template/service"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/durationpb"
)

type CompanyHandlerDeps struct {
	fx.In

	client.Instrumenter
	service.Company
}

type CompanyHandlerOutputs struct {
	fx.Out

	CompanyServer pb.TwirpServer `group:"handlers"`
}

type companyHandler struct {
	deps CompanyHandlerDeps
}

func (s *companyHandler) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	s.deps.Logger().Info("received search request", zap.String("req_name", req.GetName()))
	company, err := s.deps.Company.Search(ctx, req.GetName())
	if err != nil {
		return nil, err
	}

	if company == nil {
		return &pb.SearchResponse{
			Instances: nil,
		}, nil
	}

	var resp []*pb.Company
	resp = append(resp, &pb.Company{
		Id:               company.ID,
		Name:             company.Name,
		SinceLastUpdated: durationpb.New(100),
	})
	return &pb.SearchResponse{
		Instances: resp,
	}, nil
}

func (s *companyHandler) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	s.deps.Logger().Info("received add request")

	var resp []*model.Company
	for _, c := range req.Instances {
		resp = append(resp, &model.Company{
			ID:   c.GetId(),
			Name: c.GetName(),
		})
	}

	err := s.deps.Company.Add(ctx, resp)
	return &pb.AddResponse{
		Status: err == nil,
	}, nil
}

func NewCompanyHandler(deps CompanyHandlerDeps) CompanyHandlerOutputs {
	server := &companyHandler{deps}
	return CompanyHandlerOutputs{
		CompanyServer: pb.NewCompanyServiceServer(server),
	}
}
