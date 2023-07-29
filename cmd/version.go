package cmd

import "runtime/debug"

var version string

func init() {
	// version
	if version == "" {
		if info, ok := debug.ReadBuildInfo(); ok {
			version = info.Main.Version
		}
	}
	rootCmd.Version = version
}
