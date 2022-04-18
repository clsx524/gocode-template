package repositories

import (
	"context"
	"github.com/Rippling/gocode-template/clients"
	"github.com/Rippling/gocode-template/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	database   = "sample"
	collection = "company"
)

// CompanyRepoDeps defines all dependencies used in Company repository
type CompanyRepoDeps struct {
	fx.In

	clients.Instrumenter
	clients.Mongo
}

type Company interface {
	Search(ctx context.Context, name string) (*models.Company, error)
	AddAll(ctx context.Context, companies []*models.Company) error
}

type company struct {
	deps CompanyRepoDeps
}

func (m *company) Search(ctx context.Context, name string) (*models.Company, error) {
	coll := m.deps.GetMongoClient(ctx).Database(database).Collection(collection)

	var result models.Company
	err := coll.FindOne(ctx, bson.D{{"name", name}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		m.deps.Logger().Warn("No document was found with the name", zap.String("name", name))
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (m *company) AddAll(ctx context.Context, companies []*models.Company) error {
	coll := m.deps.GetMongoClient(ctx).Database(database).Collection(collection)

	var many []interface{}
	for _, company := range companies {
		many = append(many, &bson.D{
			{Key: "name", Value: company.Name},
			{Key: "id", Value: company.ID},
		})
	}
	_, err := coll.InsertMany(ctx, many)
	return err
}

// New return an instance of Company repository
func New(deps CompanyRepoDeps) Company {
	return &company{deps}
}
