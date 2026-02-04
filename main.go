package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

func ConnectSQL() *pgxpool.Pool {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		"niflheim", "niflguard", "postgres_ds", "5432", "ds_db")

	poolconfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	ctx := context.Background()

	pool, err := pgxpool.NewWithConfig(ctx, poolconfig)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	if err = pool.Ping(ctx); err != nil {
		log.Fatal(err)
		panic(err)
	}
	return pool
}
func InitRoute() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))

	return mux
}
func InitKafka() *kafka.Writer {
	// To connect to Kafka, we use the external listener you configured.
	// The address is the public IP of your server and the external port.
	kafkaURL := "194.147.142.26:9094"
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    "my-topic", 	
		Balancer: &kafka.LeastBytes{},
	}
	log.Println("Kafka writer initialized")
	return writer
}

func main() {
	db := ConnectSQL()
	kafkaWriter := InitKafka()
	defer kafkaWriter.Close()

	server := &http.Server{
		Addr:    ":8081",
		Handler: InitRoute(),
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
