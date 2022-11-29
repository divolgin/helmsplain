package log

import "fmt"

var isDebug = false

func SetDebug(debug bool) {
	isDebug = debug
}

func Debugf(format string, args ...interface{}) {
	if !isDebug {
		return
	}
	fmt.Printf(format, args...)
}
