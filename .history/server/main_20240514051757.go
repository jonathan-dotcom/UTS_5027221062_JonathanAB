package main

import (
	"context"
	"log"
	"net"
	"strings"

	"github.com/jonathan-dotcom/asset-portfolio-management/server/asset"
	"github.com/jonathan-dotcom/asset-portfolio-management/server/mongodb"
	"github.com/piquette/finance-go/crypto"
	"github.com/piquette/finance-go/quote"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

const (
  port = ":50051"
  mongoURI = "mongodb://localhost:27017"
)

type server struct {
  asset.UnimplementedAssetServiceServer
  mongoClient *mongo.Client
}

type AssetInfo struct {
  ID       string  `json:"id"`
  Symbol   string  `json:"symbol"`
  Quantity int32   `json:"quantity"`
  Price    float64 `json:"price"`
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
    update := bson.M{
        "$set": bson.M{
            "symbol":   req.Symbol,
            "quantity": req.Quantity,
            "price":    req.Price,
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

  var stockSymbols []string
  var cryptoSymbols []string
  assets := make([]*asset.Asset, 0)

  for cursor.Next(ctx) {
    var assetDB struct {
      ID       primitive.ObjectID `bson:"_id"`
      Symbol   string             `bson:"symbol"`
      Quantity int32              `bson:"quantity"`
    }
    err := cursor.Decode(&assetDB)
    if err != nil {
      log.Printf("Error decoding asset from database: %v", err)
      continue
    }
    if strings.HasPrefix(assetDB.Symbol, "USD") {
      cryptoSymbols = append(cryptoSymbols, assetDB.Symbol)
    } else {
      stockSymbols = append(stockSymbols, assetDB.Symbol)
    }
    assets = append(assets, &asset.Asset{
      Id:       assetDB.ID.Hex(),
      Symbol:   assetDB.Symbol,
      Quantity: assetDB.Quantity,
    })
  }

  // Fetch stock quotes
  stockPrices := make(map[string]float64)
  stockIter := quote.List(stockSymbols)
  for stockIter.Next() {
    q := stockIter.Quote()
    stockPrices[q.Symbol] = q.RegularMarketPrice
  }

  // Fetch cryptocurrency quotes
  cryptoPrices := make(map[string]float64)
  cryptoIter := crypto.List(cryptoSymbols)
  for cryptoIter.Next() {
    c := cryptoIter.CryptoPair()
    cryptoPrices[c.Symbol] = c.RegularMarketPrice
  }

  // Convert assets to AssetWithPrice struct and assign real-time prices
  assetsInfo := make([]*AssetInfo, len(assets))
  for i, asset := range assets {
    assetsInfo[i] = &AssetInfo{
      ID:       asset.Id,
      Symbol:   asset.Symbol,
      Quantity: asset.Quantity,
      Price:    0,
    }
    if price, ok := stockPrices[asset.Symbol]; ok {
      assetsInfo[i].Price = price
    } else if price, ok := cryptoPrices[asset.Symbol]; ok {
      assetsInfo[i].Price = price
    }
  }

  // Convert assetsInfo to []*asset.AssetInfo
  assetInfoList := make([]*asset.AssetInfo, len(assetsInfo))
  for i, info := range assetsInfo {
    assetInfoList[i] = &asset.AssetInfo{
      Id:       info.ID,
      Symbol:   info.Symbol,
      Quantity: info.Quantity,
      Price:    info.Price,
    }
  }

  if err := cursor.Err(); err != nil {
    log.Printf("Error iterating over assets: %v", err)
  }
  return &asset.AssetList{Assets: assets, AssetsInfo: assetInfoList}, nil
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