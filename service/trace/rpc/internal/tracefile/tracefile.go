package tracefile

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"time"
)

type TraceData struct {
	Time time.Time
	Name string
	Data any
}

func (t *TraceData) Marshal() string {
	return fmt.Sprintf("%s %s %s\n", t.Time.Format("2006-01-02 15:04:05"), t.Name, t.Data)
}

type TraceFile struct {
	Dir        string
	Prefix     string
	File       *os.File
	writeBuf   *bufio.Writer
	CreateTime time.Time
	MsgChan    chan *TraceData
	StopChan   chan bool
	IsAlive    bool
}

func NewFile(filePath, prefix string) (f *TraceFile, err error) {
	f = &TraceFile{
		Dir:     filePath,
		Prefix:  prefix,
		MsgChan: make(chan *TraceData),
	}

	err = f.createFile()
	if nil != err {
		return nil, err
	}

	//开启接收数据的协程
	f.Start()
	return f, nil
}

func (f *TraceFile) Start() {
	go func(f *TraceFile) {
		f.IsAlive = true
		for {
			select {
			case data := <-f.MsgChan:
				f.doWrite(data)
			case <-f.StopChan:
				logx.Error("trace file writer stopped!")
				f.IsAlive = false
				f.Flush()
				f.CloseFile()
				return
			}
		}
	}(f)
}

func (f *TraceFile) Stop() {
	go func() {
		f.StopChan <- true
	}()
}

func (f *TraceFile) Flush() {
	if nil != f.writeBuf {
		_ = f.writeBuf.Flush()
	}
}

func (f *TraceFile) CloseFile() {
	if nil != f.File {
		err := f.File.Close()
		if err != nil {
			logx.Error("trace file close failed, err: ", err.Error())
		}
	}
}

func (f *TraceFile) createFile() error {
	oldFile := f.File
	defer func(oldFile *os.File) {
		err := oldFile.Close()
		if err != nil {
			logx.Error("trace file close failed, err: ", err.Error())
		}
	}(oldFile)

	f.CreateTime = time.Now()
	timeTag := time.Now().Format("2006-01-02")
	newFile, err := os.OpenFile(fmt.Sprintf("%s/%s%s.log", f.Dir, f.Prefix, timeTag), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	f.File = newFile
	f.writeBuf = bufio.NewWriter(f.File)

	return nil
}

func (f *TraceFile) checkSplit() {
	//按日期拆分文件
	nonce := time.Now()
	if nonce.Month() == f.CreateTime.Month() && nonce.Day() == f.CreateTime.Day() {
		return
	}
	err := f.createFile()
	if err != nil {
		logx.Error("switch trace file failed, err: ", err.Error())
	}
}

func (f *TraceFile) doWrite(data *TraceData) {
	f.checkSplit()

	_, err := f.writeBuf.WriteString(data.Marshal())
	if err != nil {
		logx.Error(err)
	}
	//zglog.Info("doWrite: ", data)
	//Flush将缓存的文件真正写入到文件中
	_ = f.writeBuf.Flush()
}

func (f *TraceFile) WriteTrace(traceName string, data any) error {
	if !f.IsAlive {
		return errors.New("trace file already stopped")
	}
	//构造数据包
	td := &TraceData{
		Time: time.Now(),
		Name: traceName,
		Data: data,
	}
	go func(td *TraceData) {
		f.MsgChan <- td
	}(td)
	return nil
}
