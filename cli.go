package main

import (
	"errors"
	"flag"
	"os"
)

// FlagsCLI store cli flags after parse
type FlagsCLI struct {
	ConfigPath  string
	DebugOutput bool
}

// NewFlagsCLI parse cli args and return FlagsCLI struct
func NewFlagsCLI() *FlagsCLI {
	flagConfig := flag.String("config", "", "path to config file")
	flagDebug := flag.Bool("debug", false, "addition output to syslog")
	flag.Parse()

	return &FlagsCLI{
		ConfigPath:  *flagConfig,
		DebugOutput: *flagDebug,
	}
}

// VerifyPath verifies ConfigPath
func (f *FlagsCLI) VerifyPath() error {
	if len(f.ConfigPath) == 0 {
		return errors.New("path to config file does not passed")
	}

	if _, err := os.Stat(f.ConfigPath); os.IsNotExist(err) {
		return errors.New("file does not exists")
	}

	return nil
}
