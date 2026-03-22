package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yy/len/internal/loader"
)

func TestLoaderLoadsTransitiveImports(t *testing.T) {
	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "lang", "l1", "core", "math", "set.l1"), "type Seq\n")
	mustWriteFile(t, filepath.Join(root, "entry.l1"), "import core.math.set\ntype Local\n")

	program := loader.Loader{Root: root}.LoadPaths([]string{filepath.Join(root, "entry.l1")})
	if len(program.Diagnostics) != 0 {
		t.Fatalf("LoadPaths diagnostics = %#v, want none", program.Diagnostics)
	}
	if len(program.Units) != 2 {
		t.Fatalf("units = %d, want 2", len(program.Units))
	}
}

func TestLoaderReportsMissingImport(t *testing.T) {
	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "entry.l1"), "import core.missing.module\n")

	program := loader.Loader{Root: root}.LoadPaths([]string{filepath.Join(root, "entry.l1")})
	if len(program.Diagnostics) == 0 {
		t.Fatal("expected missing import diagnostic")
	}
	if got := program.Diagnostics[0].Code; got != "loader.import.missing" {
		t.Fatalf("diagnostic code = %q, want loader.import.missing", got)
	}
}

func mustWriteFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}
