package main

import (
	"bufio"
	"code.google.com/p/log4go"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	str "strings"
	"time"
)

var log log4go.Logger

func init() {
	log = make(log4go.Logger)
	log.AddFilter("stdout", log4go.INFO, log4go.NewFormatLogWriter(os.Stdout, "[%D %T][%L][Conn](%S) %M"))
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
	bufReader    *bufio.Reader
}

type NuclideActivity struct {
	Name  string
	Act   float32
	IsLLD bool
}

func (self *GammaRPT) Parse(filePath string) error {
	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	file := string(f)
	self.file = &file
	self.fReader = str.NewReader(file)
	self.bufReader = bufio.NewReader(self.fReader)

	volStart := str.Index(file, "Sample Size")
	if volStart < 0 {
		return errors.New("未找到样品体积")
	}
	self.fReader.Seek(int64(volStart), 0)
	volStr, _, err := self.bufReader.ReadLine()
	if err != nil {
		log.Debug("未读到体积行,已经找到的为:", volStr)
		panic(err)
	}
	volParts := str.Split(string(volStr), ":")
	if len(volParts) != 2 {
		panic("体积解析出错")
	}
	volVals := str.Split(str.TrimSpace(volParts[1]), " ")
	volValStr := volVals[0]
	volUnit := volVals[1]

	fmt.Sscanf(volValStr, "%f", &self.SVol)
	fmt.Sscanf(volUnit, "%s", &self.SVolUnit)

	titleStart := str.Index(file, "Sample Title")
	if titleStart < 0 {
		return errors.New("未找到样品名称")
	}
	self.fReader.Seek(int64(titleStart), 0)
	return nil
}
