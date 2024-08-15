package main

import (
	"fmt"
	"net/http"
	"net/rpc"
)

type Param struct {
	Param1 float64 `json:"param1"`
	Param2 float64 `json:"param2"`
}

type E struct{}

func (e *E) MultiTwoNumber(p Param, res *float64) error {
	*res = p.Param1 * p.Param2
	return nil
}

func main() {
	var server = new(E)
	err := rpc.Register(server)
	if err != nil {
		fmt.Println(err)
		return
	}
	rpc.HandleHTTP()
	err = http.ListenAndServe(":8880", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
