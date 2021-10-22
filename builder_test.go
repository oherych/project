package project

import (
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestNewBuilder(t *testing.T) {
	testTemplate := template.Must(template.New("test").Parse("test file"))
	testData := "test Data"

	tests := map[string]struct {
		fn  func(t *testing.T, b *Builder)
		exp []Job
	}{
		"empty": {
			fn:  func(t *testing.T, b *Builder) {},
			exp: nil,
		},
		"traversing": {
			fn: func(t *testing.T, b *Builder) {
				b.Touch("empty_file.txt", 0651)
				b.DirIn("dir_1", 0652)
				b.String("string_file.txt", "string content", 0653)
				b.Back()
				b.Dir("dir_2", 0654)
				b.In("dir_1")
				b.Bytes("bytes_file.txt", []byte{12, 14, 15}, 0655)
				b.Template("from_template.txt", 0656, testTemplate, testData)
				b.Root()
				b.Touch("in_roo.txt", 0657)

			},
			exp: []Job{
				CreateFile{Path: "empty_file.txt", Data: []byte{}, Perm: 0651},
				CreateFolder{Path: "dir_1", Perm: 0652},
				CreateFile{Path: "dir_1/string_file.txt", Data: []byte("string content"), Perm: 0653},
				CreateFolder{Path: "dir_2", Perm: 0654},
				CreateFile{Path: "dir_1/bytes_file.txt", Data: []byte{12, 14, 15}, Perm: 0655},
				CreateFileFromTemplate{Path: "dir_1/from_template.txt", Template: testTemplate, Data: testData, Perm: 0656},
				CreateFile{Path: "in_roo.txt", Data: []byte{}, Perm: 0657},
			},
		},
	}

	for testCase, tt := range tests {
		t.Run(testCase, func(t *testing.T) {
			b := NewBuilder()

			tt.fn(t, b)

			assert.Equal(t, tt.exp, b.Done())
		})
	}
}

func TestBuilder_Back(t *testing.T) {
	b := NewBuilder()

	assert.False(t, b.Back())

	assert.Equal(t, []string(nil), b.root)

	b.DirIn("level1", 0)

	assert.Equal(t, []string{"level1"}, b.root)

	b.Dir("level2", 0)

	assert.Equal(t, []string{"level1"}, b.root)

	b.In("level2")

	assert.Equal(t, []string{"level1", "level1/level2"}, b.root)

	assert.True(t, b.Back())
	assert.Equal(t, []string{"level1"}, b.root)

	assert.True(t, b.Back())

	assert.Equal(t, []string{}, b.root)

	assert.False(t, b.Back())

	assert.Equal(t, []string{}, b.root)
}
