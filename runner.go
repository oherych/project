package project

import (
	"errors"
	"fmt"
)

const (
	defaultPos = -1
)

var (
	// ErrEndOfJobs says there is no more jobs in the list
	ErrEndOfJobs = errors.New("end of jobs")
)

type Runner struct {
	jobs       []Job
	fileSystem DriverInterface
	output     OutputInterface
	pos        int
}

// NewRunner create new NewRunner
func NewRunner(jobs []Job, fileSystem DriverInterface, output OutputInterface) *Runner {
	return &Runner{
		jobs:       jobs,
		fileSystem: fileSystem,
		output:     output,
		pos:        defaultPos,
	}
}

// Run all jobs one by one
// can return RunError
func (r *Runner) Run() error {
	for {
		err := r.Next()
		if err == ErrEndOfJobs {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// Next run a next job in the list
// can return ErrEndOfJobs
// can return RunError
func (r *Runner) Next() error {
	if r.pos >= len(r.jobs)-1 {
		return ErrEndOfJobs
	}

	r.pos++

	return r.exec(r.jobs[r.pos])
}

// Reset runner to default state
func (r *Runner) Reset() {
	r.pos = defaultPos
}

// Pos return number of curren job
//  0 is first job
// -1 is default value
func (r Runner) Pos() int {
	return r.pos
}

func (r *Runner) exec(j Job) (derr error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		e, ok := r.(error)
		if !ok {
			e = fmt.Errorf("panic: %v", r)
		}

		derr = RunError{Job: j, E: e}
	}()

	if err := j.Up(r.fileSystem); err != nil {
		r.output.Error(j.Comment(), err.Error())

		return RunError{Job: j, E: err}
	}

	r.output.Success(j.Comment())

	return nil
}
