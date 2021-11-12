package project

import (
	"fmt"
	"io"
)

var (
	_ OutputInterface = ConsoleOutput{}
	_ OutputInterface = NOPOutput{}
)

type OutputInterface interface {
	Success(text string)
	Error(text string, errMsg string)
}

type ConsoleOutput struct {
	Writer io.Writer
}

func (c ConsoleOutput) Success(text string) {
	_, _ = fmt.Fprintln(c.Writer, "√", text)
}

func (c ConsoleOutput) Error(text string, errMsg string) {
	_, _ = fmt.Fprintln(c.Writer, "X", text)
	_, _ = fmt.Fprintln(c.Writer, "\tErr:", errMsg)
}

type NOPOutput struct{}

func (NOPOutput) Success(_ string)         {}
func (NOPOutput) Error(_ string, _ string) {}
