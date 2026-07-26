package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yuan-shuo/copy-paste-ai/cmd"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cpa",
		Short: "Copy-Paste AI - 快速生成项目上下文供 AI 使用",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Long: `Copy-Paste AI (cpa) 是一个命令行工具，用于生成项目文件树和代码上下文，
方便将项目信息粘贴给 AI 进行分析和开发。

子命令:
  init  初始化配置文件
  gen   生成 .md 文件
  tree  在终端打印文件树`,
	}

	rootCmd.AddCommand(cmd.NewInitCmd())
	rootCmd.AddCommand(cmd.NewGenCmd())
	rootCmd.AddCommand(cmd.NewTreeCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}
