package pool

// import (
// 	"context"
// 	"net"

// 	"golang.org/x/net/context"
// )

// type Config struct {
// 	Endpoints []string
// }

// type Client struct {
// 	conn   net.Conn //TODO: connection pool
// 	ctx    context.Context
// 	cancel context.CancelFunc
// }

// func (c *Client) Close() error {

// }

// func NewClient(c Config) (*Client, error) {
// 	conn, err := net.Dial("tcp", c.Endpoints[0])
// 	if err != nil {
// 		return nil, err
// 	}

// 	cl := &Client{conn: conn}
// 	ctx, cancel := context.WithCancel(context.Background())
// 	go cl.readLoop(ctx)
// 	go cl.writeLoop(ctx)
// 	cl.cancel = cancel
// 	cl.ctx = ctx
// 	return cl, nil
// }

// func (c *Client) Get(key string) (interface{}, error) {
// 	//Get a connection from pool
// 	//Retry if fail
// 	//Include retry backoff
// 	//Write to socket
// 	//Read from socket
// 	return nil, nil
// }

// func (c *Client) readLoop() {

// }
// func (c *Client) writeLoop() {

// }
