package downloader

import (
	"errors"
	"runtime"
	"strings"
)

const traceDeep = 3

func trace() (name string, err error) {
	pc := make([]uintptr, traceDeep)
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	if len(pc) != traceDeep {
		return "", errors.New("trace error")
	}
	if !strings.Contains(f.Name(), ".") {
		return "", errors.New("trace error")
	}
	list := strings.Split(f.Name(), ".")
	name = list[len(list)-1]
	return
}
