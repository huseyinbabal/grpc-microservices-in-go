package main

import (
	"context"
	"fmt"

	payment "listing_2.2"

	"google.golang.org/grpc"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial("http:/localhost:8080", opts...)
	if err != nil {
		fmt.Println("It is fine, this is not a complete example.")
	}

	defer conn.Close()

	paymentClient := payment.NewPaymentServiceClient(conn)
	ctx := context.Background()
	result, err := paymentClient.Create(ctx, &payment.CreatePaymentRequest{Price: 23})
	if err != nil {
		fmt.Println("Don't worry, we don't expect to see it is working.")
	}
	fmt.Printf("Result %v", result)
}
