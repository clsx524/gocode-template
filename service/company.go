package service

import (
	"context"
	"github.com/Rippling/gocode-template/model"
	"github.com/Rippling/gocode-template/repository"
	"go.uber.org/fx"
)

type Dep struct {
	fx.In

	repository.Company
}

type Company interface {
	Search(ctx context.Context, name string) (*model.Company, error)
	Add(ctx context.Context, companies []*model.Company) error
}

type company struct {
	deps Dep
}

func (c *company) Search(ctx context.Context, name string) (*model.Company, error) {
	return c.deps.Company.Search(ctx, name)
}

func (c *company) Add(ctx context.Context, companies []*model.Company) error {
	return c.deps.AddAll(ctx, companies)
}

func New(deps Dep) Company {
	return &company{deps}
}
