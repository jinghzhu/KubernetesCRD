package panic

import (
	"fmt"
	"runtime"
)

const (
	stackBuffer        = 10000
	maxCallerLevel     = 4
	defaultCallerLevel = 1
)

func CatchPanic(err *error) {
	buf := make([]byte, stackBuffer)
	runtime.Stack(buf, false)
	name, file, line := GetCallerInfo(2)
	if r := recover(); r != nil {
		fmt.Printf("%s %s ln%d: PANIC Defered : %v\n", name, file, line, r)
		fmt.Printf("%s %s ln%d: Stack Trace : %s", name, file, line, string(buf))
		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	} else if err != nil && *err != nil {
		fmt.Printf("%s %s ln%d: ERROR : %v\n", name, file, line, *err)
		fmt.Printf("%s %s ln%d: Stack Trace : %s", name, file, line, string(buf))
	}
}

func GetCallerInfo(level int) (caller string, fileName string, lineNum int) {
	if level < 1 || level > maxCallerLevel {
		level = defaultCallerLevel
	}

	pc, file, line, ok := runtime.Caller(level)
	fileDefault := ""
	lineDefault := -1
	nameDefault := ""
	if ok {
		fileDefault = file
		lineDefault = line
	}
	details := runtime.FuncForPC(pc)
	if details != nil {
		nameDefault = details.Name()
	}

	return nameDefault, fileDefault, lineDefault
}
