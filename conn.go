package kevago

import (
	"sync/atomic"
	"time"
)

func (c *Conn) lastUsed() time.Time {
	unix := atomic.LoadInt64(&c.usedAt)
	return time.Unix(unix, 0)
}
func (c *Conn) IsManaged() bool {
	return c.managed
}
