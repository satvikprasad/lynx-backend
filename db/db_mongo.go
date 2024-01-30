package db

import (
	"context"
	"fmt"
	"lynx-backend/models"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
}

func CreateMongoDB() (DB, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, fmt.Errorf("$MONGODB_URI must be set")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &MongoDB{
		client: client,
	}, nil
}

func (db *MongoDB) CreateMetric(m *models.Metric) error {
	collection := db.client.Database("lynx").Collection("metrics")

	if _, err := collection.InsertOne(context.TODO(), m); err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) Metrics() ([]models.Metric, error) {
	collection := db.client.Database("lynx").Collection("metrics")

	filter := bson.M{"time": bson.M{"$exists": true}}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var metrics []models.Metric
	if err = cur.All(context.TODO(), &metrics); err != nil {
		return nil, err
	}

	return metrics, nil
}

func (db *MongoDB) Disconnect() error {
	if err := db.client.Disconnect(context.TODO()); err != nil {
		return err
	}
	return nil
}
