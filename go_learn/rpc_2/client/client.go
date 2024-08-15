package main

import (
	"fmt"
	"net/rpc/jsonrpc"
)

type Params struct {
	Param1 float64 `json:"param1"`
	Param2 float64 `json:"param2"`
}

func main() {
	dial, err := jsonrpc.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println(err)
	}
	var area float64
	err = dial.Call("Rect.RectangleArea", Params{85.6, 19.8}, &area)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("area is", area)

	var c float64
	err = dial.Call("Rect.RectangleC", Params{85.6, 19.8}, &c)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("c is", c)
}
