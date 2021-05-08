package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/duongcongtoai/kevago"
	"github.com/duongcongtoai/kevago/pool"
)

func main() {
	popt := pool.Options{
		PoolTimeout: time.Second,
		PoolSize:    20,
		MinIdleConn: 5,
		Dialer: func(ctx context.Context) (net.Conn, error) {
			conn, err := net.Dial("tcp", "localhost:6767")
			if err != nil {
				return nil, err
			}
			return conn, err
		},
		IdleTimeout:        time.Minute * 5,
		MaxConnAge:         time.Minute * 10,
		IdleCheckFrequency: time.Minute * 5,
	}
	cl, err := kevago.NewClient(kevago.Config{
		Pool: popt,
		// Endpoints: []string{}
	})
	if err != nil {
		panic(err)
	}
	// res, err := cl.P
	res, err := cl.Set("key1", "value1")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	res, err = cl.Get("key1")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
