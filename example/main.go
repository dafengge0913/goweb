package main

import (
	"fmt"
	"github.com/dafengge0913/goweb"
	"html/template"
)

func main() {
	s := goweb.NewServer()
	s.AddRouter("/hello1/haha", Haha)
	s.AddRouter("/hello[0-9]+", HelloN)
	s.AddRouter("/hello1", Hello1)
	s.AddStaticRouter("/css/", "example/css")
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
	name := ""
	for k, v := range ctx.Params {
		fmt.Printf("Hello1 : %v -> %v \n", k, v)
		name = v
	}
	if tpl, err := template.ParseFiles("example/index.html"); err != nil {
		fmt.Println("parse template error: ", err)
	} else {
		data := make(map[string]string)
		data["Name"] = name
		ctx.ResponseTemplate(tpl, data)
	}
}

func Haha(ctx *goweb.Context) {
	for k, v := range ctx.Params {
		fmt.Printf("Haha : %v -> %v \n", k, v)
	}
}
