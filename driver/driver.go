package driver

import (
	"context"
	"github.com/yusuf/go-app/modules/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var app config.GoAppTools

func Connection(URI string) *mongo.Client {

	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		app.ErrorLogger.Panicln(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		app.ErrorLogger.Fatalln(err)
	}

	return client
}
