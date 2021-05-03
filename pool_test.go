package kevago

import (
	"context"
	"net"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func dummyDialer(ctx context.Context) (net.Conn, error) {
	return &net.TCPConn{}, nil
}

var _ = Describe("MinIdleConns", func() {
	const poolSize = 100
	// ctx := context.Background()
	createPoolFunc := func(minIdleConns int) *ConnPool {
		connPool, err := NewConnPool(Options{
			Dialer:             dummyDialer,
			PoolSize:           poolSize,
			MinIdleConn:        minIdleConns,
			PoolTimeout:        100 * time.Millisecond,
			IdleTimeout:        -1,
			IdleCheckFrequency: -1,
		})
		Eventually(func() int {
			return connPool.TotalConn()
		}).Should(Equal(minIdleConns))
		Eventually(func() int {
			return connPool.TotalIdleConns()
		}).Should(Equal(minIdleConns))
		Expect(err).NotTo(HaveOccurred())
		return connPool
	}

	assertPoolGetPut := func(minIdleConn int) {
		var connPool *ConnPool
		BeforeEach(func() {
			connPool = createPoolFunc(minIdleConn)
		})
		AfterEach(func() {
			connPool.Close()
		})

		It("has idle connections", func() {
			Expect(connPool.TotalConn()).To(Equal(minIdleConn))
			Expect(connPool.TotalIdleConns()).To(Equal(minIdleConn))
		})

		When("Get one", func() {
			var cn *Conn
			BeforeEach(func() {
				incn, err := connPool.Get()
				Expect(err).NotTo(HaveOccurred())
				cn = incn

				// //wait for min idle to be ensured
				Eventually(func() int {
					return connPool.TotalIdleConns()
				}).Should(Equal(minIdleConn))

			})
			It("has idle connections", func() {
				Eventually(func() int {
					return connPool.TotalIdleConns()
				}).Should(Equal(minIdleConn))

				Expect(connPool.TotalIdleConns()).To(Equal(minIdleConn))

			})
			When("Put back", func() {
				BeforeEach(func() {
					connPool.Put(cn)
				})

				It("has idle connections", func() {
					Expect(connPool.TotalIdleConns()).To(Equal(minIdleConn + 1))
				})
			})

		})

	}
	Context("minIdleConns = 1", func() {
		assertPoolGetPut(1)
	})

	Context("minIdleConns = 32", func() {
		assertPoolGetPut(32)
	})

})

// var _ = Describe("conns reaper", func() {
// 	const idleTimeout = time.Minute
// 	const maxAge = time.Hour
// 	var closedConns []*Conn

// 	assert := func(typ string) {
// 		BeforeEach(func() {
// 			connPool, err := NewConnPool(Options{
// 				Dialer:             dummyDialer,
// 				PoolSize:           10,
// 				IdleTimeout:        idleTimeout,
// 				MaxConnAge:         maxAge,
// 				PoolTimeout:        time.Second,
// 				IdleCheckFrequency: time.Hour,
// 				OnConnClosed: func(cn *Conn) error {
// 					closedConns = append(closedConns, cn)
// 					return nil
// 				},
// 			})
// 			Expect(err).NotTo(HaveOccurred())
// 		})
// 	}
// })
