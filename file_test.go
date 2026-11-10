package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
	err = walkInternal(testFS, ".", info, func(path string, _ fs.FileInfo, _ error) error {
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
	err = walkInternal(testFS, ".", info, func(path string, _ fs.FileInfo, _ error) error {
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

	err = walkInternal(testFS, ".", info, func(_ string, _ fs.FileInfo, _ error) error {
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
	err = walkInternal(testFS, ".", info, func(_ string, _ fs.FileInfo, err error) error {
		if err != nil {
			return testErr
		}
		return nil
	})
	require.ErrorIs(t, err, testErr)
}

// Test reading a file directly from the OS when no filesystem is provided
func Test_ReadFile_NoFS(t *testing.T) {
	t.Parallel()

	data, err := ReadFile(".github/tests/john.txt", nil)
	require.NoError(t, err)
	require.Contains(t, string(data), "doe")
}

func Test_ReadFile_SubDir(t *testing.T) {
	t.Parallel()

	data, err := ReadFile("example/example1.txt", http.FS(os.DirFS(".github/tests")))
	require.NoError(t, err)
	require.NotNil(t, data)
}

func Test_ReadFile_StdlibParity(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	contents := map[string][]byte{
		"empty.txt":  {},
		"small.txt":  []byte("hello\n"),
		"binary.bin": {0, 1, 2, 0xff, 0xfe, '\n', '\r'},
		"large.bin":  make([]byte, 1<<20+7),
	}
	for i := range contents["large.bin"] {
		contents["large.bin"][i] = byte(i * 31)
	}
	for name, data := range contents {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), data, 0o600))
	}
	require.NoError(t, os.Mkdir(filepath.Join(dir, "sub"), 0o750))

	fsys := http.Dir(dir)
	for name := range contents {
		want, err := os.ReadFile(filepath.Join(dir, name)) // #nosec G304
		require.NoError(t, err)

		fromOS, err := ReadFile(filepath.Join(dir, name), nil)
		require.NoError(t, err)
		require.Equal(t, want, fromOS, name)

		fromFS, err := ReadFile("/"+name, fsys)
		require.NoError(t, err)
		require.Equal(t, want, fromFS, name)
	}

	_, stdErr := os.ReadFile(filepath.Join(dir, "missing")) // #nosec G304
	require.ErrorIs(t, stdErr, fs.ErrNotExist)
	_, err := ReadFile(filepath.Join(dir, "missing"), nil)
	require.ErrorIs(t, err, fs.ErrNotExist)
	_, err = ReadFile("/missing", fsys)
	require.ErrorIs(t, err, fs.ErrNotExist)

	_, stdErr = os.ReadFile(filepath.Join(dir, "sub")) // #nosec G304
	require.Error(t, stdErr)
	_, err = ReadFile(filepath.Join(dir, "sub"), nil)
	require.Error(t, err)
	_, err = ReadFile("/sub", fsys)
	require.Error(t, err)

	empty, err := ReadFile("/empty.txt", fsys)
	require.NoError(t, err)
	require.Empty(t, empty)
}

// Walk promises filepath.Walk's contract, so each test walks one tree with both under every
// (trigger entry, walkFn result) pair and requires equal traces and errors. A real filesystem cannot
// be made to fail a read or a stat on demand (not even as root), so for those the oracle is refWalk,
// which earns the role by matching filepath.Walk on every healthy case (Test_Walk_RefWalkMatchesFilepath).

var (
	errWalkBoom    = errors.New("walkFn boom")
	errWalkReaddir = errors.New("readdir failed")
	errWalkOpen    = errors.New("open failed")
)

type walkEvent struct {
	Path   string
	Name   string
	IsDir  bool
	NoInfo bool
	HasErr bool
}

func walkFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, p := range []string{
		"a/a1.txt", "a/a2.txt", "a/sub/s1.txt", "a/sub/s2.txt",
		"b/b1.txt", "c.txt", "d/empty/", "z.txt",
	} {
		full := filepath.Join(root, filepath.FromSlash(p))
		if strings.HasSuffix(p, "/") {
			require.NoError(t, os.MkdirAll(full, 0o750))
			continue
		}
		require.NoError(t, os.MkdirAll(filepath.Dir(full), 0o750))
		require.NoError(t, os.WriteFile(full, []byte("x"), 0o600))
	}
	return root
}

var walkTreePaths = []string{
	".", "a", "a/a1.txt", "a/a2.txt", "a/sub", "a/sub/s1.txt", "a/sub/s2.txt",
	"b", "b/b1.txt", "c.txt", "d", "d/empty", "z.txt", "nope",
}

// walkReturns: wrapped sentinels are ordinary errors, since filepath.Walk compares with ==.
var walkReturns = map[string]error{
	"SkipDir":         filepath.SkipDir,
	"SkipAll":         filepath.SkipAll,
	"error":           errWalkBoom,
	"wrapped SkipDir": fmt.Errorf("wrapped: %w", filepath.SkipDir),
	"wrapped SkipAll": fmt.Errorf("wrapped: %w", filepath.SkipAll),
}

type walkerFunc func(root string, walkFn func(rel string, info fs.FileInfo, err error) error) error

func traceWalk(walker walkerFunc, root, trigger string, ret error) ([]walkEvent, error) {
	var events []walkEvent
	err := walker(root, func(rel string, info fs.FileInfo, err error) error {
		ev := walkEvent{Path: rel, NoInfo: info == nil, HasErr: err != nil}
		if info != nil {
			ev.IsDir = info.IsDir()
			if rel != "." { // the root is named after the temporary directory
				ev.Name = info.Name()
			}
		}
		events = append(events, ev)
		if rel == trigger {
			return ret
		}
		return nil
	})
	return events, err
}

func relPath(prefix, path string) string {
	rel := strings.TrimPrefix(filepath.ToSlash(path), filepath.ToSlash(prefix))
	rel = strings.TrimPrefix(rel, "/")
	if rel == "" {
		return "."
	}
	return rel
}

func stdWalker(fixture string) walkerFunc {
	return func(root string, walkFn func(string, fs.FileInfo, error) error) error {
		return filepath.Walk(filepath.Join(fixture, filepath.FromSlash(root)), func(path string, info fs.FileInfo, err error) error {
			return walkFn(relPath(fixture, path), info, err)
		})
	}
}

func httpRoot(root string) string {
	if root == "." {
		return "/"
	}
	return "/" + root
}

func utilsWalker(fsys http.FileSystem) walkerFunc {
	return func(root string, walkFn func(string, fs.FileInfo, error) error) error {
		return Walk(fsys, httpRoot(root), func(path string, info fs.FileInfo, err error) error {
			return walkFn(relPath("/", path), info, err)
		})
	}
}

func refWalker(fsys http.FileSystem) walkerFunc {
	return func(root string, walkFn func(string, fs.FileInfo, error) error) error {
		return refWalk(fsys, httpRoot(root), func(path string, info fs.FileInfo, err error) error {
			return walkFn(relPath("/", path), info, err)
		})
	}
}

// refWalk transcribes filepath.Walk statement for statement; Lstat and the directory read go through fsys.
func refWalk(fsys http.FileSystem, root string, walkFn filepath.WalkFunc) error {
	info, err := refStat(fsys, root)
	if err != nil {
		err = walkFn(root, nil, err)
	} else {
		err = refWalkDir(fsys, root, info, walkFn)
	}
	if err == filepath.SkipDir || err == filepath.SkipAll { //nolint:errorlint,err113 // transcribed from filepath.Walk
		return nil
	}
	return err
}

func refWalkDir(fsys http.FileSystem, path string, info fs.FileInfo, walkFn filepath.WalkFunc) error {
	if !info.IsDir() {
		return walkFn(path, info, nil)
	}

	names, err := refReadDirNames(fsys, path)
	err1 := walkFn(path, info, err)
	if err != nil || err1 != nil {
		return err1
	}

	for _, name := range names {
		filename := strings.TrimSuffix(path, "/") + "/" + name
		fileInfo, err := refStat(fsys, filename)
		if err != nil {
			if err := walkFn(filename, fileInfo, err); err != nil && err != filepath.SkipDir { //nolint:errorlint,err113 // transcribed from filepath.Walk
				return err
			}
		} else {
			err = refWalkDir(fsys, filename, fileInfo, walkFn)
			if err != nil {
				if !fileInfo.IsDir() || err != filepath.SkipDir { //nolint:errorlint,err113 // transcribed from filepath.Walk
					return err
				}
			}
		}
	}
	return nil
}

func refStat(fsys http.FileSystem, name string) (fs.FileInfo, error) {
	f, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close() //nolint:errcheck // read-only handle
	return f.Stat()
}

func refReadDirNames(fsys http.FileSystem, name string) ([]string, error) {
	f, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close() //nolint:errcheck // read-only handle
	infos, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(infos))
	for _, info := range infos {
		names = append(names, info.Name())
	}
	for i := 1; i < len(names); i++ {
		for j := i; j > 0 && names[j] < names[j-1]; j-- {
			names[j], names[j-1] = names[j-1], names[j]
		}
	}
	return names, nil
}

func requireSameWalk(t *testing.T, want, got walkerFunc, root, trigger string, ret error) {
	t.Helper()
	wantEvents, wantErr := traceWalk(want, root, trigger, ret)
	gotEvents, gotErr := traceWalk(got, root, trigger, ret)
	require.Equal(t, wantEvents, gotEvents, "root %q, %v at %q", root, ret, trigger)
	require.Equal(t, wantErr, gotErr, "root %q, %v at %q", root, ret, trigger)
}

func forEveryWalk(fn func(root, trigger string, ret error)) {
	for _, root := range []string{".", "a", "a/sub", "c.txt", "d/empty", "nope"} {
		fn(root, "", nil)
		for _, trigger := range walkTreePaths {
			for _, ret := range walkReturns {
				fn(root, trigger, ret)
			}
		}
	}
}

func Test_Walk_RefWalkMatchesFilepath(t *testing.T) {
	t.Parallel()
	fixture := walkFixture(t)
	fsys := http.Dir(fixture)
	forEveryWalk(func(root, trigger string, ret error) {
		requireSameWalk(t, stdWalker(fixture), refWalker(fsys), root, trigger, ret)
	})
}

func Test_Walk_MatchesFilepathWalk(t *testing.T) {
	t.Parallel()
	fixture := walkFixture(t)
	fsys := http.Dir(fixture)
	forEveryWalk(func(root, trigger string, ret error) {
		requireSameWalk(t, stdWalker(fixture), utilsWalker(fsys), root, trigger, ret)
	})
}

// Test_Walk_ControlFlow pins the contract in words, so a sweep failure reads as a behavior, not a diff.
func Test_Walk_ControlFlow(t *testing.T) {
	t.Parallel()
	fixture := walkFixture(t)
	walker := utilsWalker(http.Dir(fixture))
	paths := func(events []walkEvent) []string {
		out := make([]string, len(events))
		for i, ev := range events {
			out[i] = ev.Path
		}
		return out
	}

	t.Run("SkipDir on a directory skips only that directory", func(t *testing.T) {
		t.Parallel()
		events, err := traceWalk(walker, ".", "a", filepath.SkipDir)
		require.NoError(t, err)
		require.Equal(t, []string{".", "a", "b", "b/b1.txt", "c.txt", "d", "d/empty", "z.txt"}, paths(events))
	})
	t.Run("SkipDir on a file skips the rest of its directory", func(t *testing.T) {
		t.Parallel()
		events, err := traceWalk(walker, ".", "a/a1.txt", filepath.SkipDir)
		require.NoError(t, err)
		require.Equal(t, []string{".", "a", "a/a1.txt", "b", "b/b1.txt", "c.txt", "d", "d/empty", "z.txt"}, paths(events))
	})
	t.Run("SkipDir on a file directly under the root is not an error", func(t *testing.T) {
		t.Parallel()
		events, err := traceWalk(walker, ".", "c.txt", filepath.SkipDir)
		require.NoError(t, err)
		require.Equal(t, []string{".", "a", "a/a1.txt", "a/a2.txt", "a/sub", "a/sub/s1.txt", "a/sub/s2.txt", "b", "b/b1.txt", "c.txt"}, paths(events))
	})
	t.Run("SkipAll ends the walk without an error", func(t *testing.T) {
		t.Parallel()
		events, err := traceWalk(walker, ".", "a/sub", filepath.SkipAll)
		require.NoError(t, err)
		require.Equal(t, []string{".", "a", "a/a1.txt", "a/a2.txt", "a/sub"}, paths(events))
	})
	t.Run("another error ends the walk and is returned", func(t *testing.T) {
		t.Parallel()
		events, err := traceWalk(walker, ".", "b/b1.txt", errWalkBoom)
		require.ErrorIs(t, err, errWalkBoom)
		require.Equal(t, "b/b1.txt", events[len(events)-1].Path)
	})
	t.Run("a wrapped SkipDir is just an error", func(t *testing.T) {
		t.Parallel()
		wrapped := fmt.Errorf("wrapped: %w", filepath.SkipDir)
		_, err := traceWalk(walker, ".", "a", wrapped)
		require.Equal(t, wrapped, err)
	})
	t.Run("a missing root reaches walkFn with its error and no info", func(t *testing.T) {
		t.Parallel()
		events, err := traceWalk(walker, "nope", "", nil)
		require.NoError(t, err)
		require.Equal(t, []walkEvent{{Path: "nope", NoInfo: true, HasErr: true}}, events)
	})
}

type faultFS struct {
	http.FileSystem
	openFail    map[string]bool
	readdirFail map[string]bool
}

type faultFile struct {
	http.File
	readdirErr error
}

func (f faultFS) Open(name string) (http.File, error) {
	clean := strings.TrimPrefix(filepath.ToSlash(filepath.Clean(name)), "/")
	if clean == "" {
		clean = "."
	}
	if f.openFail[clean] {
		return nil, errWalkOpen
	}
	file, err := f.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	if f.readdirFail[clean] {
		return faultFile{File: file, readdirErr: errWalkReaddir}, nil
	}
	return file, nil
}

func (f faultFile) Readdir(count int) ([]fs.FileInfo, error) {
	if f.readdirErr != nil {
		return nil, f.readdirErr
	}
	return f.File.Readdir(count)
}

// Test_Walk_FilesystemErrors compares Walk with refWalk when a directory read or a stat fails.
func Test_Walk_FilesystemErrors(t *testing.T) {
	t.Parallel()
	fixture := walkFixture(t)
	faults := map[string]faultFS{
		"unreadable directory":         {readdirFail: map[string]bool{"a": true}},
		"unreadable root":              {readdirFail: map[string]bool{".": true}},
		"unreadable nested directory":  {readdirFail: map[string]bool{"a/sub": true}},
		"unreadable empty directory":   {readdirFail: map[string]bool{"d/empty": true}},
		"unstatable file":              {openFail: map[string]bool{"a/a2.txt": true}},
		"unstatable directory":         {openFail: map[string]bool{"b": true}},
		"unstatable root":              {openFail: map[string]bool{".": true}},
		"several at once":              {readdirFail: map[string]bool{"a/sub": true, "b": true}, openFail: map[string]bool{"c.txt": true, "z.txt": true}},
		"unreadable and unstatable":    {readdirFail: map[string]bool{"a": true}, openFail: map[string]bool{"d": true}},
		"unreadable with failing open": {readdirFail: map[string]bool{"d": true}, openFail: map[string]bool{"d/empty": true}},
	}
	for name, fault := range faults {
		fault.FileSystem = http.Dir(fixture)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			forEveryWalk(func(root, trigger string, ret error) {
				requireSameWalk(t, refWalker(fault), utilsWalker(fault), root, trigger, ret)
			})
		})
	}

	t.Run("a directory is reported once, with its read error", func(t *testing.T) {
		t.Parallel()
		fault := faultFS{FileSystem: http.Dir(fixture), readdirFail: map[string]bool{"a": true}}
		var calls []walkEvent
		err := Walk(fault, "/", func(path string, info fs.FileInfo, err error) error {
			if rel := relPath("/", path); rel == "a" {
				calls = append(calls, walkEvent{Path: rel, IsDir: info.IsDir(), HasErr: err != nil})
				require.ErrorIs(t, err, errWalkReaddir)
			}
			return nil
		})
		require.NoError(t, err)
		require.Equal(t, []walkEvent{{Path: "a", IsDir: true, HasErr: true}}, calls)
	})
}
