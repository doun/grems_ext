package main

import (
	"bufio"
	"code.google.com/p/log4go"
	"errors"
	"fmt"
	//	mlib "github.com/doun/golib"
	"io/ioutil"
	"os"
	str "strings"
	"time"
)

var log log4go.Logger

func init() {
	log = make(log4go.Logger)
	log.AddFilter("stdout", log4go.DEBUG, log4go.NewFormatLogWriter(os.Stdout, "[%D %T][%L][Conn](%S) %M"))
}

type GammaRPT struct {
	SVol         float32
	SVolUnit     string
	STitle       string
	SType        string
	AcqStartTime time.Time
	LiveSeconds  int32
	RptGenTime   time.Time
	EngCalTime   time.Time
	EffCalTime   time.Time
	Nuclides     map[string]NuclideActivity
	file         *string
	fReader      *str.Reader
}

type NuclideActivity struct {
	Name  string
	Act   float32
	IsLLD bool
}

func (self *GammaRPT) Parse(filePath string) (err error) {
	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	file := string(f)
	self.file = &file
	self.fReader = str.NewReader(file)

	volStr, err := self.parseElement(file, "Sample Size")
	if err != nil {
		return
	}
	n, err := fmt.Sscanf(volStr, "%f %s", &self.SVol, &self.SVolUnit)
	if n != 2 || err != nil {
		return
	}

	titleStr, err := self.parseElement(file, "Sample Title")
	if err != nil {
		return
	}
	self.STitle = titleStr

	startTimeStr, err := self.parseElement(file, "Acquisition Started")
	if err != nil {
		return
	}
	log.Debug("time string is:%v", startTimeStr)
	self.AcqStartTime, err = time.Parse("2006-01-02 15:04:05", startTimeStr)
	if err != nil {
		return
	}

	return
}

func (self *GammaRPT) parseElement(content, prefix string) (ele string, err error) {
	offset := str.Index(content, prefix)
	if offset < 0 {
		err = errors.New("find prefix err:" + prefix)
		log.Info("找参数%v的起点出错，找到的offset为:%v", prefix, offset)
		return
	}
	n, err := self.fReader.Seek(int64(offset), 0)
	if err != nil {
		return
	}
	if n < 0 {
		return
	}

	reader := bufio.NewReader(self.fReader)

	log.Debug("%v的offset为:%v", prefix, offset)
	line, _, err := reader.ReadLine()
	if err != nil {
		return
	}
	parts := str.Split(string(line), " : ")
	if len(parts) != 2 {
		err = errors.New("格式错误，应为prefix:content, prefix为:" + prefix + "当前行为:" + string(line))
		return
	}
	ele = str.TrimSpace(parts[1])
	return
}
