package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	cert, err := tls.LoadX509KeyPair("redis.crt", "redis_private.key")
	if err != nil {
		log.Fatal(err)
	}
	// Load CA cert
	caCert, err := os.ReadFile("redis_ca.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Username: "default",
		Password: "redisPassword",
		TLSConfig: &tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		},
	})
	ctx := context.Background()
	session := map[string]string{"name": "XYZ", "surname": "ABCD", "DOB": "31-01-2001"}
	for k, v := range session {
		err := client.HSet(ctx, "user-1", k, v).Err()
		if err != nil {
			panic(err)
		}
	}
	userSession := client.HGetAll(ctx, "user-1").Val()
	fmt.Println(userSession)
}
