package client

import (
  "context"

  "github.com/jonathan-dotcom/asset-portfolio-management/proto"
  "google.golang.org/grpc"
)

type PortfolioClient struct {
  client proto.PortfolioServiceClient
}

func NewPortfolioClient(addr string) (*PortfolioClient, error) {
  conn, err := grpc.Dial(addr, grpc.WithInsecure())
  if err != nil {
    return nil, err
  }
  client := proto.NewPortfolioServiceClient(conn)
  return &PortfolioClient{client}, nil
}

func (c *PortfolioClient) CreateAsset(name string, price float64) (*proto.Asset, error) {
  req := &proto.CreateAssetRequest{
    Name:  name,
    Price: price,
  }
  return c.client.CreateAsset(context.Background(), req)
}

// Implement other client functions (GetAsset, UpdateAsset, DeleteAsset)