package cmd

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var buildClean bool

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "编译项目",
	Long:  `读取 zyrthi.yaml 配置，调用编译器编译项目。`,
	Example: `  zyrthi build
  zyrthi build --config my-config.yaml
  zyrthi build --clean`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := findConfig()
		buildArgs := []string{"--config", configPath}
		if buildClean {
			buildArgs = append(buildArgs, "--clean")
		}
		execCommand("zyrthi-build", buildArgs)
	},
}

func init() {
	buildCmd.Flags().BoolVar(&buildClean, "clean", false, "清理编译产物")
	rootCmd.AddCommand(buildCmd)
}

func execCommand(name string, args []string) {
	c := exec.Command(name, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}

func findConfig() string {
	if cfgFile != "zyrthi.yaml" {
		return cfgFile
	}

	cwd, _ := os.Getwd()
	configPath := filepath.Join(cwd, cfgFile)
	if _, err := os.Stat(configPath); err == nil {
		return configPath
	}

	dir := cwd
	for {
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		configPath := filepath.Join(parent, cfgFile)
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
		dir = parent
	}

	return cfgFile
}