package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

var dir string
var port int
var staticHandler http.Handler
var indexHandler http.Handler

// 初始化参数
func init() {
	dir = path.Dir(os.Args[0])
	flag.IntVar(&port, "port", 80, "服务器端口")
	flag.Parse()
	staticHandler = http.FileServer(http.Dir(dir))
}

func main() {

	http.HandleFunc("/parser", parser)
	http.HandleFunc("/", StaticServer)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func parser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	_, h, e := r.FormFile("ff")
	if e != nil {
		w.Write([]byte("no..." + e.Error()))
	} else {
		w.Write([]byte("ok!" + h.Filename))
	}
}

// 静态文件处理
func StaticServer(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		fmt.Println(req.URL.String())
		staticHandler.ServeHTTP(w, req)
		return
	} else {
		fmt.Println(req.URL.String())
		http.ServeFile(w, req, dir+"/index.html")
		return
	}
}
