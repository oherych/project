package project

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestConsoleOutput(t *testing.T) {
	b := strings.Builder{}

	output := ConsoleOutput{Writer: &b}
	output.Success("im_success_text")
	output.Error("im_error_text", "im_error_msg")

	exp := `√ im_success_text
X im_error_text
	Err: im_error_msg
`
	assert.Equal(t, exp, b.String())
}

func TestNOPOutput(t *testing.T) {
	output := NOPOutput{}
	output.Success("im_success_text")
	output.Error("im_error_text", "im_error_msg")
}