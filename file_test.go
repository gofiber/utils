package utils

import (
	"fmt"
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
		require.Equal(t, "doe\r\n", string(file))
	default:
		require.Equal(t, "doe\n", string(file))
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

	expectedResults := []file{
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
	err := Walk(testFS, ".", func(path string, info fs.FileInfo, _ error) error {
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
	require.Equal(t, expectedResults, files)
}

func Test_Walk_Error(t *testing.T) {
	t.Parallel()

	testFS := http.FS(os.DirFS(".github/tests"))
	err := Walk(testFS, "nonexistent", func(path string, _ fs.FileInfo, _ error) error {
		return fmt.Errorf("file not found: %s", path)
	})

	require.Error(t, err)
}

func Test_ReadFile_Error(t *testing.T) {
	t.Parallel()

	// Test error when file does not exist
	testFS := http.FS(os.DirFS(".github/tests"))
	_, err := ReadFile("nonexistent.txt", testFS)
	require.Error(t, err)

	// Test error when file does not exist and fs is nil
	_, err = ReadFile("nonexistent.txt", nil)
	require.Error(t, err)
}
