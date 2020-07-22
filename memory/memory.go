// Copyright (c) 2017, Jeremy Jay
// All rights reserved.
// https://github.com/pbnjay/memory

// Package memory provides a single method reporting total system memory
// accessible to the kernel.
package memory

import (
	"fmt"
	"runtime"
)

// TotalMemory returns the total accessible system memory in bytes.
//
// The total accessible memory is installed physical memory size minus reserved
// areas for the kernel and hardware, if such reservations are reported by
// the operating system.
//
// If accessible memory size could not be determined, then 0 is returned.
func TotalMemory() uint64 {
	return sysTotalMemory()
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
