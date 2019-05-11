package main

import (
	"testing"
)

func setupConfigPath(path string) error {
	var config = FlagsCLI{path, false}
	return config.VerifyPath()
}

func TestVerifyPath_Empty(t *testing.T) {
	var actual = setupConfigPath("")
	if actual == nil {
		t.Fatal(actual)
	}
}

func TestVerifyPath_Space(t *testing.T) {
	var actual = setupConfigPath("\t")
	if actual == nil {
		t.Fatal(actual)
	}
}
func TestVerifyPath_Exists(t *testing.T) {
	var actual = setupConfigPath("example/vanya83.yml")
	if actual != nil {
		t.Fatal("file exists")
	}
}

func TestVerifyPath_UnExists(t *testing.T) {
	var actual = setupConfigPath("filename.yml")
	if actual == nil {
		t.Fatal(actual)
	}
}

func TestVerifyPath_UnexpectedFormat(t *testing.T) {
	var actual = setupConfigPath("assets/pipeline.json")
	if actual == nil {
		t.Fatal(actual)
	}
}
