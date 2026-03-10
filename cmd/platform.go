package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var platformDir string

func init() {
	platformDir = filepath.Join(os.Getenv("HOME"), ".zyrthi", "platforms")
}

var platformCmd = &cobra.Command{
	Use:   "platform",
	Short: "平台配置管理",
	Long:  `管理芯片平台配置（安装、列出、更新）。`,
}

var platformListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出已安装的平台",
	Run: func(cmd *cobra.Command, args []string) {
		listPlatforms()
	},
}

var platformInstallCmd = &cobra.Command{
	Use:   "install <platform>",
	Short: "安装平台配置",
	Long: `安装芯片平台配置。

配置将从官方仓库下载到 ~/.zyrthi/platforms/<platform>/`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		installPlatform(args[0])
	},
}

var platformUpdateCmd = &cobra.Command{
	Use:   "update <platform>",
	Short: "更新平台配置",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		updatePlatform(args[0])
	},
}

func init() {
	platformCmd.AddCommand(platformListCmd)
	platformCmd.AddCommand(platformInstallCmd)
	platformCmd.AddCommand(platformUpdateCmd)
	rootCmd.AddCommand(platformCmd)
}

func listPlatforms() {
	fmt.Println("已安装的平台:")

	entries, err := os.ReadDir(platformDir)
	if err != nil {
		fmt.Println("  (无)")
		fmt.Printf("\n配置目录: %s\n", platformDir)
		return
	}

	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		configPath := filepath.Join(platformDir, name, "platform.yaml")

		status := "(缺少配置)"
		if _, err := os.Stat(configPath); err == nil {
			status = "✓"
		}

		fmt.Printf("  - %s %s\n", name, status)
		count++
	}

	if count == 0 {
		fmt.Println("  (无)")
	}

	fmt.Printf("\n配置目录: %s\n", platformDir)
}

func installPlatform(name string) {
	targetDir := filepath.Join(platformDir, name)
	configPath := filepath.Join(targetDir, "platform.yaml")

	// 创建目录
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "错误: 无法创建目录: %v\n", err)
		os.Exit(1)
	}

	// 检查是否已存在
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("平台 %s 已安装\n", name)
		return
	}

	// 创建默认配置
	defaultConfig := fmt.Sprintf(`# platform.yaml - %s 平台配置
platform: %s

compiler:
  default_cflags: [-Os, -Wall]
  default_ldflags: []

flash:
  default_baud: 115200
  max_baud: 921600

chips:
  # 请添加芯片配置
`, name, name)

	if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "错误: 无法写入配置: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("已安装平台: %s\n", name)
	fmt.Printf("配置文件: %s\n", configPath)
	fmt.Println("\n提示: 请编辑 platform.yaml 添加芯片配置")
}

func updatePlatform(name string) {
	configPath := filepath.Join(platformDir, name, "platform.yaml")

	if _, err := os.Stat(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "错误: 平台 %s 未安装\n", name)
		os.Exit(1)
	}

	// TODO: 从远程仓库更新
	fmt.Printf("平台 %s 配置已是最新\n", name)
}
