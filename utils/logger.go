package utils

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"time"
)

const (
	LevelError = iota
	LevelWarning
	LevelInformational
	LevelDebug
)

var logger *Logger

// Logger 日志（非文件）
type Logger struct {
	level        int8
	infoPrint    *color.Color
	debugPrint   *color.Color
	errorPrint   *color.Color
	warningPrint *color.Color
}

// BuildLogger 根据日志级别构建 Logger
func BuildLogger(level string) {
	logger = &Logger{
		level: func() int8 {
			switch level {
			case "error":
				return LevelError
			case "warning":
				return LevelWarning
			case "info":
				return LevelInformational
			case "debug":
				return LevelDebug
			default:
				return LevelDebug
			}
		}(),
		infoPrint:    color.New(color.BgGreen),
		debugPrint:   color.New(color.BgBlue),
		errorPrint:   color.New(color.BgRed),
		warningPrint: color.New(color.BgYellow),
	}
}

// Log 获取默认 Debug 级别的 Logger，
func Log() *Logger {
	if logger == nil {
		BuildLogger("debug")
	}
	return logger
}

// Println 打印
func (l *Logger) Println(msg string) {
	fmt.Printf("%s %s\n", time.Now().Format("2006-01-02 15:04:05 -0700"), msg)
}

// Panic 极端错误
func (l *Logger) Panic(apiPath, format string, v ...interface{}) {
	if LevelError > l.level {
		return
	}
	msg := fmt.Sprintf("[Panic] APIPath="+apiPath+" "+format, v...)
	l.Println(msg)
	os.Exit(0)
}

// Error 错误
func (l *Logger) Error(apiPath, format string, v ...interface{}) {
	if LevelError > l.level {
		return
	}
	msg := fmt.Sprintf("[Error] APIPath="+apiPath+" "+format, v...)
	l.errorPrint.Printf("%s %s\n", time.Now().Format("2006-01-02 15:04:05 -0700"), msg)
}

// Warning 警告
func (l *Logger) Warning(apiPath, format string, v ...interface{}) {
	if LevelWarning > l.level {
		return
	}
	msg := fmt.Sprintf("[Warning] APIPath="+apiPath+" "+format, v...)
	l.warningPrint.Printf("%s %s\n", time.Now().Format("2006-01-02 15:04:05 -0700"), msg)
}

// Info 信息
func (l *Logger) Info(apiPath, format string, v ...interface{}) {
	if LevelInformational > l.level {
		return
	}
	msg := fmt.Sprintf("[Info] APIPath="+apiPath+" "+format, v...)
	l.infoPrint.Printf("%s %s\n", time.Now().Format("2006-01-02 15:04:05 -0700"), msg)
}

// Debug 调试
func (l *Logger) Debug(apiPath, format string, v ...interface{}) {
	if LevelDebug > l.level {
		return
	}
	msg := fmt.Sprintf("[Debug] APIPath="+apiPath+" "+format, v...)
	l.debugPrint.Printf("%s %s\n", time.Now().Format("2006-01-02 15:04:05 -0700"), msg)
}
