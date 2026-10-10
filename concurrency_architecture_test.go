package mpt_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestProductionPackagesOwnNoGoroutines(t *testing.T) {
	t.Parallel()

	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate module source")
	}
	moduleDirectory := filepath.Dir(sourceFile)
	positions, err := productionGoroutines(moduleDirectory)
	if err != nil {
		t.Fatalf("inspect production sources: %v", err)
	}
	for _, position := range positions {
		t.Errorf("production goroutine at %s", position)
	}
}

func TestProductionGoroutineScopeExcludesExternalTooling(t *testing.T) {
	t.Parallel()
	directory := t.TempDir()
	tooling := filepath.Join(directory, ".golib-tooling")
	if err := os.Mkdir(tooling, 0o700); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{
		filepath.Join(directory, "owned.go"),
		filepath.Join(tooling, "external.go"),
	} {
		if err := os.WriteFile(path, []byte("package fixture\nfunc start() { go func() {}() }\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	positions, err := productionGoroutines(directory)
	if err != nil {
		t.Fatal(err)
	}
	if len(positions) != 1 || positions[0].Filename != filepath.Join(directory, "owned.go") {
		t.Fatalf("expected only the owned production goroutine, got %v", positions)
	}
}

func productionGoroutines(moduleDirectory string) ([]token.Position, error) {
	var positions []token.Position
	fileSet := token.NewFileSet()
	err := filepath.WalkDir(
		moduleDirectory,
		func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() {
				switch entry.Name() {
				case ".ai", "_interop", "benchmarks", "node_modules", "testdata":
					if path != moduleDirectory {
						return filepath.SkipDir
					}
				case ".golib-tooling":
					// CI checks out an independent tooling module here, not package production code.
					if path == filepath.Join(moduleDirectory, ".golib-tooling") {
						return filepath.SkipDir
					}
				}
				return nil
			}
			if !strings.HasSuffix(entry.Name(), ".go") ||
				strings.HasSuffix(entry.Name(), "_test.go") {
				return nil
			}
			parsed, parseErr := parser.ParseFile(
				fileSet,
				path,
				nil,
				parser.SkipObjectResolution,
			)
			if parseErr != nil {
				return parseErr
			}
			ast.Inspect(parsed, func(node ast.Node) bool {
				statement, isGoStatement := node.(*ast.GoStmt)
				if isGoStatement {
					positions = append(positions, fileSet.Position(statement.Go))
				}
				return true
			})
			return nil
		},
	)
	return positions, err
}
