package logging

import (
	"os"
	"sync"
)

type LoggerFactory struct {
	loggers map[string]*Logging
	once    sync.Once
	m       sync.Mutex
}

func (lm *LoggerFactory) InitLogger() {
	//sync.Once 여러 고루틴에서 실행해도 해당 함수는 한번만 실행
	lm.once.Do(func() {
		if lm.loggers == nil {
			lm.loggers = make(map[string]*Logging)
		}
	})
}

/*
Get method returns the logger for a given loggerName.
If there is no logger for loggerName, then new logger is created and returned.
*/
//multi-process들이 공유할수있는 로그 파일을 만든다
//로그 key -value -> key를 주면 가져온다
func (lm *LoggerFactory) Get(loggerName string, logPath string, enableStdout bool) *Logging {
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		//panic(err)
	}
	//여러 함수에서 get -> lock 
	lm.m.Lock()
	defer lm.m.Unlock()

	l, ok := lm.loggers[loggerName]
	if ok == false {
		l = new(Logging)
		l.New(loggerName, logPath, enableStdout, false)
		lm.loggers[loggerName] = l
	}

	return l
}

/*
Get method returns the logger for a given loggerName.
If there is no logger for loggerName, then new logger is created and returned.
*/
func (lm *LoggerFactory) GetNoFlag(loggerName string, logPath string, enableStdout bool) *Logging {
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		//panic(err)
	}

	// check if logger is already in map if not create a new one.
	lm.m.Lock()
	defer lm.m.Unlock()

	l, ok := lm.loggers[loggerName]
	if ok == false {
		l = new(Logging)
		l.New(loggerName, logPath, enableStdout, true)
		lm.loggers[loggerName] = l
	}

	return l
}
