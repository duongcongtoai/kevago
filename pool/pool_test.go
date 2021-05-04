package pool

import (
	"context"
	"net"
	"sync"
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
			return connPool.TotalConns()
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
			Expect(connPool.TotalConns()).To(Equal(minIdleConn))
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

var _ = Describe("conns reaper", func() {
	const idleTimeout = 5 * time.Minute
	const maxAge = time.Hour
	const poolSize = 10
	const minIdle = 0
	closedConnL := new(sync.Mutex)
	var closedConns []*Conn
	var connPool *ConnPool

	assertConnsReaperWork := func() {
		BeforeEach(func() {
			initConnPool, err := NewConnPool(Options{
				Dialer:             dummyDialer,
				PoolSize:           poolSize,
				IdleTimeout:        idleTimeout,
				MinIdleConn:        minIdle,
				MaxConnAge:         maxAge,
				PoolTimeout:        time.Second,
				IdleCheckFrequency: time.Hour,
				OnConnClosed: func(cn *Conn) error {
					closedConnL.Lock()
					closedConns = append(closedConns, cn)
					closedConnL.Unlock()
					return nil
				},
			})
			connPool = initConnPool
			Expect(err).NotTo(HaveOccurred())
			Eventually(func() int {
				return connPool.TotalIdleConns()
			}).Should(Equal(minIdle))
			var temp []*Conn = nil

			for i := 0; i < minIdle; i++ {
				con, err := connPool.Get()
				Expect(err).NotTo(HaveOccurred())
				temp = append(temp, con)
			}
			for _, item := range temp {
				item.SetLastUsed(time.Now().Add(-2 * idleTimeout)) //make conn idle time out
				connPool.Put(item)
			}
			Eventually(func() int {
				return connPool.TotalIdleConns()
			}).Should(Equal(0))

		})
		It("staled connections are closed", func() {
			Expect(len(closedConns)).Should(Equal(minIdle))
		})
		// It("fresh connections remain equal min idles", func() {
		// 	Expect(connPool.TotalConns()).Should(Equal(minIdle))
		// 	Expect(connPool.TotalIdleConns()).Should(Equal(minIdle))
		// })
	}
	assertConnsReaperWork()
})
