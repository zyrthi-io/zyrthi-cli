package cmd

import (
	"github.com/spf13/cobra"
)

var (
	monitorPort      string
	monitorBaud      int
	monitorTimestamp bool
	monitorHex       bool
	monitorLog       string
	monitorFilter    string
)

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "串口监控",
	Long:  `启动串口监控，实时显示设备输出。`,
	Example: `  zyrthi monitor
  zyrthi monitor --port /dev/ttyUSB0 --baud 115200
  zyrthi monitor --timestamp --log output.txt`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := findConfig()
		monArgs := []string{"--config", configPath}

		if monitorPort != "" {
			monArgs = append(monArgs, "--port", monitorPort)
		}
		if monitorBaud > 0 {
			monArgs = append(monArgs, "--baud", string(rune(monitorBaud)))
		}
		if monitorTimestamp {
			monArgs = append(monArgs, "--timestamp")
		}
		if monitorHex {
			monArgs = append(monArgs, "--hex")
		}
		if monitorLog != "" {
			monArgs = append(monArgs, "--log", monitorLog)
		}
		if monitorFilter != "" {
			monArgs = append(monArgs, "--filter", monitorFilter)
		}

		execCommand("zyrthi-monitor", monArgs)
	},
}

func init() {
	monitorCmd.Flags().StringVarP(&monitorPort, "port", "p", "", "串口设备")
	monitorCmd.Flags().IntVarP(&monitorBaud, "baud", "b", 0, "波特率")
	monitorCmd.Flags().BoolVar(&monitorTimestamp, "timestamp", false, "显示时间戳")
	monitorCmd.Flags().BoolVar(&monitorHex, "hex", false, "十六进制显示")
	monitorCmd.Flags().StringVar(&monitorLog, "log", "", "日志保存文件")
	monitorCmd.Flags().StringVar(&monitorFilter, "filter", "", "过滤关键字")

	rootCmd.AddCommand(monitorCmd)
}