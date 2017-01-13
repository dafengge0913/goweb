package main

import (
	"fmt"
	"html/template"

	"github.com/dafengge0913/goweb"
)

var s *goweb.Server

func main() {
	s = goweb.NewServer()
	s.AddRouter("/hello1/haha", Haha)
	s.AddRouter("/hello[0-9]+", HelloN)
	s.AddRouter("/hello1", Hello1)
	s.AddStaticRouter("/css/", "/css/", "example/css")
	s.AddStaticRouter("/js/", "/js/", "example/js")
	s.AddStaticRouter("/.*\\.html", "/", "example") // match xxx.html not xxx.htm
	s.AddRouter("/admin", admin)                    // both admin.html and admin can be access
	s.AddRouter("/helloAjax", helloAjax)
	s.AddRouter("/addCookie", addCookie)
	s.AddRouter("/delCookie", delCookie)
	s.AddRouter("/login", login)
	s.AddRouter("/close", closeServer) // just for testing
	if err := s.Serve(":8888"); err != nil {
		fmt.Println("start server error: ", err)
	}
}

func HelloN(ctx *goweb.Context) {
	for k, v := range ctx.Params() {
		fmt.Printf("HelloN : %v -> %v \n", k, v)
	}
}

func Hello1(ctx *goweb.Context) {
	name := ""
	for k, v := range ctx.Params() {
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
	for k, v := range ctx.Params() {
		fmt.Printf("Haha : %v -> %v \n", k, v)
	}
}

func helloAjax(ctx *goweb.Context) {
	name := ""
	fmt.Printf("param len: : %d \n", len(ctx.JSONParams()))
	for k, v := range ctx.JSONParams() {
		fmt.Printf("helloAjax : %v -> %v \n", k, v)
		name = v
	}
	ctx.ResponseJSON("Ajax hello " + name)
}

func addCookie(ctx *goweb.Context) {
	ctx.SetCookie("sid", "12345678")
}

func delCookie(ctx *goweb.Context) {
	ctx.DelCookie("sid")
}

func login(ctx *goweb.Context) {
	ctx.Session().Set("username", ctx.Param("username"))
}

func admin(ctx *goweb.Context) {
	if name, fd := ctx.Session().Get("username"); fd {
		fmt.Printf("username : %s  \n", name)
	}
}

func closeServer(_ *goweb.Context) {
	s.Close()
}
