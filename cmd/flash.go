package cmd

import (
	"github.com/spf13/cobra"
)

var (
	flashPort     string
	flashBaud     int
	flashFirmware string
	flashErase    bool
	flashVerify   bool
)

var flashCmd = &cobra.Command{
	Use:   "flash",
	Short: "烧录固件",
	Long:  `将固件烧录到目标设备。支持自动检测芯片和串口。`,
	Example: `  zyrthi flash
  zyrthi flash --port /dev/ttyUSB0
  zyrthi flash --baud 921600
  zyrthi flash --erase`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := findConfig()
		flashArgs := []string{"--config", configPath}

		if flashPort != "" {
			flashArgs = append(flashArgs, "--port", flashPort)
		}
		if flashBaud > 0 {
			flashArgs = append(flashArgs, "--baud", string(rune(flashBaud)))
		}
		if flashFirmware != "" {
			flashArgs = append(flashArgs, "--firmware", flashFirmware)
		}
		if flashErase {
			flashArgs = append(flashArgs, "--erase")
		}
		if flashVerify {
			flashArgs = append(flashArgs, "--verify")
		}

		execCommand("zyrthi-flash", flashArgs)
	},
}

func init() {
	flashCmd.Flags().StringVarP(&flashPort, "port", "p", "", "串口设备")
	flashCmd.Flags().IntVarP(&flashBaud, "baud", "b", 0, "波特率")
	flashCmd.Flags().StringVarP(&flashFirmware, "firmware", "f", "", "固件文件路径")
	flashCmd.Flags().BoolVar(&flashErase, "erase", false, "烧录前全片擦除")
	flashCmd.Flags().BoolVar(&flashVerify, "verify", false, "烧录后校验")

	rootCmd.AddCommand(flashCmd)
}