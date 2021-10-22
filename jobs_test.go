package project

import (
	"os"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestJobs(t *testing.T) {
	testPath := "/test/path"
	testPerm := os.FileMode(0654)
	testTemplate := template.Must(template.New("test").Parse("section {{.}}"))

	tests := map[string]struct {
		job        Job
		expComment string
		expError   error
		expAfter   FileSystemMock
	}{
		"CreateFile": {
			job:        CreateFile{Path: testPath, Data: []byte("test content"), Perm: testPerm},
			expComment: "create file: /test/path",
			expAfter: FileSystemMock{
				{Type: "CreateFile", Name: "/test/path", Data: []byte("test content"), Perm: testPerm},
			},
		},
		"CreateFolder": {
			job:        CreateFolder{Path: testPath, Perm: testPerm},
			expComment: "create folder: /test/path",
			expAfter: FileSystemMock{
				{Type: "CreateDir", Name: testPath, Perm: testPerm},
			},
		},
		"CreateFileFromTemplate": {
			job:        CreateFileFromTemplate{Path: testPath, Perm: testPerm, Template: testTemplate, Data: 234},
			expComment: "create file: /test/path",
			expAfter: FileSystemMock{
				{Type: "CreateFile", Name: testPath, Data: []byte("section 234"), Perm: testPerm},
			},
		},
	}

	for testCase, tt := range tests {
		t.Run(testCase, func(t *testing.T) {
			var mock FileSystemMock

			assert.Equal(t, tt.expError, tt.job.Up(&mock))
			assert.Equal(t, tt.expComment, tt.job.Comment())
			assert.Equal(t, tt.expAfter, mock)
		})
	}
}

type FileSystemMock []TestEvent

type TestEvent struct {
	Type string
	Name string
	Data []byte
	Perm os.FileMode
}

func (f *FileSystemMock) CreateDir(name string, perm os.FileMode) error {
	*f = append(*f, TestEvent{Type: "CreateDir", Name: name, Perm: perm})

	return nil
}

func (f *FileSystemMock) CreateFile(name string, data []byte, perm os.FileMode) error {
	*f = append(*f, TestEvent{Type: "CreateFile", Name: name, Perm: perm, Data: data})

	return nil
}
