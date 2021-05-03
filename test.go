package kevago

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

//
func st() {
	c := redis.NewClient(&redis.Options{})
	_ = c.BRPop(context.Background(), time.Second, "hello")
	r := redis.NewRing(&redis.RingOptions{})
	_ = r.BRPop(context.Background(), time.Second, "hello")
	cl := redis.NewClusterClient(&redis.ClusterOptions{})
	_ = cl.BRPop(context.Background(), time.Second, "hello")
}
