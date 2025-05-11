package main

import (
	"github.com/codecrafters-io/redis-starter-go/app/redis"
)

func main() {
	c := redis.NewServerConfig("tcp", "0.0.0.0", "6379")
	rs := redis.NewRedisServer(*c)
	rs.Run()
}
