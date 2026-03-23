package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yy/len/internal/config"
)

func TestResolveReportsMissingConfigForAbsentDirectory(t *testing.T) {
	t.Setenv(config.EnvVar, filepath.Join(t.TempDir(), "missing-len-config"))
	_, err := config.Resolve()
	if err == nil {
		t.Fatal("Resolve returned nil error for missing config")
	}
	if _, ok := err.(*config.MissingConfigError); !ok {
		t.Fatalf("Resolve error = %T, want *config.MissingConfigError", err)
	}
}

func TestResolveReportsMissingConfigForUninitializedDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "len-config")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	t.Setenv(config.EnvVar, dir)
	_, err := config.Resolve()
	if err == nil {
		t.Fatal("Resolve returned nil error for uninitialized config")
	}
	if _, ok := err.(*config.MissingConfigError); !ok {
		t.Fatalf("Resolve error = %T, want *config.MissingConfigError", err)
	}
}

func TestInitDirBootstrapsEmptyConfigDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "len-config")
	bootstrapped, err := config.InitDir(dir)
	if err != nil {
		t.Fatalf("InitDir returned error: %v", err)
	}
	if !bootstrapped {
		t.Fatal("InitDir reported no bootstrap for an empty directory")
	}
	assertExists(t, filepath.Join(dir, "config.toml"))
	assertExists(t, filepath.Join(dir, "styles", "procedural-algorithm.quasi-style.yaml"))
}

func TestInitDirDoesNotOverwriteExistingFiles(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "len-config")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte("default_quasi_style = \"Custom\"\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	bootstrapped, err := config.InitDir(dir)
	if err != nil {
		t.Fatalf("InitDir returned error: %v", err)
	}
	if bootstrapped {
		t.Fatal("InitDir reported bootstrap for a non-empty directory")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(data) != "default_quasi_style = \"Custom\"\n" {
		t.Fatalf("config.toml was overwritten: %q", string(data))
	}
}

func TestResolveReadsInitializedConfig(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "len-config")
	bootstrapped, err := config.InitDir(dir)
	if err != nil {
		t.Fatalf("InitDir returned error: %v", err)
	}
	if !bootstrapped {
		t.Fatal("InitDir reported no bootstrap for an empty directory")
	}
	t.Setenv(config.EnvVar, dir)
	settings, err := config.Resolve()
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if settings.Dir != dir {
		t.Fatalf("settings.Dir = %q, want %q", settings.Dir, dir)
	}
	if settings.DefaultQuasiStyle != "ProceduralAlgorithm" {
		t.Fatalf("settings.DefaultQuasiStyle = %q", settings.DefaultQuasiStyle)
	}
}

func assertExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %s to exist: %v", path, err)
	}
}
