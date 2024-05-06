package server

import (
  "context"

  "github.com/jonathan-dotcom/asset-portfolio-management/database"
  "github.com/jonathan-dotcom/asset-portfolio-management/proto"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
  "google.golang.org/protobuf/types/known/emptypb"
)

type portfolioServer struct {
  proto.UnimplementedPortfolioServiceServer
}

func NewPortfolioServer() proto.PortfolioServiceServer {
  return &portfolioServer{}
}

func (s *portfolioServer) CreateAsset(ctx context.Context, req *proto.CreateAssetRequest) (*proto.Asset, error) {
  asset := &proto.Asset{
    Name:  req.GetName(),
    Price: req.GetPrice(),
  }

  err := database.CreateAsset(asset)
  if err != nil {
    return nil, status.Errorf(codes.Internal, "failed to create asset: %v", err)
  }

  return asset, nil
}
