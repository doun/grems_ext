package main

import (
	rpt ".."
	"flag"
	"fmt"
	"github.com/coocood/jas"
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

type Rpt struct {
	Raw      string
	FileName string
	Info     map[string]map[string]string
	Val      map[string]map[string]string
}

func (self *Rpt) rspWithFile(c *jas.Context, fn string) {
	fname := fn
	c.ResponseHeader.Add("Access-Control-Allow-Origin", "*")
	c.ResponseHeader.Add("Access-Control-Allow-Headers", "Content-Type")
	parsedRst := rpt.GammaRPT{}
	err := parsedRst.ReadFile(fname)
	if err != nil {
		c.Error = jas.NewInternalError(err)
		return
	}

	content, err := ioutil.ReadFile(fname)
	if err != nil {
		c.Error = jas.NewInternalError(err)
		return
	}

	rst := Rpt{}
	rst.Raw = string(content)
	rst.FileName = fname
	rst.Info = make(map[string]map[string]string)
	rst.Val = make(map[string]map[string]string)

	rst.Info["体积"] = map[string]string{"Name": "体积", "Value": fmt.Sprintf("%f", parsedRst.SVol)}
	rst.Info["体积单位"] = map[string]string{"Name": "体积单位", "Value": parsedRst.SVolUnit}
	rst.Info["测量时长"] = map[string]string{"Name": "测量时长", "Value": fmt.Sprintf("%d", parsedRst.LiveSeconds)}
	rst.Info["测量时间"] = map[string]string{"Name": "测量时间", "Value": parsedRst.AcqStartTime.Format("2006-01-02 15:04:05")}
	rst.Info["样品类型"] = map[string]string{"Name": "样品类型", "Value": parsedRst.SType}
	for key, val := range parsedRst.Nuclides {
		if val.IsLLD {
			rst.Val[key] = map[string]string{"Name": key, "Value": fmt.Sprintf("<%3.2E", val.Act)}
		} else {
			rst.Val[key] = map[string]string{"Name": key, "Value": fmt.Sprintf("%3.2E", val.Act)}
		}
	}
	c.Data = nil
	_, e := c.FlushData(rst)
	if e != nil {
		var e1 error
		if rst.Raw, e1 = iconv.ConvertString(rst.Raw, "GB2312", "UTF-8"); e1 == nil {
			if _, e1 = c.FlushData(rst); e1 != nil {
				c.Error = jas.NewInternalError("编码转换成功，但仍无法转换为json")
			}
		} else {
			c.Error = jas.NewInternalError("编码转换出错")
		}
	}

}

func (self *Rpt) Post(c *jas.Context) {
	_, h, e := c.Request.FormFile("ff")
	if e != nil {
		c.Data = e.Error()
	} else {
		f, err := ioutil.TempFile(dir+"/files", "rpt_")
		fname := f.Name()
		log.Println("fname(tmp)is:", fname)
		if err != nil {
			log.Println("新建临时文件出错:", err)
			c.Error = jas.NewInternalError(err)
			return
		}
		uploadedf, err := h.Open()
		if err != nil {
			log.Println("打开上传文件出错:")
			c.Error = jas.NewInternalError(err)
			return
		}
		content, err := ioutil.ReadAll(uploadedf)
		if err != nil {
			log.Println("读取上传文件内容出错:", err)
			c.Error = jas.NewInternalError(err)
			return
		}
		f.Write(content)
		f.Close()
		self.rspWithFile(c, fname)
	}
}

func (self *Rpt) Get(c *jas.Context) {

}

func main() {
	router := jas.NewRouter(&Rpt{})
	router.BasePath = "/v1/"
	log.Println(router.HandledPaths(true))
	http.Handle("/v1/", router)
	http.HandleFunc("/", StaticServer)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func parser(w http.ResponseWriter, r *http.Request) {
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
