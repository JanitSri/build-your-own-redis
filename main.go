package main

import "github.com/JanitSri/codecrafters-build-your-own-redis/redis"

func main() {
	c := redis.NewServerConfig("tcp", "0.0.0.0", "6379")
	rs := redis.NewRedisServer(*c)
	rs.Run()
}
