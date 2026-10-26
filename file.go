package utils

import (
	"io"
	"net/http"
	"os"
	pathpkg "path"
	"path/filepath"
	"slices"
)

// Walk walks the filesystem rooted at root, calling walkFn for each file or
// directory in the filesystem, including root. All errors that arise visiting files
// and directories are filtered by walkFn. The files are walked in lexical
// order.
//
// Walk follows path/filepath.Walk step for step, on an http.FileSystem: walkFn
// is called once for a directory, with the error from reading it if there was
// one; returning filepath.SkipDir skips the directory (the parent directory of
// a file) and filepath.SkipAll ends the walk, neither being reported as an
// error; any other error ends the walk and is returned. As in filepath.Walk,
// the two sentinels must be returned as they are: a wrapped one is an ordinary
// error.
//
// An http.FileSystem has no Lstat, and its files no Readdirnames, so an entry
// is statted by opening it and a directory is listed with Readdir, which stats
// every entry as well. A directory that cannot be opened is therefore reported
// with a nil FileInfo, and one holding an entry that cannot be statted with
// the error of its listing, where filepath.Walk would report each entry.
//
// Without Lstat, Walk also follows symbolic links, which filepath.Walk reports but does
// not follow. A link to a directory is walked as that directory, and a link to one of its
// own ancestors repeats the tree below it until the operating system refuses the path
// (ELOOP); with several such links that takes practically forever. To stop it, return
// filepath.SkipDir for the link. Opening each entry also means that a named pipe or
// another special file can block the walk until the system lets the open return. fs and
// walkFn must not be nil, and a file whose Stat returns neither a FileInfo nor an error
// panics the walk.
func Walk(fs http.FileSystem, root string, walkFn filepath.WalkFunc) error {
	info, err := stat(fs, root)
	if err != nil {
		err = walkFn(root, nil, err)
	} else {
		err = walkInternal(fs, root, info, walkFn)
	}
	if err == filepath.SkipDir || err == filepath.SkipAll { //nolint:errorlint,err113 // filepath.Walk compares the sentinels exactly
		return nil
	}
	return err
}

// ReadFile returns the raw content of a file
func ReadFile(path string, fs http.FileSystem) ([]byte, error) {
	if fs != nil {
		file, err := fs.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close() //nolint: errcheck // No need to check error
		return io.ReadAll(file)
	}
	return os.ReadFile(path) // #nosec G304
}

// readDirNames reads the directory named by dirname and returns
// a sorted list of directory entries.
func readDirNames(fs http.FileSystem, dirname string) ([]string, error) {
	fis, err := readDir(fs, dirname)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(fis))
	for i := range fis {
		names[i] = fis[i].Name()
	}
	slices.Sort(names)
	return names, nil
}

// walkInternal recursively descends path, calling walkFn. It is path/filepath's
// walk: it hands filepath.SkipDir and filepath.SkipAll back to the caller, except
// that a directory's own SkipDir ends there, as Walk would have it end.
func walkInternal(fs http.FileSystem, path string, info os.FileInfo, walkFn filepath.WalkFunc) error {
	if !info.IsDir() {
		return walkFn(path, info, nil)
	}

	// The directory is read before walkFn sees it, so that a failure to read it
	// is reported by that one call.
	names, err := readDirNames(fs, path)
	err1 := walkFn(path, info, err)
	// If err != nil the directory cannot be walked into; if err1 != nil walkFn
	// wants to skip it or to stop. Either way what walkFn returned goes up.
	if err != nil || err1 != nil {
		if err1 == filepath.SkipDir { //nolint:errorlint,err113 // filepath.Walk compares the sentinels exactly
			return nil
		}
		return err1
	}

	for _, name := range names {
		filename := pathpkg.Join(path, name)
		fileInfo, err := stat(fs, filename)
		if err != nil {
			if err := walkFn(filename, fileInfo, err); err != nil && err != filepath.SkipDir { //nolint:errorlint,err113 // filepath.Walk compares the sentinels exactly
				return err
			}
		} else {
			err = walkInternal(fs, filename, fileInfo, walkFn)
			if err != nil {
				if !fileInfo.IsDir() || err != filepath.SkipDir { //nolint:errorlint,err113 // filepath.Walk compares the sentinels exactly
					return err
				}
			}
		}
	}
	return nil
}

// readDir reads the contents of the directory associated with file and
// returns a slice of FileInfo values in directory order.
func readDir(fs http.FileSystem, name string) ([]os.FileInfo, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close() //nolint: errcheck // No need to check error
	return f.Readdir(0)
}

// stat returns the FileInfo structure describing file.
func stat(fs http.FileSystem, name string) (os.FileInfo, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close() //nolint: errcheck // No need to check error
	return f.Stat()
}
