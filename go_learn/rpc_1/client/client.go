package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	type Param struct {
		Param1 float64 `json:"param1"`
		Param2 float64 `json:"param2"`
	}
	http, err := rpc.DialHTTP("tcp", "localhost:8880")
	if err != nil {
		fmt.Println(err)
	}
	var res float64
	err = http.Call("E.MultiTwoNumber", Param{1.2, 2.2}, &res)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
