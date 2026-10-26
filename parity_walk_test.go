package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Walk promises path/filepath.Walk's contract on an http.FileSystem, so these
// tests hold it to the stdlib: the same tree is walked by filepath.Walk and by
// Walk under every combination of "which entry triggers" and "what walkFn
// returns there", and the traces and returned errors must be identical.
//
// A real filesystem cannot be made to fail a directory read or a stat on
// demand (and a test running as root cannot even lose permissions), so for
// those failures the oracle is refWalk, a transcription of filepath's walk.
// refWalk earns that role by agreeing with filepath.Walk on every healthy case
// first (Test_Walk_RefWalkMatchesFilepath).

var (
	errWalkBoom    = errors.New("walkFn boom")
	errWalkReaddir = errors.New("readdir failed")
	errWalkOpen    = errors.New("open failed")
)

// walkEvent is one walkFn call, reduced to what the contract specifies.
type walkEvent struct {
	Path   string
	Name   string
	IsDir  bool
	NoInfo bool
	HasErr bool
}

// walkFixture builds the tree every walk test uses and returns its root.
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

// walkTreePaths lists every slash-separated path of the fixture, "." for the
// root, followed by one that does not exist.
var walkTreePaths = []string{
	".", "a", "a/a1.txt", "a/a2.txt", "a/sub", "a/sub/s1.txt", "a/sub/s2.txt",
	"b", "b/b1.txt", "c.txt", "d", "d/empty", "z.txt", "nope",
}

// walkReturns are what a walkFn returns at its trigger path. Wrapped sentinels
// are ordinary errors to filepath.Walk, which compares them with ==.
var walkReturns = map[string]error{
	"SkipDir":         filepath.SkipDir,
	"SkipAll":         filepath.SkipAll,
	"error":           errWalkBoom,
	"wrapped SkipDir": fmt.Errorf("wrapped: %w", filepath.SkipDir),
	"wrapped SkipAll": fmt.Errorf("wrapped: %w", filepath.SkipAll),
}

// walkerFunc walks the fixture from a slash-separated root (relative to the
// fixture, "." for all of it) and reports the paths it visits relative to the
// fixture too.
type walkerFunc func(root string, walkFn func(rel string, info fs.FileInfo, err error) error) error

// traceWalk runs walker with a walkFn returning ret at trigger and nil
// elsewhere, and returns the calls it saw with the error Walk returned.
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

// relPath maps a path handed to walkFn back to a fixture-relative one.
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

// httpRoot is the http.FileSystem spelling of a fixture-relative root.
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

// refWalk is filepath.Walk and its walk helper transcribed to an
// http.FileSystem: the same statements in the same order, with Lstat and the
// directory read going through fsys.
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
	// filepath sorts with slices.Sort as well.
	for i := 1; i < len(names); i++ {
		for j := i; j > 0 && names[j] < names[j-1]; j-- {
			names[j], names[j-1] = names[j-1], names[j]
		}
	}
	return names, nil
}

// requireSameWalk walks root with both walkers, with walkFn returning ret at
// trigger, and requires the same calls and the same result.
func requireSameWalk(t *testing.T, want, got walkerFunc, root, trigger string, ret error) {
	t.Helper()
	wantEvents, wantErr := traceWalk(want, root, trigger, ret)
	gotEvents, gotErr := traceWalk(got, root, trigger, ret)
	require.Equal(t, wantEvents, gotEvents, "root %q, %v at %q", root, ret, trigger)
	require.Equal(t, wantErr, gotErr, "root %q, %v at %q", root, ret, trigger)
}

// forEveryWalk runs fn for the plain walk of each root and then for every
// (root, trigger, return value) combination.
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

// Test_Walk_ControlFlow pins the contract in words, on the fixture, so a
// failure of the sweep above reads as a behavior and not only as a diff.
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

// faultFS makes chosen opens and directory reads fail.
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

// Test_Walk_FilesystemErrors compares Walk with refWalk when reading a
// directory or statting an entry fails, which filepath.Walk reports with one
// walkFn call (the directory, with its info and the read error) or one for
// the entry (with no info), and then carries on with the siblings.
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
