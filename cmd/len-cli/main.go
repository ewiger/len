package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yy/len/internal/config"
	"github.com/yy/len/internal/diag"
	"github.com/yy/len/internal/loader"
	"github.com/yy/len/internal/validator"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "validate":
		if err := runValidate(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	case "config":
		if err := runConfig(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	default:
		printUsage()
		os.Exit(2)
	}
}

func runConfig(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("config requires a subcommand")
	}
	switch args[0] {
	case "init":
		return runConfigInit(args[1:])
	default:
		return fmt.Errorf("unknown config subcommand %q", args[0])
	}
}

func runConfigInit(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("config init does not accept positional arguments")
	}
	dir, err := config.ResolveDir()
	if err != nil {
		return err
	}
	bootstrapped, err := config.InitDir(dir)
	if err != nil {
		return err
	}
	if bootstrapped {
		fmt.Fprintf(os.Stdout, "initialized LEN config in %s\n", dir)
		return nil
	}
	fmt.Fprintf(os.Stdout, "LEN config already present in %s\n", dir)
	return nil
}

func runValidate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("validate requires at least one .l1 path")
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	paths, err := expandPaths(args)
	if err != nil {
		return err
	}
	settings, err := config.Resolve()
	if err != nil {
		return err
	}
	program := loader.Loader{Root: cwd}.LoadPaths(paths)
	diagnostics := validator.Validator{
		ProfileDir:        settings.ProfileDir,
		DefaultQuasiStyle: settings.DefaultQuasiStyle,
	}.Validate(program)
	printDiagnostics(diagnostics)
	if hasErrors(diagnostics) {
		os.Exit(1)
	}
	return nil
}

func expandPaths(args []string) ([]string, error) {
	paths := make([]string, 0)
	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			paths = append(paths, arg)
			continue
		}
		err = filepath.WalkDir(arg, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() {
				return nil
			}
			if filepath.Ext(path) == ".l1" {
				paths = append(paths, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return paths, nil
}

func printDiagnostics(items []diag.Diagnostic) {
	for _, item := range items {
		fmt.Fprintln(os.Stderr, item.String())
	}
}

func hasErrors(items []diag.Diagnostic) bool {
	for _, item := range items {
		if item.Severity == diag.SeverityError {
			return true
		}
	}
	return false
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage: len-cli validate <path-or-directory> [more paths]")
	fmt.Fprintln(os.Stderr, "       len-cli config init")
}
