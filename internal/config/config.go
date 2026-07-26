package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"

	"github.com/yuan-shuo/copy-paste-ai/assets"
)

type Config struct {
	Default     DefaultConfig       `toml:"default"`
	FileAliases map[string][]string `toml:"file_aliases"`
	Prompt      PromptConfig        `toml:"prompt"`
	Gitignore   GitignoreConfig     `toml:"gitignore"`
}

type DefaultConfig struct {
	Files []string `toml:"files"`
}

type PromptConfig struct {
	Content string `toml:"content"`
}

type GitignoreConfig struct {
	Enabled *bool `toml:"enabled"`
}

func defaultConfig() Config {
	defaultEnabled := true
	return Config{
		Default: DefaultConfig{
			Files: []string{},
		},
		FileAliases: make(map[string][]string),
		Prompt: PromptConfig{
			Content: assets.DefaultPrompt(),
		},
		Gitignore: GitignoreConfig{
			Enabled: &defaultEnabled,
		},
	}
}

func Load(rootDir string) Config {
	cpaDir := filepath.Join(rootDir, ".cpa")
	configPath := filepath.Join(cpaDir, "config.toml")

	defaultCfg := defaultConfig()

	if err := os.MkdirAll(cpaDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "警告: 创建 .cpa 目录失败: %v\n", err)
		return defaultCfg
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.WriteFile(configPath, []byte(assets.ConfigTemplate()), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "警告: 写入默认配置文件失败: %v\n", err)
		}
		fmt.Printf("已创建默认配置文件（所有配置已注释）: %s\n", configPath)
		fmt.Println("提示: 编辑该文件取消注释以启用配置")
		return defaultCfg
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "警告: 读取配置文件失败: %v，使用默认配置\n", err)
		return defaultCfg
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "警告: 解析配置文件失败: %v，使用默认配置\n", err)
		return defaultCfg
	}

	if cfg.Prompt.Content == "" {
		cfg.Prompt.Content = defaultCfg.Prompt.Content
	}
	if cfg.FileAliases == nil {
		cfg.FileAliases = make(map[string][]string)
	}
	if cfg.Default.Files == nil {
		cfg.Default.Files = []string{}
	}
	if cfg.Gitignore.Enabled == nil {
		defaultEnabled := true
		cfg.Gitignore.Enabled = &defaultEnabled
	}

	gitignoreStatus := "启用"
	if !*cfg.Gitignore.Enabled {
		gitignoreStatus = "禁用"
	}

	fmt.Printf("已加载配置: %s\n", configPath)
	fmt.Printf("Gitignore: %s\n", gitignoreStatus)
	if len(cfg.Default.Files) > 0 {
		fmt.Printf("默认文件: %v\n", cfg.Default.Files)
	}

	return cfg
}
