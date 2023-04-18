package utils

import (
	"io/fs"
	"net/http"
	"os"
	"testing"
)

func Test_ReadFile(t *testing.T) {
	t.Parallel()

	testFS := http.FS(os.DirFS(".github/tests"))
	file, err := ReadFile("john.txt", testFS)

	AssertEqual(t, string(file), "doe\n")
	AssertEqual(t, err, nil)
}

func Test_Walk(t *testing.T) {
	t.Parallel()

	type file struct {
		path  string
		name  string
		isDir bool
	}
	var files []file

	neededResults := []file{
		{
			path:  "example",
			name:  "example",
			isDir: true,
		},
		{
			path:  "example/example1.txt",
			name:  "example1.txt",
			isDir: false,
		},
		{
			path:  "john.txt",
			name:  "john.txt",
			isDir: false,
		},
	}

	testFS := http.FS(os.DirFS(".github/tests"))
	err := Walk(testFS, ".", func(path string, info fs.FileInfo, err error) error {
		if path != "." {
			files = append(files, file{
				path:  path,
				name:  info.Name(),
				isDir: info.IsDir(),
			})
		}

		return nil
	})

	AssertEqual(t, err, nil)
	AssertEqual(t, files, neededResults)
}
