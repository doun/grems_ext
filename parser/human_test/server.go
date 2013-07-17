package main

import (
	rpt ".."
	"encoding/json"
	"flag"
	"fmt"
	"github.com/djimenez/iconv-go"
	"io/ioutil"
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

type RPT struct {
	Raw  string
	Info map[string]map[string]string
	Val  map[string]map[string]string
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
		f, err := ioutil.TempFile(dir+"/files", "rpt_")
		fname := f.Name()
		if err != nil {
			panic(err)
		}
		uploadedf, err := h.Open()
		if err != nil {
			panic(err)
		}
		content, err := ioutil.ReadAll(uploadedf)
		if err != nil {
			panic(err)
		}
		f.Write(content)
		f.Close()
		parsedRst := rpt.GammaRPT{}
		err = parsedRst.ReadFile(fname)
		if err != nil {
			rst, _ := json.Marshal(map[string]string{"error": "true", "description": err.Error()})
			w.Write(rst)
			return
		}

		rst := RPT{}
		rst.Raw = string(content)
		rst.Info = make(map[string]map[string]string)
		rst.Val = make(map[string]map[string]string)

		rst.Info["体积"] = map[string]string{"Name":"体积","Value":fmt.Sprintf("%f", parsedRst.SVol)}
		rst.Info["体积单位"] = map[string]string{"Name":"体积单位","Value":parsedRst.SVolUnit}
		rst.Info["测量时长"] = map[string]string{"Name":"测量时长","Value":fmt.Sprintf("%d", parsedRst.LiveSeconds)}
		rst.Info["测量时间"] = map[string]string{"Name":"测量时间","Value":parsedRst.AcqStartTime.Format("2006-01-02 15:04:05")}
		rst.Info["样品类型"] = map[string]string{"Name":"样品类型","Value":parsedRst.SType}
		for key, val := range parsedRst.Nuclides {
			if val.IsLLD {
				rst.Val[key] = map[string]string{"Name":key,"Value":fmt.Sprintf("<%3.2E", val.Act)}
			} else {
				rst.Val[key] = map[string]string{"Name":key,"Value":fmt.Sprintf("%3.2E", val.Act)}
			}
		}
		if js, e := json.Marshal(rst); e == nil {
			w.Write(js)
		} else {
			var e1 error
			if rst.Raw, e1 = iconv.ConvertString(rst.Raw, "GB2312", "UTF-8"); e1 == nil {
				if js, e2 := json.Marshal(rst); e2 == nil {
					w.Write(js)
				} else {
					w.Write([]byte("converted, but still failed:" + e2.Error()))
				}
			} else {
				w.Write([]byte("convert failed:" + e1.Error()))
			}
		}
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
