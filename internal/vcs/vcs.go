package vcs

import (
	"fmt"
	"runtime/debug"
)

func Version() string {

	var (
		revision string
		modified bool
	)

	// Read build metadata when it is embedded by the Go toolchain.
	bi, ok := debug.ReadBuildInfo()
	if ok {
		// Pull out the commit hash and dirty flag recorded at build time.
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				revision = s.Value
			case "vcs.modified":
				modified = s.Value == "true"
			}
		}
	}

	// Mark locally modified builds so the reported version matches the binary.
	if modified {
		return fmt.Sprintf("%s-dirty", revision)
	}
	return revision
}
