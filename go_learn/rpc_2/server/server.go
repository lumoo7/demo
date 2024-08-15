package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Params struct {
	Param1 float64 `json:"param1"`
	Param2 float64 `json:"param2"`
}

type Rect struct{}

func (r *Rect) RectangleArea(p Params, res *float64) error {
	*res = p.Param1 * p.Param2
	return nil
}

func (r *Rect) RectangleC(p Params, res *float64) error {
	*res = 2 * (p.Param1 + p.Param2)
	return nil
}

func main() {
	var rect = new(Rect)
	err := rpc.Register(rect)
	if err != nil {
		fmt.Println(err)
	}
	listen, err := net.Listen("tcp", ":1234")
	for {
		accept, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go func(c net.Conn) {
			jsonrpc.ServeConn(accept)
		}(accept)
	}
}
