package common

import (
	"fmt"
	"os"
)

type Logger interface {
	Info(a ...interface{})
	Warn(a ...interface{})
	Error(a ...interface{})
}

var Log Logger = &ConsoleLogger{}



type ConsoleLogger struct {

}

func (lg *ConsoleLogger) Info(a ...interface{}) {
	fmt.Fprintf(os.Stdout, "INFO:| %v\n", a)
}

func (lg *ConsoleLogger) Warn(a ...interface{}) {
	fmt.Fprintf(os.Stdout, "WARN:| %v\n", a)
}
func (lg *ConsoleLogger) Error(a ...interface{}) {
	// fmt.Printf(a[0].(string) + "\n", a[1:])
	fmt.Fprintf(os.Stderr, "ERROR:| %v\n", fmt.Sprintf(a[0].(string), a[1:]))
	panic("")
}