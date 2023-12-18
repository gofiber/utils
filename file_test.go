package utils

import (
	"io/fs"
	"net/http"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ReadFile(t *testing.T) {
	t.Parallel()

	testFS := http.FS(os.DirFS(".github/tests"))
	file, err := ReadFile("john.txt", testFS)

	switch runtime.GOOS {
	case "windows":
		require.Equal(t, string(file), "doe\r\n")
	default:
		require.Equal(t, string(file), "doe\n")
	}

	require.NoError(t, err)
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

	require.NoError(t, err)
	require.Equal(t, files, neededResults)
}
