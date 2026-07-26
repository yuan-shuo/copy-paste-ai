package gitignore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewMatcherFromFile_NoFile(t *testing.T) {
	dir := t.TempDir()
	m, err := NewMatcherFromFile(dir)
	if err != nil {
		t.Fatalf("无 .gitignore 文件不应报错: %v", err)
	}
	if m == nil {
		t.Fatal("返回值不应为 nil")
	}
	if m.Match("any/file.go", false) {
		t.Error("无 .gitignore 时不应忽略任何文件")
	}
}

func TestNewMatcherFromFile_WithPatterns(t *testing.T) {
	dir := t.TempDir()
	gitignorePath := filepath.Join(dir, ".gitignore")
	content := "*.log\n.vscode/\nnode_modules/\n"
	if err := os.WriteFile(gitignorePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	m, err := NewMatcherFromFile(dir)
	if err != nil {
		t.Fatalf("加载 .gitignore 失败: %v", err)
	}

	if !m.Match("app.log", false) {
		t.Error("应忽略 .log 文件")
	}

	if !m.Match(".vscode", true) {
		t.Error("应忽略 .vscode 目录")
	}

	if !m.Match("node_modules/pkg/index.js", false) {
		t.Error("应忽略 node_modules 下的文件")
	}

	if m.Match("main.go", false) {
		t.Error("不应忽略 main.go")
	}

	if m.Match("internal/app/test.go", false) {
		t.Error("不应忽略子目录中的 .go 文件")
	}
}

func TestNewMatcherFromFile_WithCommentsAndEmptyLines(t *testing.T) {
	dir := t.TempDir()
	gitignorePath := filepath.Join(dir, ".gitignore")
	content := "# 注释行\n\n*.tmp\n# 另一个注释\n"
	if err := os.WriteFile(gitignorePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	m, err := NewMatcherFromFile(dir)
	if err != nil {
		t.Fatalf("加载失败: %v", err)
	}

	if !m.Match("file.tmp", false) {
		t.Error("应忽略 .tmp 文件")
	}
	if m.Match("file.go", false) {
		t.Error("不应忽略 .go 文件")
	}
}

func TestNewMatcherFromFile_WithDotPrefix(t *testing.T) {
	dir := t.TempDir()
	gitignorePath := filepath.Join(dir, ".gitignore")
	content := "./dist/\n"
	if err := os.WriteFile(gitignorePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	m, err := NewMatcherFromFile(dir)
	if err != nil {
		t.Fatalf("加载失败: %v", err)
	}

	if !m.Match("dist", true) {
		t.Error("应忽略 dist 目录（带 ./ 前缀）")
	}
}

func TestMatcher_Match_NilMatcher(t *testing.T) {
	m := &Matcher{}
	if m.Match("any/path", false) {
		t.Error("nil matcher 不应匹配任何路径")
	}
}
