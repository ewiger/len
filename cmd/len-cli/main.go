package main

import (
	"fmt"
	"os"
	"path/filepath"

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
	default:
		printUsage()
		os.Exit(2)
	}
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
	program := loader.Loader{Root: cwd}.LoadPaths(paths)
	diagnostics := validator.Validator{
		ProfileDir:        filepath.Join(cwd, "doc", "proposals", "accepted", "lip-0001-cli-parser-n-validation"),
		DefaultQuasiStyle: "ProceduralAlgorithm",
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
}
