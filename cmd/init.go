package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yuan-shuo/copy-paste-ai/internal/config"
)

func NewInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "初始化 .cpa 目录及其默认配置文件",
		Long: `在项目根目录创建 .cpa/config.toml 配置文件模板。
如果配置文件已存在则不会覆盖。

使用示例:
  cpa init`,
		Run: runInit,
	}
}

func runInit(cmd *cobra.Command, args []string) {
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取当前目录失败: %v\n", err)
		os.Exit(1)
	}

	if config.ConfigExists(rootDir) {
		fmt.Println("配置文件已存在，跳过初始化")
		return
	}

	if err := config.Init(rootDir); err != nil {
		fmt.Fprintf(os.Stderr, "初始化失败: %v\n", err)
		os.Exit(1)
	}
}
