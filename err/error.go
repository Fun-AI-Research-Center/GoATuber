package err

import (
	"log"
	"os"
	"sync"
	"time"
)

//错误处理

// 错误等级
const (
	Normal = 1 //可忽略的错误
	Fatal  = 2 //应当中断程序执行的错误
)

// 错误日志变量
var (
	errorLog *os.File   //错误日志文件
	mu       sync.Mutex //互斥锁
)

// Error 错误处理。入参：错误信息，错误等级
func Error(err error, level int) {
	switch level {
	case Normal:
		handelNormalError(err)
	case Fatal:
		handelFatalError(err)
	}
}

// handelNormalError 处理可忽略的错误
func handelNormalError(e error) {
	var err error
	//如果错误日志不存在，则创建错误日志
	if errorLog == nil {
		currentTime := time.Now().Format("2006-01-02 15-04-05")
		errorLog, err = os.Create("./" + currentTime + " err.log")
		if err != nil {
			panic(err)
		}
	}
	//将错误信息追加到日志的尾行
	mu.Lock()
	file, err := os.OpenFile(errorLog.Name(), os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
		mu.Unlock()
	}()

	_, err = file.Write([]byte(time.Now().Format("2006-01-02 15:04:05") + ":" + e.Error() + "\n"))
	if err != nil {
		return
	}
	if err != nil {
		log.Println(err)
	}

	log.Println(e)
}

// handelFatalError 处理应当中断程序执行的错误
func handelFatalError(e error) {
	var err error
	//如果错误日志不存在，则创建错误日志
	if errorLog == nil {
		currentTime := time.Now().Format("2006-01-02 15-04-05")
		errorLog, err = os.Create("./" + currentTime + " err.log")
		if err != nil {
			panic(err)
		}
	}
	//将错误信息追加到日志的尾行
	mu.Lock()
	file, err := os.OpenFile(errorLog.Name(), os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
		mu.Unlock() //还有什么开锁的必要吗？
	}()
	_, err = file.Write([]byte(time.Now().Format("2006-01-02 15:04:05") + ":" + e.Error() + "\n"))
	if err != nil {
		return
	}
	if err != nil {
		log.Println(err)
	}

	log.Println(e, "。程序于十秒后退出")
	time.Sleep(10 * time.Second)
	os.Exit(1)
}
