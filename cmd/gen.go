package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/yuan-shuo/copy-paste-ai/internal/config"
	"github.com/yuan-shuo/copy-paste-ai/internal/content"
	"github.com/yuan-shuo/copy-paste-ai/internal/gitignore"
	"github.com/yuan-shuo/copy-paste-ai/internal/tree"
)

func NewGenCmd() *cobra.Command {
	genCmd := &cobra.Command{
		Use:   "gen",
		Short: "生成 .md 文件，包含项目文件树、指定文件代码和提示词",
		Long: `生成 .md 文件到 .cpa/prompt/ 目录下，文件内容包括:
1. 项目文件树（遵循 .gitignore 规则）
2. 指定文件的完整代码
3. 配置中的提示词

使用示例:
  cpa gen -f main.go,internal/test.go
  cpa gen --files main.go,test.go`,
		Run: runGen,
	}
	genCmd.Flags().StringP("files", "f", "", "指定要包含的文件列表，用逗号分隔，例如: internal/test.go,main.go")
	return genCmd
}

func runGen(cmd *cobra.Command, args []string) {
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取当前目录失败: %v\n", err)
		os.Exit(1)
	}

	cfg := config.Load(rootDir)

	filesFlag, _ := cmd.Flags().GetString("files")

	var fileList []string
	if filesFlag != "" {
		parts := strings.Split(filesFlag, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				if expanded, ok := cfg.FileAliases[p]; ok {
					fileList = append(fileList, expanded...)
				} else {
					fileList = append(fileList, p)
				}
			}
		}
	}

	existing := make(map[string]bool)
	for _, f := range fileList {
		existing[f] = true
	}
	for _, f := range cfg.Default.Files {
		if !existing[f] {
			fileList = append(fileList, f)
		}
	}

	if len(fileList) > 0 {
		for _, relPath := range fileList {
			fullPath := filepath.Join(rootDir, relPath)
			if _, err := os.Stat(fullPath); err != nil {
				fmt.Fprintf(os.Stderr, "错误: 文件不存在: %s\n", relPath)
				os.Exit(1)
			}
		}
	}

	gitignoreEnabled := true
	if cfg.Gitignore.Enabled != nil {
		gitignoreEnabled = *cfg.Gitignore.Enabled
	}

	var matcher *gitignore.Matcher
	if gitignoreEnabled {
		var err error
		matcher, err = gitignore.NewMatcherFromFile(rootDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "警告: %v\n", err)
		}
	}

	treeStr := tree.Generate(rootDir, gitignoreEnabled, matcher)

	cpaDir := filepath.Join(rootDir, ".cpa")
	promptDir := filepath.Join(cpaDir, "prompt")
	if err := os.MkdirAll(promptDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "创建目录失败: %v\n", err)
		os.Exit(1)
	}

	gitignorePath := filepath.Join(cpaDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte("*\n"), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "创建 .gitignore 失败: %v\n", err)
		os.Exit(1)
	}

	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s.md", timestamp)
	filePath := filepath.Join(promptDir, filename)

	mdContent, err := content.Build(treeStr, rootDir, fileList, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "构建内容失败: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(filePath, []byte(mdContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "写入文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n已生成: %s\n", filePath)
	if len(fileList) > 0 {
		fmt.Printf("包含文件: %v\n", fileList)
	}
}
