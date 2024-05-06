package database

import (
  "context"
  "log"
  "time"

  "github.com/jonathan-dotcom/asset-portfolio-management/proto"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func Init(uri string) {
  client, err := mongo.NewClient(options.Client().ApplyURI(uri))
  if err != nil {
    log.Fatalf("failed to create MongoDB client: %v", err)
  }
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()
  err = client.Connect(ctx)
  if err != nil {
    log.Fatalf("failed to connect to MongoDB: %v", err)
  }
  db = client.Database("portfolio_management")
}

func CreateAsset(asset *proto.Asset) error {
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  collection := db.Collection("assets")
  result, err := collection.InsertOne(ctx, bson.M{
    "name":  asset.GetName(),
    "price": asset.GetPrice(),
  })
  if err != nil {
    return err
  }
  asset.Id = result.InsertedID.(primitive.ObjectID).Hex()
  return nil
}

// Implement other database functions (GetAsset, UpdateAsset, DeleteAsset)