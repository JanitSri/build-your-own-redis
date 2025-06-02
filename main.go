package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/JanitSri/codecrafters-build-your-own-redis/data"
	"github.com/JanitSri/codecrafters-build-your-own-redis/redis"
	"github.com/JanitSri/codecrafters-build-your-own-redis/replication"
)

func main() {
	var d string
	flag.StringVar(&d, "dir", "", "the path to the directory where the RDB file is stored (example: /tmp/redis-data)")
	var db string
	flag.StringVar(&db, "dbfilename", "", "the name of the RDB file (example: rdbfile)")
	var port string
	flag.StringVar(&port, "port", "6379", "redis server port number")
	var replicaOf string
	flag.StringVar(&replicaOf, "replicaof", "", "redis server port number")

	flag.Parse()

	var wg sync.WaitGroup

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := <-sigChan
		log.Println("main received shutdown signal:", sig)
		cancel()
	}()

	op := replication.NewOperator()
	role := "master"
	if replicaOf != "" {
		role = "slave"
	}
	leader := createRedisServer(d, db, port, role)

	op.Join(leader)

	op.Start(ctx, &wg)

	<-ctx.Done()

	wg.Wait()
	log.Println("main exiting")
}

func createRedisServer(dir, dbFilename, port, role string) *redis.RedisServer {
	rc := data.NewRedisConfig(dir, dbFilename)
	sc := redis.NewServerConfig("tcp", "0.0.0.0", port)

	rr := &data.Replication{
		Role:             role,
		MasterReplid:     "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		MasterReplOffset: 0,
	}
	ri := &data.RedisInfo{
		Replication: rr,
	}

	return redis.NewRedisServer(*sc, *rc, ri)
}
