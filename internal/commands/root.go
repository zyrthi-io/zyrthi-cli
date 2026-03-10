package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version 版本号
	Version = "0.1.0"
	// 用于全局的配置文件路径
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "zyrthi",
	Short: "嵌入式开发工具链",
	Long: `zyrthi - 嵌入式开发工具链

提供编译、烧录、监控一站式开发体验。`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "zyrthi.yaml", "配置文件路径")
}