package main

import (
	"log"
	"net/http"

	"github.com/exporter/redis_exporter/controllera"
	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr     = "10.0.0.98:6379"
	httpaddr = "0.0.0.0:9998"
)

func NewClient() *redis.Client {
	client := redis.NewClient(
		&redis.Options{
			Network: "tcp",
			Addr:    addr,
		},
	)
	pong, err := client.Ping().Result()
	if err != nil || pong != "PONG" {
		log.Fatal("[error]connection redis is failed：", err)
	}
	return client
}

func main() {
	client := NewClient()

	prometheus.Register(controllera.NewClientConnCollector(client))
	prometheus.Register(controllera.NewUpRedisCollector(client))

	http.Handle("/metrics/", promhttp.Handler())

	http.ListenAndServe(httpaddr, nil)
}
