package logger

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	logsyslog "log/syslog"

	"github.com/sirupsen/logrus"
	syslog "github.com/sirupsen/logrus/hooks/syslog"
)

type Fields map[string]interface{}

type LoggingRemoteOpts struct {
	RemoteProtocol string             `json:"protocol"`
	RemoteServer   string             `json:"remote_server"`
	Flag           int                `json:"flag"`
	Priority       logsyslog.Priority `json:"priority"`
	Tag            string             `json:"tag"`
}

// Only support three log levels
const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	ERROR = "ERROR"
)

var (
	LogLevel string
	console  *log.Logger
)

func init() {
	initConsole()
	initLogrus()
}

func initLogrus() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// Disable logrun to output logs in local machine
	logrus.SetOutput(ioutil.Discard)

	if os.Getenv("HTTPLOGSERVER") != "" {
		logopts := &LoggingRemoteOpts{
			RemoteProtocol: "tcp",
			RemoteServer:   os.Getenv("HTTPLOGSERVER"),
			Flag:           log.LstdFlags,
			Priority:       logsyslog.LOG_INFO,
		}
		hook, err := syslog.NewSyslogHook(
			logopts.RemoteProtocol,
			logopts.RemoteServer,
			logopts.Priority,
			logopts.Tag)
		if err != nil {
			fmt.Printf("Could not format logger %+v\n", err)
		} else {
			logrus.AddHook(hook)
		}
	}
}

func initConsole() {
	console = GetConsoleLogger()
}

// SetLogLevel sets log level for logrus and local console. Accept info and error level.
// By default, it is info level.
func SetLogLevel(level string) {
	if strings.ToUpper(level) == INFO || !isValidLevel(level) {
		// by default, info level
		logrus.SetLevel(logrus.InfoLevel)
		LogLevel = INFO
	} else {
		logrus.SetLevel(logrus.ErrorLevel)
		LogLevel = ERROR
	}
	msg := "The log level is " + LogLevel
	logrus.Infoln(msg)
	console.Println(msg)
}

func isValidLevel(level string) bool {
	level = strings.ToUpper(level)
	if level != INFO && level != ERROR {
		return false
	}
	return true
}

// Print log in simple way
func Info(msg interface{}) {
	if LogLevel == INFO {
		output(msg, INFO)
	}
}

func Error(msg interface{}) {
	output(msg, ERROR)
}

func output(msg interface{}, prefix string) {
	logrus.Println(msg)
	file, line := Locate(3)
	console.Println(
		fmt.Sprintf("[%s] %s Ln%d %+v", prefix, file, line, msg),
	)
}

// InfoFields prints log with fields
func InfoFields(msg string, fields Fields) {
	if LogLevel == INFO {
		outputFields(msg, fields, INFO)
	}
}

func ErrorFields(msg string, fields Fields) {
	outputFields(msg, fields, ERROR)
}

func outputFields(msg string, fields Fields, prefix string) {
	e := logrus.WithFields(logrus.Fields(fields))
	e.Time = time.Now()
	e.Println(msg)
	data, err := e.String()
	if err != nil {
		console.Println("Fail to get string representation for logrus.entry. " + err.Error())
		data = fmt.Sprintf("%v", fields)
	}
	file, line := Locate(3)
	console.Println(fmt.Sprintf("[%s] %s Ln%d %s %s", prefix, file, line, msg, data))
}

func Locate(skip int) (filename string, line int) {
	if skip < 0 {
		return "", skip
	}
	_, path, line, ok := runtime.Caller(skip)
	file := ""
	if ok {
		_, file = filepath.Split(path)
	} else {
		fmt.Println("Fail to get method caller")
		line = -1
	}
	return file, line
}
