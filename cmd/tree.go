package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yuan-shuo/copy-paste-ai/internal/config"
	"github.com/yuan-shuo/copy-paste-ai/internal/gitignore"
	"github.com/yuan-shuo/copy-paste-ai/internal/tree"
)

func NewTreeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tree",
		Short: "在终端打印项目文件树（不生成任何文件）",
		Long: `在终端直接打印项目文件树，遵循以下规则:
- 自动忽略 .git 和 .cpa 目录
- 遵循 .gitignore 规则（默认启用，可在配置中关闭）
- 不生成任何文件

使用示例:
  cpa tree`,
		Run: runTree,
	}
}

func runTree(cmd *cobra.Command, args []string) {
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取当前目录失败: %v\n", err)
		os.Exit(1)
	}

	cfg := config.Load(rootDir)

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

	fmt.Print(treeStr)
}
