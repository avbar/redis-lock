package main

import (
	"log"
	"net/http"

	"github.com/avbar/redis-lock/internal/handler"
	"github.com/avbar/redis-lock/internal/redis/locker"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	locker := locker.NewRedisLocker(client)
	handler := handler.NewHandler(locker)

	r := mux.NewRouter()

	r.HandleFunc("/lock/{key}", handler.Lock).Methods("POST")
	r.HandleFunc("/unlock/{key}", handler.Unlock).Methods("POST")

	log.Println("server is listening at 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
