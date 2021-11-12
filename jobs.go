package project

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

// Job interface
type Job interface {
	Up(f DriverInterface) error
	Comment() string
}

type CreateFile struct {
	Path string
	Data []byte
	Perm os.FileMode
}

func (j CreateFile) Up(f DriverInterface) error {
	return f.CreateFile(j.Path, j.Data, j.Perm)
}

func (j CreateFile) Comment() string {
	return fmt.Sprintf("create file: %s", j.Path)
}

type CreateFolder struct {
	Path string
	Perm os.FileMode
}

func (j CreateFolder) Up(f DriverInterface) error {
	return f.CreateDir(j.Path, j.Perm)
}

func (j CreateFolder) Comment() string {
	return fmt.Sprintf("create folder: %s", j.Path)
}

type CreateFileFromTemplate struct {
	Path     string
	Template *template.Template
	Data     interface{}
	Perm     os.FileMode
}

func (j CreateFileFromTemplate) Up(f DriverInterface) error {
	var b bytes.Buffer

	if err := j.Template.Execute(&b, j.Data); err != nil {
		return err
	}

	return f.CreateFile(j.Path, b.Bytes(), j.Perm)
}

func (j CreateFileFromTemplate) Comment() string {
	return fmt.Sprintf("create file: %s", j.Path)
}
