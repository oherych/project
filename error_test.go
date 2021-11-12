package project

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunError_Error(t *testing.T) {
	testBaseError := errors.New("im_error")
	testError := RunError{Job: nil, E: testBaseError}

	assert.Equal(t, testBaseError.Error(), testError.Error())
}
