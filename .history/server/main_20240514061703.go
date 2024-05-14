package main

import (
	"context"
	"log"
	"net"
	"strings"
  "fmt"	

	"github.com/jonathan-dotcom/asset-portfolio-management/server/asset"
	"github.com/jonathan-dotcom/asset-portfolio-management/server/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"

	"github.com/piquette/finance-go/crypto"
	"github.com/piquette/finance-go/quote"
)

const (
  port = ":50051"
  mongoURI = "mongodb://localhost:27017"
)

type server struct {
  asset.UnimplementedAssetServiceServer
  mongoClient *mongo.Client
}

func (s *server) CreateAsset(ctx context.Context, req *asset.CreateAssetRequest) (*asset.Asset, error) {
  assetCollection := s.mongoClient.Database("assetdb").Collection("assets")
  res, err := assetCollection.InsertOne(ctx, bson.M{
    "symbol":   req.Symbol,
    "quantity": req.Quantity,
    "price":    req.Price,
  })
  if err != nil {
    return nil, err
  }
  id := res.InsertedID.(primitive.ObjectID).Hex()
  return &asset.Asset{
    Id:       id,
    Symbol:   req.Symbol,
    Quantity: req.Quantity,
    Price:    req.Price,
  }, nil
}

func (s *server) GetAsset(ctx context.Context, req *asset.GetAssetRequest) (*asset.Asset, error) {
    assetCollection := s.mongoClient.Database("assetdb").Collection("assets")
    var result asset.Asset
    objID, err := primitive.ObjectIDFromHex(req.Id)
    if err != nil {
        return nil, err
    }
    err = assetCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
    if err != nil {
        return nil, err
    }
    result.Id = req.Id
    return &result, nil
}

func (s *server) UpdateAsset(ctx context.Context, req *asset.UpdateAssetRequest) (*asset.Asset, error) {
  assetCollection := s.mongoClient.Database("assetdb").Collection("assets")
  objID, err := primitive.ObjectIDFromHex(req.Id)
  if err != nil {
      return nil, err
  }

  // Retrieve the latest price based on the asset type
  var latestPrice float64
  if strings.Contains(req.Symbol, "-USD") {
      // Cryptocurrency pair
      pair := crypto.Get(req.Symbol)
      latestPrice = pair.RegularMarketPrice
      fmt.Printf("Current Price: $%v\n", pair.RegularMarketPrice)
  } else if strings.HasPrefix(req.Symbol, "$") {
      // Stock symbol
      quote, err := quote.Get(req.Symbol[1:]) // Remove the "$" prefix
      if err != nil {
          return nil, err
      }
      latestPrice = quote.RegularMarketPrice
      fmt.Printf("Current Price: $%v\n", quote.RegularMarketPrice)
  } else {
      return nil, fmt.Errorf("unsupported asset symbol format: %s", req.Symbol)
  }

  update := bson.M{
      "$set": bson.M{
          "symbol":   req.Symbol,
          "quantity": req.Quantity,
          "price":    latestPrice,
      },
  }
  _, err = assetCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
  if err != nil {
      return nil, err
  }
  return s.GetAsset(ctx, &asset.GetAssetRequest{Id: req.Id})
}

func (s *server) DeleteAsset(ctx context.Context, req *asset.DeleteAssetRequest) (*asset.Empty, error) {
    assetCollection := s.mongoClient.Database("assetdb").Collection("assets")
    objID, err := primitive.ObjectIDFromHex(req.Id)
    if err != nil {
        return nil, err
    }
    _, err = assetCollection.DeleteOne(ctx, bson.M{"_id": objID})
    if err != nil {
        return nil, err
    }
    return &asset.Empty{}, nil
}

func (s *server) ListAssets(ctx context.Context, _ *asset.Empty) (*asset.AssetList, error) {
  assetCollection := s.mongoClient.Database("assetdb").Collection("assets")
  cursor, err := assetCollection.Find(ctx, bson.M{})
  if err != nil {
      return nil, err
  }
  defer cursor.Close(ctx)
  var assets []*asset.Asset
  for cursor.Next(ctx) {
      var assetDB struct {
          ID       primitive.ObjectID `bson:"_id"`
          Symbol   string             `bson:"symbol"`
          Quantity int32              `bson:"quantity"`
          Price    float64            `bson:"price"`
      }
      err := cursor.Decode(&assetDB)
      if err != nil {
          return nil, err
      }
      asset := &asset.Asset{
          Id:       assetDB.ID.Hex(),
          Symbol:   assetDB.Symbol,
          Quantity: assetDB.Quantity,
          Price:    assetDB.Price,
      }
      assets = append(assets, asset)
  }
  return &asset.AssetList{Assets: assets}, nil
}

func main() {
  mongoClient, err := mongodb.NewClient(mongoURI)
  if err != nil {
    log.Fatalf("Failed to connect to MongoDB: %v", err)
  }
  defer mongoClient.Disconnect(context.Background())

  lis, err := net.Listen("tcp", port)
  if err != nil {
    log.Fatalf("Failed to listen: %v", err)
  }

  s := grpc.NewServer()
  asset.RegisterAssetServiceServer(s, &server{mongoClient: mongoClient})

  log.Printf("Server listening at %v", lis.Addr())
  if err := s.Serve(lis); err != nil {
    log.Fatalf("Failed to serve: %v", err)
  }
} 