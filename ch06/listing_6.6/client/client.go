package main

import (
	"context"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	order "listing_6.6"
	"listing_6.6/middleware"
	"log"
	"time"
)

var cb *gobreaker.CircuitBreaker

func main() {
	cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "demo",
		MaxRequests: 3,
		Timeout:     4,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio >= 0.1
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Circuit Breaker: %s, changed from %v, to %v", name, from, to)
		},
	})
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithUnaryInterceptor(middleware.CircuitBreakerClientInterceptor(cb)))
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("Failed to connect order service. Err: %v", err)
	}

	defer conn.Close()

	orderClient := order.NewOrderServiceClient(conn)
	for {
		log.Println("Creating order...")
		orderResponse, errCreate := orderClient.Create(context.Background(), &order.CreateOrderRequest{
			UserId:    23,
			ProductId: 1234,
			Price:     3.2,
		})

		if errCreate != nil {
			log.Printf("Failed to create order. Err: %v", errCreate)
		} else {
			log.Printf("Order %d is created successfully.", orderResponse.OrderId)
		}
		time.Sleep(1 * time.Second)
	}

}
