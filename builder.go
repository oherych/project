package project

import (
	"os"
	"path"
	"text/template"
)

// Builder is project structure builder
type Builder struct {
	root []string
	jobs []Job
}

// NewBuilder create new project builder
func NewBuilder() *Builder {
	return &Builder{}
}

// Done return list of jobs
// if list is empty method will return nil
func (b *Builder) Done() []Job {
	return b.jobs
}

// In move cursor into folder
func (b *Builder) In(name string) {
	b.root = append(b.root, b.abs(name))
}

// Back move cursor to previous position
// can panic if no more positions in the stack
func (b *Builder) Back() bool {
	if len(b.root) == 0 {
		return false
	}

	b.root = b.root[:len(b.root)-1]

	return true
}

// Root move cursor to root position
func (b *Builder) Root() {
	b.root = b.root[:0]
}

// Job add custom job into the list of jobs
func (b *Builder) Job(j Job) {
	b.jobs = append(b.jobs, j)
}

// Dir add create directory job to the list of jobs
func (b *Builder) Dir(name string, perm os.FileMode) {
	b.Job(CreateFolder{Path: b.abs(name), Perm: perm})
}

// DirIn add create directory job to the list of jobs
// and move custom just created directory
func (b *Builder) DirIn(name string, perm os.FileMode) {
	b.Dir(name, perm)
	b.In(name)
}

// Bytes add new create file job with provided content
func (b *Builder) Bytes(fileName string, data []byte, perm os.FileMode) {
	b.Job(CreateFile{Path: b.abs(fileName), Data: data, Perm: perm})
}

// String add new create file job with provided content
func (b *Builder) String(fileName string, data string, perm os.FileMode) {
	b.Bytes(fileName, []byte(data), perm)
}

// Touch add new create file job with empty content
func (b *Builder) Touch(fileName string, perm os.FileMode) {
	b.Bytes(fileName, []byte{}, perm)
}

func (b *Builder) Template(fileName string, perm os.FileMode, t *template.Template, data interface{}) {
	b.Job(CreateFileFromTemplate{
		Path:     b.abs(fileName),
		Template: t,
		Data:     data,
		Perm:     perm,
	})
}

func (b *Builder) abs(name string) string {
	if len(b.root) == 0 {
		return name
	}

	return path.Join(b.root[len(b.root)-1], name)
}
