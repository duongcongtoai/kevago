package pool

import (
	"bufio"
	"fmt"
)

type commander struct {
	internal map[string]cmdHandlers
}
type cmdHandlers struct {
	read  func(r *reader, c cmd) (result, error)
	write func(w *writer, c cmd) error
}

type result struct{}

type reader bufio.Reader

type writer bufio.Writer

type cmd struct {
	name string
	args []string
}

var globalCmd = commander{
	internal: make(map[string]cmdHandlers),
}

func registercmd(name string, h cmdHandlers) {
	globalCmd.internal[name] = h
}

func (c commander) execute(conn *Conn, comd cmd) (result, error) {
	hs, exist := c.internal[comd.name]
	if !exist {
		return result{}, fmt.Errorf("command %s not found", comd.name)
	}
	err := hs.write(conn.w, comd)
	if err != nil {
		return result{}, err
	}
	return hs.read(conn.r, comd)
}
