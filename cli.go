package main

import (
	"flag"
)

// FlagsCLI store cli flags after parse
type FlagsCLI struct {
	ConfigPath  string
	DebugOutput bool
}

// NewFlagsCLI parse cli args and return FlagsCLI struct
func NewFlagsCLI() *FlagsCLI {
	flagConfig := flag.String("config", "", "path to config file")
	flagDebug := flag.Bool("debug", false, "addition output to stderr")
	flag.Parse()

	return &FlagsCLI{
		ConfigPath:  *flagConfig,
		DebugOutput: *flagDebug,
	}
}
