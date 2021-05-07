package kevago

import "github.com/duongcongtoai/kevago/proto"

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

var getHandler = CmdHandlers{
	read: func(r *proto.Reader, c Cmd) error {
		return nil
	},
	write: func(w *proto.Writer, c Cmd) error {
		return nil
	},
}
