package replication

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/JanitSri/codecrafters-build-your-own-redis/redis"
)

var noLeaderAvailableException = errors.New("no redis leader available")

type Operator struct {
	leader    *redis.RedisServer
	followers []*redis.RedisServer
}

func NewOperator() *Operator {
	return &Operator{}
}

func (op *Operator) Join(rs *redis.RedisServer) {
	role := rs.RedisContext.RedisInfo.Replication.Role

	if role == "" || role == "master" {
		op.leader = rs
		return
	}

	op.followers = append(op.followers, rs)
}

func (op *Operator) Start(ctx context.Context, wg *sync.WaitGroup) {
	if op.leader == nil {
		log.Fatalln(noLeaderAvailableException)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		op.leader.StartupTasks()
		op.leader.Run(ctx)
	}()

	for _, f := range op.followers {
		wg.Add(1)
		go func(f *redis.RedisServer) {
			defer wg.Done()
			f.Run(ctx)
		}(f)
	}
}
