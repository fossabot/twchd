package main

import (
	"errors"
	"flag"
	"os"
	"strings"
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

// VerifyPath verifies existance ConfigPath
func (f *FlagsCLI) VerifyPath() error {
	if _, err := os.Stat(f.ConfigPath); os.IsNotExist(err) {
		return errors.New("file does not exists")
	}

	if !strings.HasSuffix(f.ConfigPath, ".yml") && !strings.HasSuffix(f.ConfigPath, ".yaml") {
		return errors.New("unsupported file format")
	}

	return nil
}
