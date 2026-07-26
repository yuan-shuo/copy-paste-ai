package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()

	if cfg.Default.Files == nil {
		t.Error("默认 Files 不应为 nil")
	}
	if len(cfg.Default.Files) != 0 {
		t.Errorf("默认 Files 应为空切片，实际: %v", cfg.Default.Files)
	}

	if cfg.FileAliases == nil {
		t.Error("默认 FileAliases 不应为 nil")
	}

	if cfg.Prompt.Content == "" {
		t.Error("默认 Prompt.Content 不应为空")
	}

	if cfg.Gitignore.Enabled == nil {
		t.Fatal("默认 Gitignore.Enabled 不应为 nil")
	}
	if !*cfg.Gitignore.Enabled {
		t.Error("默认 Gitignore.Enabled 应为 true")
	}
}

func TestInit_CreatesDefaultConfig(t *testing.T) {
	dir := t.TempDir()

	err := Init(dir)
	if err != nil {
		t.Fatalf("Init 失败: %v", err)
	}

	configPath := filepath.Join(dir, ".cpa", "config.toml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("应创建默认配置文件")
	}

	gitignorePath := filepath.Join(dir, ".cpa", ".gitignore")
	data, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Errorf("应创建 .cpa/.gitignore: %v", err)
	}
	if string(data) != "*\n" {
		t.Errorf(".cpa/.gitignore 内容应为 '*\\n'，实际: '%s'", string(data))
	}

	cfg := Load(dir)
	if cfg.Gitignore.Enabled == nil || !*cfg.Gitignore.Enabled {
		t.Error("默认配置 Gitignore.Enabled 应为 true")
	}
	if cfg.Prompt.Content == "" {
		t.Error("默认配置应有 Prompt.Content")
	}
}

func TestInit_AlreadyExists(t *testing.T) {
	dir := t.TempDir()

	err := Init(dir)
	if err != nil {
		t.Fatalf("Init 第一次失败: %v", err)
	}

	err = Init(dir)
	if err != nil {
		t.Fatalf("Init 第二次不应失败: %v", err)
	}

	gitignorePath := filepath.Join(dir, ".cpa", ".gitignore")
	data, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("应存在 .cpa/.gitignore: %v", err)
	}
	if string(data) != "*\n" {
		t.Errorf("多次 Init 不应覆盖 .gitignore，实际: '%s'", string(data))
	}
}

func TestConfigExists(t *testing.T) {
	dir := t.TempDir()

	if ConfigExists(dir) {
		t.Error("不存在配置文件时应返回 false")
	}

	if err := Init(dir); err != nil {
		t.Fatal(err)
	}

	if !ConfigExists(dir) {
		t.Error("存在配置文件时应返回 true")
	}
}

func TestLoad_NoConfigReturnsDefaults(t *testing.T) {
	dir := t.TempDir()

	cfg := Load(dir)

	if cfg.Gitignore.Enabled == nil || !*cfg.Gitignore.Enabled {
		t.Error("无配置时 Gitignore.Enabled 应默认为 true")
	}
	if cfg.Prompt.Content == "" {
		t.Error("无配置时应有默认 Prompt.Content")
	}
}

func TestLoad_LoadsExistingConfig(t *testing.T) {
	dir := t.TempDir()
	cpaDir := filepath.Join(dir, ".cpa")
	if err := os.MkdirAll(cpaDir, 0755); err != nil {
		t.Fatal(err)
	}

	configPath := filepath.Join(cpaDir, "config.toml")
	content := `
[default]
files = ["main.go", "test.go"]

[file_aliases]
"core" = ["core.go"]

[prompt]
content = "custom prompt"

[gitignore]
enabled = false
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := Load(dir)

	if len(cfg.Default.Files) != 2 {
		t.Errorf("期望 2 个默认文件，实际: %d", len(cfg.Default.Files))
	}
	if cfg.Default.Files[0] != "main.go" {
		t.Errorf("期望 main.go，实际: %s", cfg.Default.Files[0])
	}
	if cfg.Default.Files[1] != "test.go" {
		t.Errorf("期望 test.go，实际: %s", cfg.Default.Files[1])
	}

	if len(cfg.FileAliases) != 1 {
		t.Errorf("期望 1 个别名，实际: %d", len(cfg.FileAliases))
	}
	if aliases, ok := cfg.FileAliases["core"]; !ok || aliases[0] != "core.go" {
		t.Error("别名解析错误")
	}

	if cfg.Prompt.Content != "custom prompt" {
		t.Errorf("期望 'custom prompt'，实际: '%s'", cfg.Prompt.Content)
	}

	if cfg.Gitignore.Enabled == nil || *cfg.Gitignore.Enabled {
		t.Error("期望 Gitignore.Enabled 为 false")
	}
}

func TestLoad_EmptyConfigUsesDefaults(t *testing.T) {
	dir := t.TempDir()
	cpaDir := filepath.Join(dir, ".cpa")
	if err := os.MkdirAll(cpaDir, 0755); err != nil {
		t.Fatal(err)
	}

	configPath := filepath.Join(cpaDir, "config.toml")
	if err := os.WriteFile(configPath, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := Load(dir)

	if cfg.Prompt.Content == "" {
		t.Error("空配置应使用默认 Prompt.Content")
	}
	if cfg.FileAliases == nil {
		t.Error("空配置 FileAliases 不应为 nil")
	}
	if cfg.Default.Files == nil {
		t.Error("空配置 Default.Files 不应为 nil")
	}
	if cfg.Gitignore.Enabled == nil || !*cfg.Gitignore.Enabled {
		t.Error("空配置 Gitignore.Enabled 应默认为 true")
	}
}

func TestLoad_InvalidConfigUsesDefaults(t *testing.T) {
	dir := t.TempDir()
	cpaDir := filepath.Join(dir, ".cpa")
	if err := os.MkdirAll(cpaDir, 0755); err != nil {
		t.Fatal(err)
	}

	configPath := filepath.Join(cpaDir, "config.toml")
	if err := os.WriteFile(configPath, []byte("invalid toml content {{{}"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := Load(dir)

	if cfg.Prompt.Content == "" {
		t.Error("无效配置应回退到默认 Prompt.Content")
	}
	if cfg.Gitignore.Enabled == nil || !*cfg.Gitignore.Enabled {
		t.Error("无效配置 Gitignore.Enabled 应默认为 true")
	}
}

func TestLoad_CustomPromptWithoutContent(t *testing.T) {
	dir := t.TempDir()
	cpaDir := filepath.Join(dir, ".cpa")
	if err := os.MkdirAll(cpaDir, 0755); err != nil {
		t.Fatal(err)
	}

	configPath := filepath.Join(cpaDir, "config.toml")
	content := `
[prompt]
content = ""
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := Load(dir)
	if cfg.Prompt.Content == "" {
		t.Error("空 prompt content 应回退到默认内容")
	}
}

func TestSanitizePath(t *testing.T) {
	if got := SanitizePath(`cmd\init.go`); got != "cmd/init.go" {
		t.Errorf("应将反斜杠转为正斜杠，期望 cmd/init.go，实际: %s", got)
	}
	if got := SanitizePath("cmd/init.go"); got != "cmd/init.go" {
		t.Errorf("正斜杠应保持不变，实际: %s", got)
	}
	if got := SanitizePath(`  cmd\init.go  `); got != "cmd/init.go" {
		t.Errorf("应同时处理空格和反斜杠，实际: %s", got)
	}
}

func TestLoad_BackslashPathsInConfig(t *testing.T) {
	dir := t.TempDir()
	cpaDir := filepath.Join(dir, ".cpa")
	if err := os.MkdirAll(cpaDir, 0755); err != nil {
		t.Fatal(err)
	}

	configPath := filepath.Join(cpaDir, "config.toml")
	content := `
[default]
files = ["cmd\init.go", "internal\config\config.go"]

[file_aliases]
"core" = ["cmd\gen.go", "cmd\tree.go"]
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := Load(dir)

	for _, f := range cfg.Default.Files {
		if strings.Contains(f, `\`) {
			t.Errorf("默认文件路径不应包含反斜杠: %s", f)
		}
	}
	if len(cfg.Default.Files) != 2 {
		t.Fatalf("期望 2 个默认文件，实际: %d", len(cfg.Default.Files))
	}

	aliases := cfg.FileAliases["core"]
	for _, f := range aliases {
		if strings.Contains(f, `\`) {
			t.Errorf("别名路径不应包含反斜杠: %s", f)
		}
	}
	if len(aliases) != 2 {
		t.Fatalf("期望 2 个别名文件，实际: %d", len(aliases))
	}
}
