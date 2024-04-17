package logger

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"
)

var today = ""
var infoFileHandel *os.File = nil
var errFileHandel *os.File = nil
var systemFileHandel *os.File = nil
var AccessLogFilePath = "log/access.log"

func InitLog() {
	_ = os.Mkdir("log", os.ModePerm)
	systemFileHandel, _ = os.OpenFile("log/system.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	var err error

	date := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	initDir := func(date string, yesterday string, isFirst bool) {
		today = date
		_ = os.Mkdir("log/"+date, os.ModePerm)

		if infoFileHandel != nil {
			err := infoFileHandel.Close()
			if err != nil {
				fmt.Println("日志服务错误【1】" + err.Error())
				return
			}
		}
		infoFileName := "log/" + date + "/info.log"
		infoFileHandel, err = os.OpenFile(infoFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			fmt.Println("日志服务错误【3】" + err.Error())
			return
		}

		if errFileHandel != nil {
			err := errFileHandel.Close()
			if err != nil {
				fmt.Println("日志服务错误【2】" + err.Error())
				return
			}
		}
		errFileName := "log/" + date + "/err.log"
		errFileHandel, err = os.OpenFile(errFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			fmt.Println("日志服务错误【4】" + err.Error())
			return
		}

		if !isFirst {
			_ = os.Mkdir("log/"+yesterday, os.ModePerm)
			fromAccessFile, _ := os.OpenFile(AccessLogFilePath, os.O_RDWR, os.ModePerm)
			toAccessFile, _ := os.OpenFile("log/"+yesterday+"/access.log", os.O_RDWR|os.O_CREATE, os.ModePerm)
			defer func() {
				err := fromAccessFile.Truncate(0)
				if err != nil {
					fmt.Println(err)
					fmt.Println("文件清空失败")
				}
				_, err = fromAccessFile.Seek(0, 0)
				if err != nil {
					fmt.Println(err)
					fmt.Println("文件重置偏移失败")
				}
				_ = fromAccessFile.Close()
				_ = toAccessFile.Close()
			}()
			_, _ = io.Copy(toAccessFile, fromAccessFile)
		}
	}

	initDir(date, yesterday, true)

	go func() {
		for {
			date := time.Now().Format("2006-01-02")
			yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
			if date == today {
				time.Sleep(time.Second)
				continue
			} else {
				initDir(date, yesterday, false)
			}
		}
	}()
	System("LOGGER INIT SUCCESS")
}

func Trace() string {
	b := new(bytes.Buffer)
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		_, err := fmt.Fprintf(b, "%s:%d (0x%x)\n", file, line, pc)
		if err != nil {
			return ""
		}
	}
	return b.String()
}

func Info(msg string, append ...any) {
	fmt.Println(msg, append)
	if infoFileHandel != nil {
		logHandel := slog.New(slog.NewTextHandler(infoFileHandel, nil))
		slog.SetDefault(logHandel)
		slog.Info(msg, append...)
	}
}

func System(msg string, append ...any) {
	fmt.Println(msg, append)
	logHandel := slog.New(slog.NewTextHandler(systemFileHandel, nil))
	slog.SetDefault(logHandel)
	slog.Info(msg, append...)
}

func Error(msg string, appends ...any) {
	t := Trace()
	fmt.Println(t)
	fmt.Print(msg, appends)
	if errFileHandel != nil {
		logHandel := slog.New(slog.NewTextHandler(errFileHandel, nil))
		slog.SetDefault(logHandel)
		n := append([]any{"trace", t}, appends...)
		slog.Error(msg, n...)
	}
}
