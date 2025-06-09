package replication

import (
	"context"
	"log"
	"sync"

	"github.com/JanitSri/codecrafters-build-your-own-redis/customerror"
	"github.com/JanitSri/codecrafters-build-your-own-redis/redis"
)

type Operator struct {
	leader    *redis.RedisServer
	followers []*redis.RedisServer
}

func NewOperator() *Operator {
	return &Operator{}
}

func (op *Operator) Join(rs *redis.RedisServer) {
	// uncomment to support leader and followers at the same time
	// codecrafter tests only require either a leader or follower
	// role := rs.RedisContext.RedisInfo.Replication.Role

	// if role == "" || role == "master" {
	// 	op.leader = rs
	// 	return
	// }

	// op.followers = append(op.followers, rs)

	op.leader = rs
}

func (op *Operator) Start(ctx context.Context, wg *sync.WaitGroup) {
	if op.leader == nil {
		log.Fatalln(customerror.NoLeaderAvailableError{})
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
