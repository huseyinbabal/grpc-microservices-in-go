package main

import (
	"context"
	payment "listing_2.2"
	"log"

	"google.golang.org/grpc"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial("http:/localhost:8080", opts...)
	if err != nil {
		log.Println("It is fine, this is not a complete example.")
	}

	defer conn.Close()

	paymentClient := payment.NewPaymentServiceClient(conn)
	ctx := context.Background()
	_, err = paymentClient.Create(ctx, &payment.CreatePaymentRequest{Price: 23})
	if err != nil {
		log.Println("Don't worry, we don't expect to see it is working.")
	}
}
