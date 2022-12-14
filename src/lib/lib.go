package lib

import (
	"fmt"
	"log"
	"time"
)

// type RetryFunc func(any) (*amqp.Connection, error)

func Retry[T any](retryFunc func(string) (T, error), retries, delay int) func(string) (T, error) {
	return func(args string) (T, error) {
		for c := 1; ; c++ {
			resp, err := retryFunc(args)
			if err == nil {
				log.Println("Got Response.")
				return resp, nil
			}
			if c > retries {
				return resp, fmt.Errorf("not able to connect after %v retires", c)
			}
			delay *= 2
			log.Printf("Retry: %v, sleeping for %v", c, delay)
			time.Sleep(time.Duration(delay) * time.Second)
		}
	}
}
