package cmd

import (
	"bytes"
	"flag"
	"os"
	"strings"
	"testing"
)

func TestParseFlags(t *testing.T) {
	resetFlags := func() {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}

	t.Run("default config path", func(t *testing.T) {
		resetFlags()
		os.Args = []string{"cmd"}
		configPath := ParseFlags()
		expectedPath := "./config.json"
		if configPath != expectedPath {
			t.Errorf("Expected config path to be '%s', got '%s'", expectedPath, configPath)
		}
	})

	t.Run("custom config path", func(t *testing.T) {
		resetFlags()
		customPath := "/custom/config.json"
		os.Args = []string{"cmd", "-config", customPath}
		configPath := ParseFlags()
		if configPath != customPath {
			t.Errorf("Expected config path to be '%s', got '%s'", customPath, configPath)
		}
	})

	t.Run("help flag", func(t *testing.T) {
		resetFlags()
		os.Args = []string{"cmd", "-help"}
		var buf bytes.Buffer
		flag.CommandLine.SetOutput(&buf)
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected os.Exit to be called")
			}
		}()
		ParseFlags()
		expectedHelp := "Usage of cmd:\n  -config string\n    \tpath to config file (default \"./config.json\")\n  -help\n    \tprint help and exit\n  -version\n    \tprint version and exit\n"
		if !strings.Contains(buf.String(), expectedHelp) {
			t.Errorf("Expected help output to contain:\n%s\nGot:\n%s", expectedHelp, buf.String())
		}
	})

	t.Run("version flag", func(t *testing.T) {
		resetFlags()
		os.Args = []string{"cmd", "-version"}
		var buf bytes.Buffer
		flag.CommandLine.SetOutput(&buf)
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected os.Exit to be called")
			}
		}()
		ParseFlags()
		expectedVersion := "0.0.1\n"
		if buf.String() != expectedVersion {
			t.Errorf("Expected version output to be '%s', got '%s'", expectedVersion, buf.String())
		}
	})
}
