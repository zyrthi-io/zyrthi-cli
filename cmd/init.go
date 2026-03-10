package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	initPlatform string
	initChip     string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化项目，生成 zyrthi.yaml",
	Long: `初始化项目，从平台配置生成 zyrthi.yaml。

会从 ~/.zyrthi/platforms/<platform>/platform.yaml 提取芯片配置。`,
	Example: `  zyrthi init --platform esp32 --chip esp32s3
  zyrthi init --platform esp32 --chip esp32`,
	Run: func(cmd *cobra.Command, args []string) {
		if initPlatform == "" || initChip == "" {
			fmt.Fprintln(os.Stderr, "错误: 必须指定 --platform 和 --chip")
			cmd.Help()
			os.Exit(1)
		}

		if err := generateConfig(initPlatform, initChip); err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("已生成 %s\n", cfgFile)
	},
}

func init() {
	initCmd.Flags().StringVarP(&initPlatform, "platform", "p", "", "平台名称 (必填)")
	initCmd.Flags().StringVarP(&initChip, "chip", "c", "", "芯片名称 (必填)")
	initCmd.MarkFlagRequired("platform")
	initCmd.MarkFlagRequired("chip")

	rootCmd.AddCommand(initCmd)
}

func generateConfig(platform, chip string) error {
	cfg, err := generateProjectConfig(platform, chip)
	if err != nil {
		return err
	}
	return writeProjectConfig(cfg, cfgFile)
}