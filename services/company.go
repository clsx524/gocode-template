package services

import (
	"context"
	client "github.com/clsx524/gocode-template/clients"
	"github.com/clsx524/gocode-template/models"
	"github.com/clsx524/gocode-template/repositories"
	"go.uber.org/fx"
)

// CompanySvcDeps defines the dependencies used in Company service
type CompanySvcDeps struct {
	fx.In

	client.Instrumenter
	repositories.Company
}

type Company interface {
	Search(ctx context.Context, name string) (*models.Company, error)
	Add(ctx context.Context, companies []*models.Company) error
}

type company struct {
	deps CompanySvcDeps
}

func (c *company) Search(ctx context.Context, name string) (*models.Company, error) {
	return c.deps.Company.Search(ctx, name)
}

func (c *company) Add(ctx context.Context, companies []*models.Company) error {
	return c.deps.AddAll(ctx, companies)
}

// New returns an instance of Company service
func New(deps CompanySvcDeps) Company {
	return &company{deps}
}
