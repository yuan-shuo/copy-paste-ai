package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/pelletier/go-toml/v2"
)

// Config 配置结构
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

// 全局匹配器
var matcher gitignore.Matcher

func main() {
	// 解析命令行参数
	var filesFlag string
	flag.StringVar(&filesFlag, "files", "", "指定要包含的文件列表，用逗号分隔，例如: internal/test.go,main.go")
	flag.StringVar(&filesFlag, "f", "", "指定要包含的文件列表，用逗号分隔（简写）")
	flag.Parse()

	// 获取当前目录
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取当前目录失败: %v\n", err)
		os.Exit(1)
	}

	// 加载配置
	config := loadConfig(rootDir)

	// 解析文件列表
	var fileList []string
	if filesFlag != "" {
		parts := strings.Split(filesFlag, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				// 检查是否是别名
				if expanded, ok := config.FileAliases[p]; ok {
					fileList = append(fileList, expanded...)
				} else {
					fileList = append(fileList, p)
				}
			}
		}
	}

	// 添加默认文件（去重）
	existing := make(map[string]bool)
	for _, f := range fileList {
		existing[f] = true
	}
	for _, f := range config.Default.Files {
		if !existing[f] {
			fileList = append(fileList, f)
		}
	}

	// 验证文件是否存在
	if len(fileList) > 0 {
		for _, relPath := range fileList {
			fullPath := filepath.Join(rootDir, relPath)
			if _, err := os.Stat(fullPath); err != nil {
				fmt.Fprintf(os.Stderr, "错误: 文件不存在: %s\n", relPath)
				os.Exit(1)
			}
		}
	}

	// 判断是否启用 gitignore（默认启用，只有显式设置为 false 才关闭）
	gitignoreEnabled := true
	if config.Gitignore.Enabled != nil {
		gitignoreEnabled = *config.Gitignore.Enabled
	}

	// 加载 .gitignore 规则
	if gitignoreEnabled {
		if err := loadGitignore(rootDir); err != nil {
			fmt.Fprintf(os.Stderr, "警告: %v\n", err)
		}
	}

	// 生成文件树
	tree := generateTree(rootDir, gitignoreEnabled)

	// 创建 .cpa/prompt 目录
	cpaDir := filepath.Join(rootDir, ".cpa")
	promptDir := filepath.Join(cpaDir, "prompt")
	if err := os.MkdirAll(promptDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "创建目录失败: %v\n", err)
		os.Exit(1)
	}

	// 创建 .gitignore 文件
	gitignorePath := filepath.Join(cpaDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte("*\n"), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "创建 .gitignore 失败: %v\n", err)
		os.Exit(1)
	}

	// 生成文件名: 时间戳.md
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s.md", timestamp)
	filePath := filepath.Join(promptDir, filename)

	// 构建内容
	content := buildContent(tree, rootDir, fileList, config)

	// 写入文件
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "写入文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n已生成: %s\n", filePath)
	if len(fileList) > 0 {
		fmt.Printf("包含文件: %v\n", fileList)
	}
}

// loadGitignore 加载 .gitignore 文件
func loadGitignore(rootDir string) error {
	gitignorePath := filepath.Join(rootDir, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(gitignorePath)
	if err != nil {
		return fmt.Errorf("打开 .gitignore 失败: %w", err)
	}
	defer file.Close()

	var patterns []gitignore.Pattern
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// 去掉开头的 ./ 前缀（gitignore 标准格式不支持）
		line = strings.TrimPrefix(line, "./")
		pattern := gitignore.ParsePattern(line, nil)
		patterns = append(patterns, pattern)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取 .gitignore 失败: %w", err)
	}

	matcher = gitignore.NewMatcher(patterns)
	return nil
}

// shouldIgnore 检查路径是否应该被忽略
func shouldIgnore(relPath string, isDir bool) bool {
	// 硬编码忽略 .git 目录
	if strings.HasPrefix(relPath, ".git") || strings.HasPrefix(relPath, ".git/") {
		return true
	}

	if matcher == nil {
		return false
	}

	// 将路径转换为正斜杠分隔
	parts := strings.Split(filepath.ToSlash(relPath), "/")
	return matcher.Match(parts, isDir)
}

// generateTree 生成文件树
func generateTree(rootDir string, gitignoreEnabled bool) string {
	var builder strings.Builder
	builder.WriteString(rootDir)
	builder.WriteString("\n")

	filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		// 硬编码跳过 .cpa 目录
		if d.IsDir() && d.Name() == ".cpa" {
			return fs.SkipDir
		}

		// 硬编码跳过 .git 目录
		if d.IsDir() && d.Name() == ".git" {
			return fs.SkipDir
		}

		relPath, _ := filepath.Rel(rootDir, path)
		if relPath == "." {
			return nil
		}

		// 应用 gitignore 规则
		if gitignoreEnabled {
			if shouldIgnore(relPath, d.IsDir()) {
				if d.IsDir() {
					return fs.SkipDir
				}
				return nil
			}
		}

		depth := strings.Count(relPath, string(filepath.Separator))
		prefix := strings.Repeat("    ", depth)

		parent := filepath.Dir(relPath)
		entries, err := os.ReadDir(filepath.Join(rootDir, parent))
		if err != nil {
			return nil
		}

		var names []string
		for _, e := range entries {
			if e.Name() == ".cpa" || e.Name() == ".git" {
				continue
			}
			eRelPath, _ := filepath.Rel(rootDir, filepath.Join(parent, e.Name()))
			if gitignoreEnabled {
				if shouldIgnore(eRelPath, e.IsDir()) {
					continue
				}
			}
			names = append(names, e.Name())
		}

		isLast := false
		for i, name := range names {
			if name == d.Name() {
				if i == len(names)-1 {
					isLast = true
				}
				break
			}
		}

		connector := "├── "
		if isLast {
			connector = "└── "
		}

		if depth == 0 {
			rootEntries, err := os.ReadDir(rootDir)
			if err == nil {
				var rootNames []string
				for _, e := range rootEntries {
					if e.Name() == ".cpa" || e.Name() == ".git" {
						continue
					}
					eRelPath, _ := filepath.Rel(rootDir, filepath.Join(rootDir, e.Name()))
					if gitignoreEnabled {
						if shouldIgnore(eRelPath, e.IsDir()) {
							continue
						}
					}
					rootNames = append(rootNames, e.Name())
				}
				for i, name := range rootNames {
					if name == d.Name() && i == len(rootNames)-1 {
						connector = "└── "
					}
				}
			}
		}

		name := d.Name()
		if d.IsDir() {
			name = name + "/"
		}

		builder.WriteString(prefix)
		builder.WriteString(connector)
		builder.WriteString(name)
		builder.WriteString("\n")

		return nil
	})

	return builder.String()
}

// buildContent 构建文件内容
func buildContent(tree, rootDir string, fileList []string, config Config) string {
	var builder strings.Builder

	builder.WriteString("# 项目文件树\n\n")
	builder.WriteString("```\n")
	builder.WriteString(tree)
	builder.WriteString("```\n")

	if len(fileList) > 0 {
		builder.WriteString("\n# 文件代码\n\n")
		for _, relPath := range fileList {
			fullPath := filepath.Join(rootDir, relPath)
			content, err := os.ReadFile(fullPath)
			if err != nil {
				continue
			}
			ext := filepath.Ext(relPath)
			lang := strings.TrimPrefix(ext, ".")
			if lang == "" {
				lang = "text"
			}
			builder.WriteString(fmt.Sprintf("## %s\n\n", relPath))
			builder.WriteString(fmt.Sprintf("```%s\n", lang))
			builder.WriteString(string(content))
			if len(content) > 0 && content[len(content)-1] != '\n' {
				builder.WriteString("\n")
			}
			builder.WriteString("```\n\n")
		}
	}

	builder.WriteString(config.Prompt.Content)

	return builder.String()
}

// loadConfig 加载配置文件
func loadConfig(rootDir string) Config {
	cpaDir := filepath.Join(rootDir, ".cpa")
	configPath := filepath.Join(cpaDir, "config.toml")

	// 默认启用 gitignore
	defaultEnabled := true

	defaultConfig := Config{
		Default: DefaultConfig{
			Files: []string{},
		},
		FileAliases: make(map[string][]string),
		Prompt: PromptConfig{
			Content: `# 提示

请根据以上项目文件树和代码，完成用户的需求。

1. 对于需要改动的文件，给出完整的代码，禁止删除用户原有的注释
2. 最终总结所有改动
3. 检查代码是否符合语言规范
`,
		},
		Gitignore: GitignoreConfig{
			Enabled: &defaultEnabled,
		},
	}

	// 确保 .cpa 目录存在
	if err := os.MkdirAll(cpaDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "警告: 创建 .cpa 目录失败: %v\n", err)
		return defaultConfig
	}

	// 如果配置文件不存在，创建默认配置（所有内容注释状态）
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		commentedConfig := `# ============================================================
# Copy-Paste AI 配置文件
# 取消注释即可启用对应功能
# ============================================================

# [default]
# files = ["main.go", "go.mod"]

# [file_aliases]
# "main" = ["main.go"]
# "test" = ["internal/test.go"]
# "core" = ["pkg/core/core.go", "pkg/core/utils.go"]
# "all" = ["main.go", "internal/test.go", "pkg/core/core.go"]

# [prompt]
# content = """
# # 提示
# 
# 请根据以上项目文件树和代码，完成用户的需求。
# 
# 1. 对于需要改动的文件，给出完整的代码，禁止删除用户原有的注释
# 2. 最终总结所有改动
# 3. 检查代码是否符合语言规范
# 4. 给出可执行的终端命令
# """

# [gitignore]
# enabled = true
`
		if err := os.WriteFile(configPath, []byte(commentedConfig), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "警告: 写入默认配置文件失败: %v\n", err)
		}
		fmt.Printf("已创建默认配置文件（所有配置已注释）: %s\n", configPath)
		fmt.Println("提示: 编辑该文件取消注释以启用配置")
		return defaultConfig
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "警告: 读取配置文件失败: %v，使用默认配置\n", err)
		return defaultConfig
	}

	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		fmt.Fprintf(os.Stderr, "警告: 解析配置文件失败: %v，使用默认配置\n", err)
		return defaultConfig
	}

	// 如果配置中 Prompt.Content 为空，使用默认
	if config.Prompt.Content == "" {
		config.Prompt.Content = defaultConfig.Prompt.Content
	}
	if config.FileAliases == nil {
		config.FileAliases = make(map[string][]string)
	}
	if config.Default.Files == nil {
		config.Default.Files = []string{}
	}
	// 如果 gitignore 未配置，默认启用
	if config.Gitignore.Enabled == nil {
		defaultEnabled := true
		config.Gitignore.Enabled = &defaultEnabled
	}

	gitignoreStatus := "启用"
	if !*config.Gitignore.Enabled {
		gitignoreStatus = "禁用"
	}

	fmt.Printf("已加载配置: %s\n", configPath)
	fmt.Printf("Gitignore: %s\n", gitignoreStatus)
	if len(config.Default.Files) > 0 {
		fmt.Printf("默认文件: %v\n", config.Default.Files)
	}

	return config
}
