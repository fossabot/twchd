package main

import (
	"os"
	"testing"
)

func TestVerifyPath_Empty(t *testing.T) {
	var actual = VerifyPath("")
	if actual == nil {
		t.Fatal(actual)
	}
}

func TestVerifyPath_Space(t *testing.T) {
	var actual = VerifyPath("\t")
	if actual == nil {
		t.Fatal(actual)
	}
}
func TestVerifyPath_Exists(t *testing.T) {
	var actual = VerifyPath("config/vanya83.yml")
	if actual != nil {
		t.Fatal("file exists")
	}
}

func TestVerifyPath_UnExists(t *testing.T) {
	var actual = VerifyPath("filename.yml")
	if actual == nil {
		t.Fatal(actual)
	}
}

func TestVerifyPath_UnexpectedFormat(t *testing.T) {
	var actual = VerifyPath("assets/pipeline.json")
	if actual == nil {
		t.Fatal(actual)
	}
}
func TestGetAccountName_String(t *testing.T) {
	var config = &BotConfig{
		AccountName: "vanya123",
	}
	accountName, _ := config.GetAccountName()
	var expected = "vanya123"
	if accountName != expected {
		t.Fatal(expected, "expected, gotten", accountName)
	}
}

func TestGetAccountName_Env(t *testing.T) {
	var config = &BotConfig{
		AccountName: "${TWITCH_ACCOUNT}",
	}
	os.Setenv("TWITCH_ACCOUNT", "vanya123")
	accountName, _ := config.GetAccountName()
	var expected = "vanya123"
	if accountName != expected {
		t.Fatal(expected, "expected, gotten", accountName)
	}
}
func TestGetAccountName_EmptyEnv(t *testing.T) {
	var config = &BotConfig{
		AccountName: "${TWITCH_ACCOUNT}",
	}
	os.Setenv("TWITCH_ACCOUNT", "")
	_, err := config.GetAccountName()
	if err == nil {
		t.Fatal("empty env var should not be read")
	}
}

func TestGetAccountName_EmptyStr(t *testing.T) {
	var config = &BotConfig{
		AccountName: "",
	}
	_, err := config.GetAccountName()
	if err == nil {
		t.Fatal("empty field should not be read")
	}
}

func TestGetTokene_String(t *testing.T) {
	var config = &BotConfig{
		AccountToken: "oauth:123",
	}
	accountName, _ := config.GetToken()
	var expected = "oauth:123"
	if accountName != expected {
		t.Fatal(expected, "expected, gotten", accountName)
	}
}

func TestGetToken_EmptyStr(t *testing.T) {
	var config = &BotConfig{
		AccountToken: "",
	}
	_, err := config.GetToken()
	if err == nil {
		t.Fatal("empty field should not be read")
	}
}

func TestGetToken_Env(t *testing.T) {
	var config = &BotConfig{
		AccountToken: "${TWITCH_OAUTH}",
	}
	os.Setenv("TWITCH_OAUTH", "oauth:123")
	accountName, _ := config.GetToken()
	var expected = "oauth:123"
	if accountName != expected {
		t.Fatal(expected, "expected, gotten", accountName)
	}
}
func TestGetToken_EmptyEnv(t *testing.T) {
	var config = &BotConfig{
		AccountToken: "${TWITCH_OAUTH}",
	}
	os.Setenv("TWITCH_OAUTH", "")
	_, err := config.GetToken()
	if err == nil {
		t.Fatal("empty env var should not be read")
	}
}
