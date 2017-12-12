package logger

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/jinghzhu/GoUtils/mail"
)

const (
	monitorInterval = 6 * time.Hour
	logServer       = "test.com"
)

var (
	connLogger *log.Logger
	onceConn   sync.Once
	Conn       *net.TCPConn
)

func GetConnLogger() *log.Logger {
	onceConn.Do(func() {
		Conn, err := getConn()
		if err != nil {
			Conn.Close()
			fmt.Println("Fail to connect to remote log server. " + err.Error())
			return
		}
		go monitorConn(monitorInterval)
		connLogger = log.New(Conn, "", log.Ldate|log.Ltime)
	})
	return connLogger
}

func getConn() (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", logServer)
	if err != nil {
		errMsg := "Error in esablish the connection(ResolveTCPAddr) to log server: " + err.Error()
		fmt.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	var (
		success    = false
		dialErr    error
		conn       *net.TCPConn
		retryRound = 1
	)

	for !success || retryRound <= 3 {
		conn, dialErr = net.DialTCP("tcp", nil, tcpAddr)
		if dialErr != nil {
			fmt.Println("Encounter error in DialTCP: " + dialErr.Error())
			fmt.Printf("Have tried %d rounds. Sleep 30 seconds.\n", retryRound)
			retryRound++
			time.Sleep(time.Second * 30)
		} else {
			success = true
		}
	}

	if !success {
		errMsg := "Error in esablish the connection(DialTCP) to log server: " + err.Error()
		fmt.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	return conn, nil
}

func monitorConn(interval time.Duration) {
	fmt.Printf("Start log server connection monitor. The interval is %s.\n", interval.String())
	for {
		conn, err := getConn()
		defer conn.Close()
		if err != nil {
			GetConsoleLogger().Println("Encounter error while monitoring connection. " + err.Error())
			mail.Send("test", "test@test.com", []string{"test@test.com"}, "alert", "alert")
		} else {
			time.Sleep(interval)
		}
	}
}
