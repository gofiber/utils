// +build !linux,!darwin,!windows,!freebsd,!dragonfly,!netbsd,!openbsd

// Copyright (c) 2017, Jeremy Jay
// All rights reserved.
// https://github.com/pbnjay/memory

package memory

func sysTotalMemory() uint64 {
	return 0
}
