package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"

	pb "github.com/booler007/gRPC_currentrate/pb"
	"github.com/booler007/gRPC_currentrate/storage"
)

const apiURL = "https://grinex.io/api/v1/spot/depth?symbol=usdta7a5"

type ratesServer struct {
	pb.UnimplementedRatesServiceServer
	repo storage.Repository
}

func NewRatesServer(repo storage.Repository) *ratesServer {
	return &ratesServer{repo: repo}
}

type orderBookEntry struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
}

type depthResponse struct {
	Timestamp int64            `json:"timestamp"`
	Asks      []orderBookEntry `json:"asks"`
	Bids      []orderBookEntry `json:"bids"`
}

func (s *ratesServer) HealthCheck(_ context.Context, _ *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Status: "ok"}, nil
}

func (s *ratesServer) GetRates(_ context.Context, _ *pb.GetRatesRequest) (*pb.GetRatesResponse, error) {
	client := resty.New()
	resp, err := client.R().Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}

	var book depthResponse
	if err := json.Unmarshal(resp.Body(), &book); err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}

	if len(book.Asks) == 0 || len(book.Bids) == 0 {
		return nil, fmt.Errorf("empty asks or bids in response")
	}

	bestAsk := book.Asks[0].Price
	bestBid := book.Bids[0].Price

	id, err := s.repo.InsertRate(book.Timestamp, bestAsk, bestBid)
	if err != nil {
		return nil, fmt.Errorf("db insert failed: %w", err)
	}

	return &pb.GetRatesResponse{
		Id:        id,
		Timestamp: book.Timestamp,
		BestAsk:   bestAsk,
		BestBid:   bestBid,
	}, nil
}
