package project

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunner_Run(t *testing.T) {
	testError := errors.New("im_test_error")
	testJob1 := BuildCustomJob("f1", "comment 1")
	testJob2 := BuildCustomJob("f2", "comment 2")
	testJob3 := BuildCustomJob("f3", "comment 3")
	testJobError := &CustomJob{
		MockUp: func(f DriverInterface) error {
			return testError
		},
		MockComment: func() string {
			return "comment 2"
		},
	}
	testJobStringPanic := &CustomJob{
		MockUp: func(f DriverInterface) error {
			panic("im_panic")
		},
		MockComment: func() string {
			return "comment 2"
		},
	}
	testJobErrorPanic := &CustomJob{
		MockUp: func(f DriverInterface) error {
			panic(testError)
		},
		MockComment: func() string {
			return "comment 2"
		},
	}

	tests := map[string]struct {
		jobs []Job

		expAfter []string
		expErr   error
	}{
		"empty": {
			jobs:   []Job{},
			expErr: nil,
		},
		"one": {
			jobs: []Job{
				BuildCustomJob("f1", "comment 1"),
			},
			expAfter: []string{"f1"},
			expErr:   nil,
		},
		"many": {
			jobs:     []Job{testJob1, testJob2, testJob3},
			expAfter: []string{"f1", "f2", "f3"},
			expErr:   nil,
		},
		"with error": {
			jobs:     []Job{testJob1, testJobError, testJob3},
			expAfter: []string{"f1"},
			expErr:   RunError{Job: testJobError, E: testError},
		},
		"with string panic": {
			jobs:     []Job{testJob1, testJobStringPanic, testJob3},
			expAfter: []string{"f1"},
			expErr:   RunError{Job: testJobStringPanic, E: errors.New("panic: im_panic")},
		},
		"with error panic": {
			jobs:     []Job{testJob1, testJobErrorPanic, testJob3},
			expAfter: []string{"f1"},
			expErr:   RunError{Job: testJobErrorPanic, E: testError},
		},
	}

	for testCase, tt := range tests {
		t.Run(testCase, func(t *testing.T) {
			driver := &DriverMock{}
			output := &OutputMock{}

			runner := NewRunner(tt.jobs, driver, output)

			err := runner.Run()

			assert.Equal(t, tt.expAfter, driver.After)
			assert.Equal(t, tt.expErr, err)
		})
	}

}

func TestBySteps(t *testing.T) {
	driver := &DriverMock{}
	output := &OutputMock{}
	testJobs := []Job{
		BuildCustomJob("f1", "comment 1"),
		BuildCustomJob("f2", "comment 2"),
	}

	runner := NewRunner(testJobs, driver, output)

	assert.Equal(t, defaultPos, runner.Pos())
	assert.Nil(t, runner.Next())

	assert.Equal(t, 0, runner.Pos())
	assert.Nil(t, runner.Next())

	assert.Equal(t, 1, runner.Pos())
	assert.Equal(t, ErrEndOfJobs, runner.Next())

	assert.Equal(t, 1, runner.Pos())
	assert.Equal(t, ErrEndOfJobs, runner.Next())

	runner.Reset()

	assert.Equal(t, defaultPos, runner.Pos())
	assert.Nil(t, runner.Next())

	assert.Equal(t, 0, runner.Pos())
}

type CustomJob struct {
	MockUp      func(f DriverInterface) error
	MockComment func() string
}

func BuildCustomJob(name string, comment string) CustomJob {
	return CustomJob{
		MockUp: func(f DriverInterface) error {
			return f.CreateFile(name, nil, 0)
		},
		MockComment: func() string {
			return comment
		},
	}
}

func (c CustomJob) Up(f DriverInterface) error {
	return c.MockUp(f)
}

func (c CustomJob) Comment() string {
	return c.MockComment()
}

type DriverMock struct {
	After []string
}

func (d *DriverMock) CreateDir(_ string, _ os.FileMode) error {
	return nil
}

func (d *DriverMock) CreateFile(filePath string, _ []byte, _ os.FileMode) error {
	d.After = append(d.After, filePath)
	return nil
}

type OutputMock struct {
}

func (f OutputMock) Success(text string) {

}

func (f OutputMock) Error(text string, errMsg string) {

}
