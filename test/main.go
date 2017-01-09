package main

import (
	"fmt"
	"github.com/dafengge0913/goweb"
	"github.com/hoisie/web"
	"regexp"
)

func main2() {
	//web.Get("/hello[0-9]+", HelloWeb)
	web.Get("/hello[2-9]+", HelloWeb)
	web.Get("/hello1", HelloWeb2)
	web.Run("0.0.0.0:8888")
}

func main() {
	re1, _ := regexp.Compile("/hello1/haha")
	re2, _ := regexp.Compile("/hello1")
	re3, _ := regexp.Compile("/hello[0-9]+")
	res := []*regexp.Regexp{re1, re2, re3}
	strs := []string{"/hello1/haha", "/hello1"}
	for i, re := range res {
		for _, s := range strs {
			ss := re.FindStringSubmatch(s)
			b := re.MatchString(s)
			fmt.Println(i, " ", ss, " ", b)
			//for _,ts := range ss {
			//}
		}

	}

}

func main3() {
	s := goweb.NewServer()
	s.AddRouter("/hello2", Hello2)
	s.AddRouter("/hello1/haha", Haha)
	s.AddRouter("/hello[0-9]+", Hello1)
	s.AddRouter("/hello3", Hello3)
	s.AddRouter("/hello/d", Hello2)
	if err := s.Serve(":8888"); err != nil {
		fmt.Println("start server error: ", err)
	}
}

func Hello1(ctx *goweb.Context) {
	for k, v := range ctx.Params() {
		fmt.Printf("Hello1 : %v -> %v \n", k, v)
	}
}

func Hello2(ctx *goweb.Context) {
	for k, v := range ctx.Params() {
		fmt.Printf("Hello2 : %v -> %v \n", k, v)
	}
}

func Hello3(ctx *goweb.Context) {
	for k, v := range ctx.Params() {
		fmt.Printf("Hello3 : %v -> %v \n", k, v)
	}
}

func Haha(ctx *goweb.Context) {
	for k, v := range ctx.Params() {
		fmt.Printf("Haha : %v -> %v \n", k, v)
	}
}

func HelloWeb(ctx *web.Context) {
	for k, v := range ctx.Params {
		fmt.Printf("HelloWeb : %v -> %v \n", k, v)
	}
}

func HelloWeb2(ctx *web.Context) {
	for k, v := range ctx.Params {
		fmt.Printf("HelloWeb2 : %v -> %v \n", k, v)
	}
}
