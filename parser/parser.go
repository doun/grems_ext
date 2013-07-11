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
	//体积
	volStr, err := self.readElement(file, "Sample Size")
	if err != nil {
		return
	}
	n, err := fmt.Sscanf(volStr, "%f %s", &self.SVol, &self.SVolUnit)
	if n != 2 || err != nil {
		return
	}
	//标题
	titleStr, err := self.readElement(file, "Sample Title")
	if err != nil {
		return
	}
	self.STitle = titleStr
	//样品类型
	self.SType, err = self.readElement(file, "Sample Type")
	if err != nil {
		return
	}
	//测量时长
	liveSec, err := self.readElement(file, "Live Time")
	if err != nil {
		return
	}
	n, err = fmt.Sscanf(liveSec, "%d", &self.LiveSeconds)
	if err != nil {
		return
	}
	if n != 1 {
		return errors.New("测量时长解析出错")
	}

	//测量开始时间
	startTimeStr, err := self.readElement(file, "Acquisition Started")
	if err != nil {
		return
	}
	//能量刻度时间和效率刻度时间暂时不用
	log.Debug("time string is:%v", startTimeStr)
	//经测试，一位时间也可
	self.AcqStartTime, err = time.Parse("2006-01-02 15:04:05", startTimeStr)
	if err != nil {
		return
	}

	//Corrected表
	IdentTableLines, err := self.readTable(file, "I N T E R F E R E N C E   C O R R E C T E D   R E P O R T",
		"       ? = nuclide is part of an undetermined solution")
	log.Debug(IdentTableLines)
	self.parseIdent(IdentTableLines)
	//MDA表
	MdaLines, err := self.readTable(file, "", "")
	self.parseMda(MdaLines)

	return
}

func (self *GammaRPT) parseIdent(lines string) (err error) {
	if self.Nuclides == nil {
		self.Nuclides = make(map[string]NuclideActivity)
	}
	actStart := str.Index(lines, "Uncertainty\n\n")
	reader := bufio.NewReader(str.NewReader(lines[actStart+len("Uncertainty\n\n"):]))
	for {
		lineBt, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		line := string(lineBt)
		parts := str.Split(line, " ")
		if len(parts) != 4 {
			break
		}
		var act NuclideActivity
		n,e := fmt.Sscanf(parts[2], "%f", &act.Act)
		if n!=1 || e!=nil{
			return errors.New("解析活度出错")
		}
		act.Name = parts[0]
		self.Nuclides[parts[0]] = act
	}
	return
}

func (self *GammaRPT) parseMda(mda string) (err error) {
	if self.Nuclides == nil {
		self.Nuclides = make(map[string]NuclideActivity)
	}
	return
}

func (self *GammaRPT) readElement(content, prefix string) (ele string, err error) {
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

func (self *GammaRPT) readTable(content, title, tail string) (eleLines string, err error) {
	start := str.Index(content, title)
	end := str.Index(content[start:], tail)
	return content[start : start+end], nil
}
