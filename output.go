package project

import (
	"fmt"
	"io"
)

var (
	_ OutputInterface = ConsoleOutput{}
)

type OutputInterface interface {
	Success(text string)
	Error(text string, errMsg string)
}

type ConsoleOutput struct {
	Writer io.Writer
}

func (c ConsoleOutput) Success(text string) {
	fmt.Fprintln(c.Writer, "√", text)
}

func (c ConsoleOutput) Error(text string, errMsg string) {
	fmt.Fprintln(c.Writer, "X", text)
	fmt.Fprintln(c.Writer, "Err:", errMsg)
}
