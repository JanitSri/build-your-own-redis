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
	flag.Parse()

	rc := data.NewRedisConfig(d, db)
	sc := redis.NewServerConfig("tcp", "0.0.0.0", "6379")
	rs := redis.NewRedisServer(*sc, *rc)
	rs.Run()
}
