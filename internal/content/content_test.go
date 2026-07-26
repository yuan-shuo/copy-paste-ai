package content

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yuan-shuo/copy-paste-ai/internal/config"
)

func setupTestFiles(t *testing.T) (string, []string) {
	t.Helper()
	dir := t.TempDir()

	files := []string{"main.go", "go.mod", "README.md"}
	contents := map[string]string{
		"main.go":   "package main\n\nfunc main() {}\n",
		"go.mod":    "module test\n\ngo 1.21\n",
		"README.md": "# Test Project\n",
	}
	for _, name := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(contents[name]), 0644); err != nil {
			t.Fatal(err)
		}
	}

	return dir, files
}

func TestBuild_WithFiles(t *testing.T) {
	dir, files := setupTestFiles(t)

	cfg := config.Config{
		Prompt: config.PromptConfig{
			Content: "# test prompt",
		},
	}

	result, err := Build("root/tree", dir, files, cfg)
	if err != nil {
		t.Fatalf("Build 失败: %v", err)
	}

	if !strings.Contains(result, "root/tree") {
		t.Error("结果应包含文件树")
	}
	if !strings.Contains(result, "main.go") {
		t.Error("结果应包含 main.go 文件名")
	}
	if !strings.Contains(result, "go.mod") {
		t.Error("结果应包含 go.mod 文件名")
	}
	if !strings.Contains(result, "README.md") {
		t.Error("结果应包含 README.md 文件名")
	}
	if !strings.Contains(result, "package") {
		t.Error("结果应包含文件内容")
	}
	if !strings.Contains(result, "mod") {
		t.Error("结果应包含 go.mod 的语言标记")
	}
	if !strings.Contains(result, "# test prompt") {
		t.Error("结果应包含提示词")
	}
	if !strings.Contains(result, "```") {
		t.Error("结果应包含代码块标记")
	}
}

func TestBuild_WithoutFiles(t *testing.T) {
	dir := t.TempDir()

	cfg := config.Config{
		Prompt: config.PromptConfig{
			Content: "# prompt without files",
		},
	}

	result, err := Build("tree content", dir, nil, cfg)
	if err != nil {
		t.Fatalf("Build 失败: %v", err)
	}

	if !strings.Contains(result, "tree content") {
		t.Error("结果应包含文件树")
	}
	if !strings.Contains(result, "# prompt without files") {
		t.Error("结果应包含提示词")
	}
}

func TestBuild_FileNotFound(t *testing.T) {
	dir := t.TempDir()

	cfg := config.Config{
		Prompt: config.PromptConfig{
			Content: "test",
		},
	}

	result, err := Build("tree", dir, []string{"nonexistent.go"}, cfg)
	if err != nil {
		t.Fatalf("Build 不应报错（文件不存在应跳过）: %v", err)
	}

	if !strings.Contains(result, "tree") {
		t.Error("结果应包含文件树")
	}
}

func TestBuild_LanguageDetection(t *testing.T) {
	dir := t.TempDir()

	testCases := []struct {
		filename string
		lang     string
	}{
		{"main.go", "go"},
		{"test.py", "py"},
		{"app.js", "js"},
		{"data.json", "json"},
		{"noext", "text"},
	}

	for _, tc := range testCases {
		path := filepath.Join(dir, tc.filename)
		if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}

		cfg := config.Config{
			Prompt: config.PromptConfig{Content: "test"},
		}
		result, err := Build("tree", dir, []string{tc.filename}, cfg)
		if err != nil {
			t.Fatal(err)
		}

		expectedFence := "```" + tc.lang
		if !strings.Contains(result, expectedFence) {
			t.Errorf("文件 %s 应使用语言标记 %s", tc.filename, tc.lang)
		}
	}
}

func TestFileData_Structure(t *testing.T) {
	fd := FileData{
		Path:    "test.go",
		Content: "package test",
		Lang:    "go",
	}

	if fd.Path != "test.go" {
		t.Error("Path 字段不正确")
	}
	if fd.Content != "package test" {
		t.Error("Content 字段不正确")
	}
	if fd.Lang != "go" {
		t.Error("Lang 字段不正确")
	}
}

func TestRenderData_Structure(t *testing.T) {
	rd := RenderData{
		Tree:   "tree content",
		Files:  []FileData{{Path: "a.go", Content: "code", Lang: "go"}},
		Prompt: "prompt text",
	}

	if rd.Tree != "tree content" {
		t.Error("Tree 字段不正确")
	}
	if len(rd.Files) != 1 {
		t.Error("Files 字段长度不正确")
	}
	if rd.Prompt != "prompt text" {
		t.Error("Prompt 字段不正确")
	}
}
