package logging

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Logging struct {
	name string
	log  *log.Logger
	f    *os.File
	path string
}

func (l *Logging) New(loggerName string, logPath string, enableStdout bool, noFlag bool) {
	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)

		defer f.Close()  //위에서 exit되어서 defer안됨
		os.Exit(1) // <- app 죽는거
	}

	l.f = f
	l.name = loggerName
	l.path = logPath

	writers := []io.Writer{f,}
	if enableStdout {
		writers = append(writers, os.Stdout)
	}
	wrt := io.MultiWriter(writers...)

	l.log = log.New(wrt, "[INFO]", log.LstdFlags)
}

func (l *Logging) Info(vals ...interface{}) {
	l.log.SetPrefix("[INFO] ")

	var txt string
	for _, v := range vals {
		switch v.(type) {
		case float64, float32:
			txt += fmt.Sprintf("%f ", v)
		default:
			txt += fmt.Sprintf("%v ", v)
		}
	}
	l.log.Println(txt)
}

func (l *Logging) Error(str ...interface{}) {
	l.log.SetPrefix("[ERROR] ")
	l.log.Println(str)
}

func (l *Logging) Debug(str ...interface{}) {
	l.log.SetPrefix("[DEBUG] ")
	l.log.Println(str)
}

func (l *Logging) Warn(format string, v ...interface{}) {
	l.log.SetPrefix("[WARN] ")
	l.log.Printf(format, v)
}

func (l *Logging) Close() {
	l.f.Close()
	//fmt.Printf("Logger [%s] is closed. \t [%s] \n", l.name, l.path)
}

//func InitLogger(logPath string, enableStdout bool) *Logging {
//	once.Do(func() {
//		l = new(Logging)
//		l.New(logPath, enableStdout)
//	})
//
//	return l
//}
