package main

import (
	"errors"
	"github.com/sony/gobreaker"
	"log"
	"math/rand"
)

var cb *gobreaker.CircuitBreaker

func main() {
	cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "demo",
		MaxRequests: 3,
		Timeout:     4,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Circuit Breaker: %s, changed from %v, to %v", name, from, to)
		},
	})
	cbRes, cbErr := cb.Execute(func() (interface{}, error) {
		res, isErr := isError()
		if isErr {
			return nil, errors.New("error")
		}
		return res, nil
	})
	if cbErr != nil {
		log.Fatalf("Circuit breaker error %v", cbErr)
	} else {
		log.Printf("Circuit breaker result %v", cbRes)
	}
}

func isError() (int, bool) {
	min := 10
	max := 30
	result := rand.Intn(max-min) + min
	return result, result != 20
}
