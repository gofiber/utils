package utils

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// errFS wraps another http.FileSystem and fails opening a specific path.
type errFS struct {
	http.FileSystem
	fail string
}

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

// Test for readDirNames success and error cases
func Test_readDirNames(t *testing.T) {
	t.Parallel()

	testFS := http.FS(os.DirFS(".github/tests"))

	names, err := readDirNames(testFS, ".")
	require.NoError(t, err)
	require.Equal(t, []string{"example", "john.txt"}, names)

	_, err = readDirNames(testFS, "does-not-exist")
	require.Error(t, err)
}

// Test for readDir error handling and content retrieval
func Test_readDir(t *testing.T) {
	t.Parallel()

	testFS := http.FS(os.DirFS(".github/tests"))

	fis, err := readDir(testFS, ".")
	require.NoError(t, err)
	require.Len(t, fis, 2)

	_, err = readDir(testFS, "missing")
	require.Error(t, err)
}

// Test walkInternal when walkFn returns SkipDir
func Test_walkInternal_SkipDir(t *testing.T) {
	testFS := http.FS(os.DirFS(".github/tests"))
	info, err := stat(testFS, ".")
	require.NoError(t, err)

	var visited []string
	err = walkInternal(testFS, ".", info, func(path string, info fs.FileInfo, err error) error {
		visited = append(visited, path)
		return filepath.SkipDir
	})
	require.NoError(t, err)
	require.Equal(t, []string{"."}, visited)
}

// Test walkInternal returning error from walkFn
func Test_walkInternal_Error(t *testing.T) {
	testFS := http.FS(os.DirFS(".github/tests"))
	info, err := stat(testFS, ".")
	require.NoError(t, err)

	testErr := errors.New("walk error")
	err = walkInternal(testFS, ".", info, func(path string, info fs.FileInfo, err error) error {
		if path != "." {
			return testErr
		}
		return nil
	})
	require.ErrorIs(t, err, testErr)
}

func (e errFS) Open(name string) (http.File, error) {
	if filepath.Clean(name) == e.fail {
		return nil, os.ErrNotExist
	}
	return e.FileSystem.Open(name)
}

// Test walkInternal when stat returns error for a child and walkFn returns nil
func Test_walkInternal_StatError(t *testing.T) {
	testFS := errFS{http.FS(os.DirFS(".github/tests")), filepath.Join("example", "example1.txt")}
	info, err := stat(testFS, ".")
	require.NoError(t, err)

	err = walkInternal(testFS, ".", info, func(path string, info fs.FileInfo, err error) error {
		return nil
	})
	require.NoError(t, err)
}

// Test walkInternal when stat returns error and walkFn propagates error
func Test_walkInternal_StatErrorPropagate(t *testing.T) {
	testFS := errFS{http.FS(os.DirFS(".github/tests")), filepath.Join("example", "example1.txt")}
	info, err := stat(testFS, ".")
	require.NoError(t, err)

	testErr := errors.New("stat failed")
	err = walkInternal(testFS, ".", info, func(path string, _ fs.FileInfo, err error) error {
		if err != nil {
			return testErr
		}
		return nil
	})
	require.ErrorIs(t, err, testErr)
}
