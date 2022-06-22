package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	order "listing_6.8"
	"log"
)

func getTlsCredentials() (credentials.TransportCredentials, error) {
	// Load certs from the d
	cert, err := tls.LoadX509KeyPair("cert/client.crt", "cert/client.key")
	if err != nil {
		return nil, fmt.Errorf("Could not load client key pair : %v", err)
	}

	// Create certpool from the CA
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("cert/ca.crt")
	if err != nil {
		return nil, fmt.Errorf("Could not read Cert CA : %v", err)
	}

	// Append the certs from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, fmt.Errorf("Failed to append CA cert : %v", err)
	}

	// Create transport creds based on TLS.
	creds := credentials.NewTLS(&tls.Config{
		ServerName:   "*.microservices.dev",
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	})
	return creds, nil
}

func main() {
	tlsCredentials, tlsCredentialsErr := getTlsCredentials()
	if tlsCredentialsErr != nil {
		log.Fatalf("failed to get tls credentials. %v", tlsCredentialsErr)
	}
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("Failed to connect order service. Err: %v", err)
	}

	defer conn.Close()

	orderClient := order.NewOrderServiceClient(conn)
	log.Println("Creating order...")
	orderResponse, errCreate := orderClient.Create(context.Background(), &order.CreateOrderRequest{
		UserId:    -1,
		ProductId: 0,
		Price:     2,
	})

	if errCreate != nil {
		log.Fatalf("Failed to create order. %v", errCreate)
	} else {
		log.Printf("Order %d is created successfully.", orderResponse.OrderId)
	}
}
