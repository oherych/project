package project

// RunError contain information about error and job
type RunError struct {
	Job Job
	E   error
}

// Error to string
func (e RunError) Error() string {
	return e.E.Error()
}
