package global

import "runtime/debug"

// Version - version of this binary, taken from the build info the Go toolchain
// stamps into every binary. A build from a tagged, clean checkout reports the
// tag, an untagged checkout a pseudo-version, and `go run` reports "(devel)".
func Version() (version string) {
	version = "(unknown)"
	if bi, ok := debug.ReadBuildInfo(); ok && bi.Main.Version != "" {
		version = bi.Main.Version
	}
	return version
}
