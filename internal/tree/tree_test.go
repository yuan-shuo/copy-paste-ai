package tree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yuan-shuo/copy-paste-ai/internal/gitignore"
)

func createTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	files := map[string]string{
		"main.go":            "package main\n",
		"go.mod":             "module test\n",
		"internal/app/a.go":  "package app\n",
		"internal/app/b.go":  "package app\n",
		"internal/util/u.go": "package util\n",
		"pkg/core/c.go":      "package core\n",
	}

	for path, content := range files {
		fullPath := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	return dir
}

func TestGenerate_BasicTree(t *testing.T) {
	dir := createTestDir(t)

	result := Generate(dir, false, nil)

	if !strings.Contains(result, dir) {
		t.Error("文件树应包含根目录路径")
	}
	if !strings.Contains(result, "main.go") {
		t.Error("文件树应包含 main.go")
	}
	if !strings.Contains(result, "go.mod") {
		t.Error("文件树应包含 go.mod")
	}
	if !strings.Contains(result, "internal/") {
		t.Error("文件树应包含 internal/ 目录")
	}
	if !strings.Contains(result, "pkg/") {
		t.Error("文件树应包含 pkg/ 目录")
	}
}

func TestGenerate_SkipsGitAndCpaDirs(t *testing.T) {
	dir := createTestDir(t)

	gitDir := filepath.Join(dir, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gitDir, "config"), []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	cpaDir := filepath.Join(dir, ".cpa")
	if err := os.MkdirAll(cpaDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cpaDir, "config.toml"), []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	result := Generate(dir, false, nil)

	if strings.Contains(result, ".git") {
		t.Error("文件树不应包含 .git 目录")
	}
	if strings.Contains(result, ".cpa") {
		t.Error("文件树不应包含 .cpa 目录")
	}
}

func TestGenerate_WithGitignore(t *testing.T) {
	dir := createTestDir(t)

	gitignorePath := filepath.Join(dir, ".gitignore")
	content := "*.log\ninternal/util/\n"
	if err := os.WriteFile(gitignorePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(dir, "app.log"), []byte("log"), 0644); err != nil {
		t.Fatal(err)
	}

	m, err := gitignore.NewMatcherFromFile(dir)
	if err != nil {
		t.Fatal(err)
	}

	result := Generate(dir, true, m)

	if strings.Contains(result, ".log") {
		t.Error("文件树不应包含 .log 文件")
	}
	if strings.Contains(result, "u.go") {
		t.Error("文件树不应包含被忽略目录中的文件")
	}
	if !strings.Contains(result, "a.go") {
		t.Error("文件树应包含未被忽略的文件")
	}
	if !strings.Contains(result, "b.go") {
		t.Error("文件树应包含未被忽略的文件")
	}
}

func TestGenerate_DirectorySuffix(t *testing.T) {
	dir := createTestDir(t)
	result := Generate(dir, false, nil)

	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	for _, line := range lines[1:] {
		if strings.Contains(line, "internal/") && !strings.HasSuffix(strings.TrimSpace(line), "/") {
			t.Errorf("目录行应以 / 结尾: %s", line)
		}
	}
}

func TestShouldIgnore(t *testing.T) {
	if !shouldIgnore(".git", true, nil) {
		t.Error(".git 目录应被忽略")
	}
	if !shouldIgnore(".git/config", false, nil) {
		t.Error(".git 下的文件应被忽略")
	}
	if !shouldIgnore(".gitignore", false, nil) {
		t.Error(".gitignore 应以 .git 开头被忽略")
	}
	if shouldIgnore("main.go", false, nil) {
		t.Error("main.go 不应被忽略")
	}

	m, _ := gitignore.NewMatcherFromFile(t.TempDir())
	if shouldIgnore("any/file.go", false, m) {
		t.Error("无规则时不应忽略文件")
	}
}
