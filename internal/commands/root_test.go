package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRootCommandExists(t *testing.T) {
	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}
	if rootCmd.Use != "zyrthi" {
		t.Errorf("expected Use 'zyrthi', got %s", rootCmd.Use)
	}
}

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
}

func TestVersionCommand(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"version"})

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("version command error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("version command should produce output")
	}
}

func TestHelpCommand(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("help command error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("help command should produce output")
	}
}

func TestBuildCommandExists(t *testing.T) {
	buildCmd, _, err := rootCmd.Find([]string{"build"})
	if err != nil {
		t.Fatal("build command should exist")
	}
	if buildCmd == nil {
		t.Fatal("build command should not be nil")
	}
	if buildCmd.Short == "" {
		t.Error("build command should have a short description")
	}
}

func TestBuildCommandFlags(t *testing.T) {
	buildCmd, _, _ := rootCmd.Find([]string{"build"})
	
	cleanFlag := buildCmd.Flags().Lookup("clean")
	if cleanFlag == nil {
		t.Error("build command should have --clean flag")
	}
}

func TestFindConfigDefault(t *testing.T) {
	// 测试默认配置文件查找
	cfg := findConfig()
	if cfg == "" {
		t.Error("findConfig should return a path")
	}
}

func TestFindConfigCustom(t *testing.T) {
	// 保存原始值
	origCfgFile := cfgFile
	defer func() { cfgFile = origCfgFile }()

	cfgFile = "custom.yaml"
	cfg := findConfig()
	if cfg != "custom.yaml" {
		t.Errorf("expected 'custom.yaml', got %s", cfg)
	}
}

func TestFindConfigInDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	// 创建配置文件
	configPath := filepath.Join(tmpDir, "zyrthi.yaml")
	os.WriteFile(configPath, []byte("platform: esp32\n"), 0644)

	// 保存原始值
	origCfgFile := cfgFile
	defer func() { cfgFile = origCfgFile }()
	cfgFile = "zyrthi.yaml"

	cfg := findConfig()
	// 使用 filepath.EvalSymlinks 处理符号链接
	expectedPath, _ := filepath.EvalSymlinks(configPath)
	actualPath, _ := filepath.EvalSymlinks(cfg)
	
	if actualPath != expectedPath {
		t.Errorf("expected %s, got %s", expectedPath, actualPath)
	}
}

func TestRootCmdHasSubCommands(t *testing.T) {
	commands := rootCmd.Commands()
	cmdNames := make(map[string]bool)
	for _, cmd := range commands {
		cmdNames[cmd.Name()] = true
	}

	expected := []string{"build", "version", "init", "flash", "monitor", "platform"}
	for _, name := range expected {
		if !cmdNames[name] {
			t.Errorf("root command should have '%s' subcommand", name)
		}
	}
}
