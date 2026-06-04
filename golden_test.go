package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	out, err := exec.Command("go", "build", "-o", "aa.test", ".").CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "building aa.test: %v\n%s", err, out)
		os.Exit(1)
	}
	code := m.Run()
	os.Remove("aa.test")
	os.Exit(code)
}

func TestGolden(t *testing.T) {
	files, err := filepath.Glob("testdata/*.aa")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no .aa files found in testdata/")
	}

	for _, file := range files {
		name := filepath.Base(file)
		base := strings.TrimSuffix(name, ".aa")
		expectedFile := filepath.Join("testdata", base+".expected")
		inputFile := filepath.Join("testdata", base+".input")

		t.Run(base, func(t *testing.T) {
			_ = readExpected(t, expectedFile)

			cmd := exec.Command("./aa.test", file)
			if fileExists(inputFile) {
				input, err := os.ReadFile(inputFile)
				if err != nil {
					t.Fatalf("reading input file %s: %v", inputFile, err)
				}
				cmd.Stdin = bytes.NewReader(input)
			}

			out, err := cmd.CombinedOutput()
			if err != nil {
				// If the program had a runtime error, ShowError writes to stdout
				// and os.Exit(1) causes CombinedOutput to return an error,
				// but stdout/stderr are still captured in out.
			}

			expected, err := os.ReadFile(expectedFile)
			if err != nil {
				t.Fatalf("reading expected file %s: %v", expectedFile, err)
			}

			got := string(out)
			want := string(expected)

			if got != want {
				t.Errorf("%s:\ngot:\n%s\nwant:\n%s", name, got, want)
				t.Logf("diff: got %d bytes, want %d bytes", len(got), len(want))
			}
		})
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func readExpected(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("missing expected file %s", path)
	}
	return string(data)
}
