package db

import (
	"context"
	"discord-go/config"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client     *mongo.Client
	clientOnce sync.Once
)

func Context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func Ping() time.Duration {
	ctx, cancel := Context()
	defer cancel()

	start := time.Now()
	Client.Ping(ctx, nil)

	duration := time.Since(start)
	return duration
}

func Connect() *mongo.Client {
	cfg := config.Loader()

	clientOnce.Do(func() {
		ctx, cancel := Context()
		defer cancel()

		opts := options.Client().ApplyURI(cfg.Section("database").Key("uri").String())

		client, err := mongo.Connect(ctx, opts)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}

		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("MongoDB connection failed: %v", err)
		}

		fmt.Println("Successfully connected to MongoDB!")

		Client = client
	})

	return Client
}
