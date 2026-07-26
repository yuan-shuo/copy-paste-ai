package config

import (
	"os"
	"path/filepath"
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

func TestLoad_CreatesDefaultConfig(t *testing.T) {
	dir := t.TempDir()
	cfg := Load(dir)

	if cfg.Gitignore.Enabled == nil || !*cfg.Gitignore.Enabled {
		t.Error("默认配置 Gitignore.Enabled 应为 true")
	}

	configPath := filepath.Join(dir, ".cpa", "config.toml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("应创建默认配置文件")
	}

	if cfg.Prompt.Content == "" {
		t.Error("默认配置应有 Prompt.Content")
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
