package kevago

import (
	"fmt"
	"strings"

	"github.com/duongcongtoai/kevago/proto"
)

type getCmd struct {
	input  []string
	result string
}

func (g *getCmd) Name() string {
	return "get"
}

func (g *getCmd) Args() []string {
	return g.input
}
func (g *getCmd) ReadResult(r *proto.Reader) error {
	bs, _, err := r.ReadLine()
	if err != nil {
		return err
	}
	g.result = string(bs)
	return nil
}

var getHandler = CmdHandlers{
	name: "get",
	read: func(r *proto.Reader, c Cmd) error {
		return c.ReadResult(r)
	},
	write: func(w *proto.Writer, c Cmd) error {
		_, err := w.WriteString(fmt.Sprintf("%s %s\n", c.Name(), strings.Join(c.Args(), " ")))
		if err != nil {
			return err
		}
		return w.Flush()
	},
}
