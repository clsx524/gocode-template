package repository

import (
	"context"
	"github.com/Rippling/gocode-template/client"
	"github.com/Rippling/gocode-template/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	database   = "sample"
	collection = "company"
)

type Deps struct {
	fx.In

	client.Instrumenter
	client.Mongo
}

type Company interface {
	Search(ctx context.Context, name string) (*model.Company, error)
	AddAll(ctx context.Context, companies []*model.Company) error
}

type company struct {
	deps Deps
}

func (m *company) Search(ctx context.Context, name string) (*model.Company, error) {
	coll := m.deps.GetMongoClient(ctx).Database(database).Collection(collection)

	var result model.Company
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

func (m *company) AddAll(ctx context.Context, companies []*model.Company) error {
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

func New(deps Deps) Company {
	return &company{deps}
}
