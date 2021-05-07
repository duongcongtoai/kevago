package kevago

import (
	"fmt"

	"github.com/duongcongtoai/kevago/pool"
	"github.com/duongcongtoai/kevago/proto"
)

type commander struct {
	internal map[string]CmdHandlers
}
type CmdHandlers struct {
	read  func(r *proto.Reader, c Cmd) error
	write func(w *proto.Writer, c Cmd) error
}

type Cmd interface {
	Name() string
	Args() []string

	// resultExtractor(r *reader) error
}

var globalCmd = commander{
	internal: make(map[string]CmdHandlers),
}

func registerCmd(name string, h CmdHandlers) {
	globalCmd.internal[name] = h
}

func (c commander) execute(conn *pool.Conn, comd Cmd) error {
	hs, exist := c.internal[comd.Name()]
	if !exist {
		return fmt.Errorf("command %s not found", comd.Name())
	}
	err := conn.WriteIntercept(func(w *proto.Writer) error {
		return hs.write(w, comd)
	})
	// err := hs.write(conn.w, comd)
	if err != nil {
		return err
	}
	return conn.ReadIntercept(func(w *proto.Reader) error {
		return hs.read(w, comd)
	})
}
