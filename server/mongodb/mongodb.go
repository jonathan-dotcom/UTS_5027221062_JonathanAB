package mongodb

import (
  "context"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(uri string) (*mongo.Client, error) {
  clientOptions := options.Client().ApplyURI(uri)
  client, err := mongo.Connect(context.Background(), clientOptions)
  if err != nil {
    return nil, err
  }
  return client, nil
}