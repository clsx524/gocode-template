package client

import (
	"context"
	"github.com/Rippling/gocode-template/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"os"
)

const mongoUriEnvKey = "MONGODB_URI"

type Mongo interface {
	GetMongoClient(ctx context.Context) *mongo.Client
}

type MongoDeps struct {
	fx.In

	Instrumenter
	config.Provider

	Lifecycle fx.Lifecycle
}

type mongoClient struct {
	deps MongoDeps

	client *mongo.Client
}

func (i *mongoClient) GetMongoClient(ctx context.Context) *mongo.Client {
	return i.client
}

func ProvideMongoClient(deps MongoDeps) (Mongo, error) {
	if err := os.Setenv(mongoUriEnvKey, deps.MongoConfig().URI); err != nil {
		return nil, err
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(deps.MongoConfig().URI))
	if err != nil {
		return nil, err
	}

	deps.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Disconnect(ctx)
		},
	})

	return &mongoClient{deps: deps, client: client}, nil
}
