package main

import (
	"flag"

	"github.com/JanitSri/codecrafters-build-your-own-redis/data"
	"github.com/JanitSri/codecrafters-build-your-own-redis/redis"
)

func main() {
	var d string
	flag.StringVar(&d, "dir", "", "the path to the directory where the RDB file is stored (example: /tmp/redis-data)")
	var db string
	flag.StringVar(&db, "dbfilename", "", "the name of the RDB file (example: rdbfile)")
	var port string
	flag.StringVar(&port, "port", "6379", "redis server port number")
	flag.Parse()

	rc := data.NewRedisConfig(d, db)
	sc := redis.NewServerConfig("tcp", "0.0.0.0", port)

	rr := &data.Replication{
		Role: "master",
	}
	ri := &data.RedisInfo{
		Replication: rr,
	}

	rs := redis.NewRedisServer(*sc, *rc, ri)
	rs.Run()
}
