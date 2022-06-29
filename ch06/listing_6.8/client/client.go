package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	order "listing_6.8"
	"log"
)

func getTlsCredentials() (credentials.TransportCredentials, error) {
	clientCert, clientCertErr := tls.LoadX509KeyPair("cert/client.crt", "cert/client.key")
	if clientCertErr != nil {
		return nil, fmt.Errorf("could not load client key pair : %v", clientCertErr)
	}

	certPool := x509.NewCertPool()
	caCert, caCertErr := ioutil.ReadFile("cert/ca.crt")
	if caCertErr != nil {
		return nil, fmt.Errorf("could not read Cert CA : %v", caCertErr)
	}

	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("failed to append CA cert.")
	}

	return credentials.NewTLS(&tls.Config{
		ServerName:   "*.microservices.dev",
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}), nil
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
