// internal/helpers/system.go

package helpers

import (
	"os"
	"runtime"
)

// GetCurrentPID get PID (max 4294967295)
func GetCurrentPID() uint32 {
	return uint32(os.Getpid())
}

// GetOperatingSystem returns the operating system of the current environment.
func GetOperatingSystem() string {
	return runtime.GOOS
}
