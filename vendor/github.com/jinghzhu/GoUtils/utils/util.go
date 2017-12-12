package utils

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func Struct2String(v interface{}) string {
	result, err := json.Marshal(v)
	if err != nil {
		errMsg := "Fail to translate to json"
		fmt.Println(errMsg)
		return fmt.Sprintf("%v", v)
	}
	return string(result)
}

// GetMountPoints returns all moutpoins in a string array
func GetMountPoints(server string) ([]string, error) {
	b, err := exec.Command("showmount", "-e", server).Output()
	if err != nil {
		fmt.Println("error in showmount: " + err.Error())
		return nil, err
	}
	s := strings.TrimSpace(string(b))
	// The fist line of showmount -e <server> is Exports list on <server>
	firstLine := strings.Index(s, "\n")
	sArr := strings.Split(s[firstLine+1:], "\n")
	for i := 0; i < len(sArr); i++ {
		index := strings.Index(sArr[i], " ")
		temp := sArr[i]
		sArr[i] = temp[:index]
	}
	return sArr, nil
}

// Locate returns the line number and file name in the current goroutine statck trace. The argument skip is the number of stack frames to ascend, with 0 identifying the caller of Caller.
func Locate(skip int) (filename string, line int) {
	if skip < 0 {
		skip = 2
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
