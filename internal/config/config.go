package config

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	configassets "github.com/yy/len/doc/config"
)

const (
	EnvVar            = "LEN_CONFIG_DIR"
	defaultDirName    = ".len"
	configFileName    = "config.toml"
	defaultTemplate   = "template.toml"
	defaultStyleDir   = "styles"
	defaultQuasiStyle = "ProceduralAlgorithm"
)

type Settings struct {
	Dir               string
	ConfigFile        string
	ProfileDir        string
	DefaultQuasiStyle string
}

type MissingConfigError struct {
	Dir string
}

func (e *MissingConfigError) Error() string {
	return fmt.Sprintf("LEN config is not initialized in %s; run `len-cli config init`", e.Dir)
}

func Resolve() (Settings, error) {
	dir, err := resolveConfigDir(os.Getenv)
	if err != nil {
		return Settings{}, err
	}
	initialized, err := IsInitializedDir(dir)
	if err != nil {
		return Settings{}, err
	}
	if !initialized {
		return Settings{}, &MissingConfigError{Dir: dir}
	}
	cfg, err := loadConfig(dir)
	if err != nil {
		return Settings{}, err
	}
	profileDir := cfg.StyleDir
	if !filepath.IsAbs(profileDir) {
		profileDir = filepath.Join(dir, profileDir)
	}
	return Settings{
		Dir:               dir,
		ConfigFile:        filepath.Join(dir, configFileName),
		ProfileDir:        profileDir,
		DefaultQuasiStyle: cfg.DefaultQuasiStyle,
	}, nil
}

func resolveConfigDir(getenv func(string) string) (string, error) {
	if dir := strings.TrimSpace(getenv(EnvVar)); dir != "" {
		return filepath.Abs(dir)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve %s: %w", EnvVar, err)
	}
	return filepath.Join(home, defaultDirName), nil
}

func ResolveDir() (string, error) {
	return resolveConfigDir(os.Getenv)
}

func IsInitializedDir(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if len(entries) == 0 {
		return false, nil
	}
	if _, err := os.Stat(filepath.Join(dir, configFileName)); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func InitDir(dir string) (bool, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return false, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	if len(entries) != 0 {
		return false, nil
	}

	// TODO: allow bootstrapping or refreshing official config assets from the repo or remote source.
	if err := fs.WalkDir(configassets.DefaultFS, ".", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == "." {
			return nil
		}
		target := path
		if path == defaultTemplate {
			target = configFileName
		}
		fullPath := filepath.Join(dir, filepath.FromSlash(target))
		if d.IsDir() {
			return os.MkdirAll(fullPath, 0o755)
		}
		data, err := fs.ReadFile(configassets.DefaultFS, path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			return err
		}
		return os.WriteFile(fullPath, data, 0o644)
	}); err != nil {
		return false, err
	}
	return true, nil
}

type fileConfig struct {
	DefaultQuasiStyle string
	StyleDir          string
}

func loadConfig(dir string) (fileConfig, error) {
	cfg := fileConfig{
		DefaultQuasiStyle: defaultQuasiStyle,
		StyleDir:          defaultStyleDir,
	}
	path := filepath.Join(dir, configFileName)
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return fileConfig{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, "\"")
		switch key {
		case "default_quasi_style":
			if value != "" {
				cfg.DefaultQuasiStyle = value
			}
		case "style_dir":
			if value != "" {
				cfg.StyleDir = value
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fileConfig{}, err
	}
	return cfg, nil
}
