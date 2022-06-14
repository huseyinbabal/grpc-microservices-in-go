package main

import (
	"context"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	order "listing_6.7"
	"log"
	"net"
)

type server struct {
	order.UnimplementedOrderServiceServer
}

func (s *server) Create(ctx context.Context, in *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	var validationErrors []*errdetails.BadRequest_FieldViolation
	if in.UserId < 1 {
		validationErrors = append(validationErrors, &errdetails.BadRequest_FieldViolation{
			Field:       "user_id",
			Description: "user id cannot be less than 1",
		})
	}
	if in.ProductId < 0 {
		validationErrors = append(validationErrors, &errdetails.BadRequest_FieldViolation{
			Field:       "product_id",
			Description: "product id cannot be negative",
		})
	}
	if len(validationErrors) > 0 {
		stat := status.New(400, "invalid order request")
		badRequest := &errdetails.BadRequest{}
		badRequest.FieldViolations = validationErrors
		s, _ := stat.WithDetails(badRequest)
		return nil, s.Err()
	}
	return &order.CreateOrderResponse{OrderId: 1243}, nil
}

func main() {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	order.RegisterOrderServiceServer(grpcServer, &server{})
	grpcServer.Serve(listener)
}
