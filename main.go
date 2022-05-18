package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/go-redis/redis/v8"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"time"
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)

type server struct {
	redis redis.UniversalClient
	logger *zap.Logger
}

var pingCounter = prometheus.NewCounter(
   prometheus.CounterOpts{
       Name: "ping_request_count",
       Help: "No of request handled by Ping handler",
   },
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("unable to initialize logger")
	}

	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{os.Getenv("REDIS_ADDR")},
		Password: "",
		DB:       0,
	})

	srv := &server{
		redis: rdb,
		logger: logger,
	}

	router := httprouter.New()
        
        prometheus.MustRegister(pingCounter)

	router.GET("/", srv.indexHandler)
        router.Handler("GET", "/metrics", promhttp.Handler())

	logger.Info("server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func (s *server) indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var v string
	var err error
	if v, err = s.redis.Get(context.Background(), "updated_time").Result(); err != nil {
		s.logger.Info("updated_time not found, setting it")
		v = time.Now().Format("2006-01-02 03:04:05")
		s.redis.Set(context.Background(), "updated_time", v, 5*time.Second)
	} else {
		s.logger.Info("got updated_time")
	}
        pingCounter.Inc()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello world: updated_time=%s\n", v)
}
