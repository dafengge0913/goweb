package main

import (
	"fmt"
	"github.com/dafengge0913/goweb"
)

func main() {
	s := goweb.NewServer()
	s.AddRouter("/hello1/haha", Haha)
	s.AddRouter("/hello[0-9]+", HelloN)
	s.AddRouter("/hello1", Hello1)
	if err := s.Serve(":8888"); err != nil {
		fmt.Println("start server error: ", err)
	}
}

func HelloN(ctx *goweb.Context) {
	for k, v := range ctx.Params {
		fmt.Printf("HelloN : %v -> %v \n", k, v)
	}
}

func Hello1(ctx *goweb.Context) {
	for k, v := range ctx.Params {
		fmt.Printf("Hello1 : %v -> %v \n", k, v)
	}
}

func Haha(ctx *goweb.Context) {
	for k, v := range ctx.Params {
		fmt.Printf("Haha : %v -> %v \n", k, v)
	}
}
