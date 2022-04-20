package apis

import (
	"context"
	"github.com/clsx524/gocode-template/clients"
	"github.com/clsx524/gocode-template/models"
	pb "github.com/clsx524/gocode-template/rpc/company"
	"github.com/clsx524/gocode-template/services"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/durationpb"
)

type CompanyHandlerDeps struct {
	fx.In

	clients.Instrumenter
	services.Company
}

type CompanyHandlerOutputs struct {
	fx.Out

	CompanyServer pb.TwirpServer `group:"apis"`
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

	var resp []*models.Company
	for _, c := range req.Instances {
		resp = append(resp, &models.Company{
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
