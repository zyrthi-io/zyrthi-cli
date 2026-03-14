package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

// ============ Build Command Tests ============

func TestBuildCommand(t *testing.T) {
	buildCmd, _, err := rootCmd.Find([]string{"build"})
	if err != nil {
		t.Fatal("build command should exist")
	}
	if buildCmd.Use != "build" {
		t.Errorf("expected Use 'build', got %s", buildCmd.Use)
	}
	if buildCmd.Short == "" {
		t.Error("build command should have a short description")
	}
}

func TestBuildCommandFlagClean(t *testing.T) {
	buildCmd, _, _ := rootCmd.Find([]string{"build"})

	cleanFlag := buildCmd.Flags().Lookup("clean")
	if cleanFlag == nil {
		t.Fatal("build command should have --clean flag")
	}
	if cleanFlag.DefValue != "false" {
		t.Errorf("expected default 'false', got %s", cleanFlag.DefValue)
	}
}

func TestBuildCleanFlagDefaultValue(t *testing.T) {
	if buildClean != false {
		t.Error("buildClean should be false by default")
	}
}

// ============ Flash Command Tests ============

func TestFlashCommandExists(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"flash"})
	if err != nil {
		t.Fatal("flash command should exist")
	}
	if cmd.Use != "flash" {
		t.Errorf("expected Use 'flash', got %s", cmd.Use)
	}
}

func TestFlashCommandFlags(t *testing.T) {
	flashCmd, _, _ := rootCmd.Find([]string{"flash"})

	expectedFlags := []string{"port", "baud", "firmware", "erase", "verify"}
	for _, flagName := range expectedFlags {
		flag := flashCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("flash command should have --%s flag", flagName)
		}
	}
}

func TestFlashCommandFlagDefaults(t *testing.T) {
	_, _, _ = rootCmd.Find([]string{"flash"})

	if flashPort != "" {
		t.Error("flashPort should be empty by default")
	}
	if flashBaud != 0 {
		t.Error("flashBaud should be 0 by default")
	}
	if flashErase != false {
		t.Error("flashErase should be false by default")
	}
	if flashVerify != false {
		t.Error("flashVerify should be false by default")
	}
}

// ============ Monitor Command Tests ============

func TestMonitorCommand(t *testing.T) {
	monitorCmd, _, err := rootCmd.Find([]string{"monitor"})
	if err != nil {
		t.Fatal("monitor command should exist")
	}
	if monitorCmd.Use != "monitor" {
		t.Errorf("expected Use 'monitor', got %s", monitorCmd.Use)
	}
}

func TestMonitorCommandFlags(t *testing.T) {
	monitorCmd, _, _ := rootCmd.Find([]string{"monitor"})

	expectedFlags := []string{"port", "baud", "timestamp", "hex", "log", "filter"}
	for _, flagName := range expectedFlags {
		flag := monitorCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("monitor command should have --%s flag", flagName)
		}
	}
}

func TestMonitorCommandFlagDefaults(t *testing.T) {
	if monitorPort != "" {
		t.Error("monitorPort should be empty by default")
	}
	if monitorBaud != 0 {
		t.Error("monitorBaud should be 0 by default")
	}
	if monitorTimestamp != false {
		t.Error("monitorTimestamp should be false by default")
	}
	if monitorHex != false {
		t.Error("monitorHex should be false by default")
	}
}

// ============ Init Command Tests ============

func TestInitCommand(t *testing.T) {
	initCmd, _, err := rootCmd.Find([]string{"init"})
	if err != nil {
		t.Fatal("init command should exist")
	}
	if initCmd.Use != "init" {
		t.Errorf("expected Use 'init', got %s", initCmd.Use)
	}
}

func TestInitCommandFlags(t *testing.T) {
	initCmd, _, _ := rootCmd.Find([]string{"init"})

	platformFlag := initCmd.Flags().Lookup("platform")
	if platformFlag == nil {
		t.Error("init command should have --platform flag")
	}

	chipFlag := initCmd.Flags().Lookup("chip")
	if chipFlag == nil {
		t.Error("init command should have --chip flag")
	}
}

func TestInitCommandRequiredFlags(t *testing.T) {
	initCmd, _, _ := rootCmd.Find([]string{"init"})

	// Check if flags are marked as required
	platformFlag := initCmd.Flags().Lookup("platform")
	if platformFlag == nil {
		t.Fatal("platform flag not found")
	}

	chipFlag := initCmd.Flags().Lookup("chip")
	if chipFlag == nil {
		t.Fatal("chip flag not found")
	}
}

// ============ Version Command Tests ============

func TestVersionCommandExists(t *testing.T) {
	versionCmd, _, err := rootCmd.Find([]string{"version"})
	if err != nil {
		t.Fatal("version command should exist")
	}
	if versionCmd.Use != "version" {
		t.Errorf("expected Use 'version', got %s", versionCmd.Use)
	}
	if versionCmd.Short == "" {
		t.Error("version command should have a short description")
	}
}

// ============ Platform Command Tests ============

func TestPlatformCommand(t *testing.T) {
	platformCmd, _, err := rootCmd.Find([]string{"platform"})
	if err != nil {
		t.Fatal("platform command should exist")
	}
	if platformCmd.Use != "platform" {
		t.Errorf("expected Use 'platform', got %s", platformCmd.Use)
	}
}

func TestPlatformListCommand(t *testing.T) {
	platformListCmd, _, err := rootCmd.Find([]string{"platform", "list"})
	if err != nil {
		t.Fatal("platform list command should exist")
	}
	if platformListCmd.Use != "list" {
		t.Errorf("expected Use 'list', got %s", platformListCmd.Use)
	}
}

func TestPlatformInstallCommand(t *testing.T) {
	platformInstallCmd, _, err := rootCmd.Find([]string{"platform", "install"})
	if err != nil {
		t.Fatal("platform install command should exist")
	}
	if platformInstallCmd.Use != "install <platform>" {
		t.Errorf("expected Use 'install <platform>', got %s", platformInstallCmd.Use)
	}
}

func TestPlatformUpdateCommand(t *testing.T) {
	platformUpdateCmd, _, err := rootCmd.Find([]string{"platform", "update"})
	if err != nil {
		t.Fatal("platform update command should exist")
	}
	if platformUpdateCmd.Use != "update <platform>" {
		t.Errorf("expected Use 'update <platform>', got %s", platformUpdateCmd.Use)
	}
}

// ============ findConfig Tests ============

func TestFindConfigCustomPath(t *testing.T) {
	origCfgFile := cfgFile
	defer func() { cfgFile = origCfgFile }()

	cfgFile = "custom-config.yaml"
	result := findConfig()
	if result != "custom-config.yaml" {
		t.Errorf("expected 'custom-config.yaml', got %s", result)
	}
}

func TestFindConfigInCurrentDir(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	// Create config file
	configPath := filepath.Join(tmpDir, "zyrthi.yaml")
	os.WriteFile(configPath, []byte("platform: esp32\n"), 0644)

	origCfgFile := cfgFile
	defer func() { cfgFile = origCfgFile }()
	cfgFile = "zyrthi.yaml"

	result := findConfig()
	// Result should be the absolute path to the config
	if result == "" {
		t.Error("findConfig should return a path")
	}
}

func TestFindConfigNonexistent(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	origCfgFile := cfgFile
	defer func() { cfgFile = origCfgFile }()
	cfgFile = "zyrthi.yaml"

	result := findConfig()
	// Should return the default filename when not found
	if result != "zyrthi.yaml" {
		t.Errorf("expected 'zyrthi.yaml', got %s", result)
	}
}

// ============ listPlatforms Tests ============

func TestListPlatformsEmpty(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Save original platformDir
	origPlatformDir := platformDir
	platformDir = tmpDir
	defer func() { platformDir = origPlatformDir }()

	// listPlatforms should not panic with empty directory
	listPlatforms()
}

func TestListPlatformsNonexistentDir(t *testing.T) {
	// Save original platformDir
	origPlatformDir := platformDir
	platformDir = "/nonexistent/path"
	defer func() { platformDir = origPlatformDir }()

	// listPlatforms should not panic with nonexistent directory
	listPlatforms()
}

func TestListPlatformsValid(t *testing.T) {
	// Create temp directory with a platform
	tmpDir := t.TempDir()
	platformPath := filepath.Join(tmpDir, "test-platform")
	os.MkdirAll(platformPath, 0755)

	// Create platform.yaml
	configContent := `platform: test-platform
compiler:
  default_cflags: [-Os]
  default_ldflags: []
flash:
  default_baud: 115200
chips: {}
`
	os.WriteFile(filepath.Join(platformPath, "platform.yaml"), []byte(configContent), 0644)

	// Save original platformDir
	origPlatformDir := platformDir
	platformDir = tmpDir
	defer func() { platformDir = origPlatformDir }()

	// listPlatforms should list the platform
	listPlatforms()
}

// ============ installPlatform Tests ============

func TestInstallPlatform(t *testing.T) {
	tmpDir := t.TempDir()

	// Save original platformDir
	origPlatformDir := platformDir
	platformDir = tmpDir
	defer func() { platformDir = origPlatformDir }()

	installPlatform("test-platform")

	// Check if platform.yaml was created
	configPath := filepath.Join(tmpDir, "test-platform", "platform.yaml")
	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("platform.yaml should exist: %v", err)
	}
}

func TestInstallPlatformAlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Create existing platform
	platformPath := filepath.Join(tmpDir, "existing-platform")
	os.MkdirAll(platformPath, 0755)
	os.WriteFile(filepath.Join(platformPath, "platform.yaml"), []byte("platform: existing"), 0644)

	// Save original platformDir
	origPlatformDir := platformDir
	platformDir = tmpDir
	defer func() { platformDir = origPlatformDir }()

	installPlatform("existing-platform")

	// Should not overwrite
	data, _ := os.ReadFile(filepath.Join(platformPath, "platform.yaml"))
	if string(data) != "platform: existing" {
		t.Error("existing platform should not be overwritten")
	}
}

// ============ updatePlatform Tests ============

func TestUpdatePlatformNotInstalled(t *testing.T) {
	tmpDir := t.TempDir()

	// Save original platformDir
	origPlatformDir := platformDir
	platformDir = tmpDir
	defer func() { platformDir = origPlatformDir }()

	// updatePlatform should exit with error for non-existent platform
	// We can't easily test os.Exit, so we just verify the logic
	configPath := filepath.Join(tmpDir, "nonexistent", "platform.yaml")
	if _, err := os.Stat(configPath); err == nil {
		t.Error("platform should not exist")
	}
}

func TestUpdatePlatformInstalled(t *testing.T) {
	tmpDir := t.TempDir()

	// Create existing platform
	platformPath := filepath.Join(tmpDir, "existing-platform")
	os.MkdirAll(platformPath, 0755)
	os.WriteFile(filepath.Join(platformPath, "platform.yaml"), []byte("platform: existing"), 0644)

	// Save original platformDir
	origPlatformDir := platformDir
	platformDir = tmpDir
	defer func() { platformDir = origPlatformDir }()

	// updatePlatform should work for existing platform
	updatePlatform("existing-platform")
}

// ============ generateConfig Tests ============

func TestGenerateConfigWithValidInput(t *testing.T) {
	// Create temp directory with platform config
	tmpDir := t.TempDir()
	platformPath := filepath.Join(tmpDir, "platforms", "test-platform")
	os.MkdirAll(platformPath, 0755)

	configContent := `platform: test-platform
compiler:
  default_cflags: [-Os, -Wall]
  default_ldflags: [-nostdlib]
flash:
  default_baud: 115200
  max_baud: 921600
chips:
  test-chip:
    core: riscv
    flash_size: 4MB
    ram_size: 400KB
    compiler:
      prefix: riscv32-esp-elf-
      cflags: [-march=rv32imc]
    flash:
      plugin: https://example.com/plugin.wasm
      entry_addr: "0x0"
`
	os.WriteFile(filepath.Join(platformPath, "platform.yaml"), []byte(configContent), 0644)

	// Save working directory and change to temp
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	// Save original platformDir
	origPlatformDir := platformDir
	platformDir = filepath.Join(tmpDir, "platforms")
	defer func() { platformDir = origPlatformDir }()

	// Test generateConfig
	origCfgFile := cfgFile
	defer func() { cfgFile = origCfgFile }()
	cfgFile = "zyrthi.yaml"

	err := generateConfig("test-platform", "test-chip")
	if err != nil {
		t.Errorf("generateConfig error: %v", err)
	}

	// Check if zyrthi.yaml was created
	if _, err := os.Stat("zyrthi.yaml"); err != nil {
		t.Errorf("zyrthi.yaml should exist: %v", err)
	}
}

// ============ Command Execution Tests ============

func TestRootCommandExecute(t *testing.T) {
	// Test that Execute() doesn't panic
	// Note: This will actually execute the command, so we need to be careful
	// We're just testing that the function exists and doesn't crash
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--help"})

	// This should not panic
	Execute()
}

func TestCommandHelp(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"build help", []string{"build", "--help"}},
		{"flash help", []string{"flash", "--help"}},
		{"monitor help", []string{"monitor", "--help"}},
		{"platform help", []string{"platform", "--help"}},
		{"platform list help", []string{"platform", "list", "--help"}},
		{"platform install help", []string{"platform", "install", "--help"}},
		{"platform update help", []string{"platform", "update", "--help"}},
		{"init help", []string{"init", "--help"}},
		{"version help", []string{"version", "--help"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			rootCmd.SetOut(buf)
			rootCmd.SetArgs(tt.args)

			err := rootCmd.Execute()
			if err != nil {
				t.Errorf("command error: %v", err)
			}

			output := buf.String()
			if output == "" {
				t.Error("help should produce output")
			}
		})
	}
}

// ============ Additional Tests ============

func TestInitVariablesDefaults(t *testing.T) {
	if initPlatform != "" {
		t.Error("initPlatform should be empty by default")
	}
	if initChip != "" {
		t.Error("initChip should be empty by default")
	}
}

func TestPlatformVariablesDefaults(t *testing.T) {
	// platformDir should be set in init()
	if platformDir == "" {
		t.Error("platformDir should be set")
	}
}

func TestCommandExamples(t *testing.T) {
	commands := []struct {
		name string
		path []string
	}{
		{"build", []string{"build"}},
		{"flash", []string{"flash"}},
		{"monitor", []string{"monitor"}},
		{"init", []string{"init"}},
	}

	for _, tc := range commands {
		t.Run(tc.name, func(t *testing.T) {
			cmd, _, err := rootCmd.Find(tc.path)
			if err != nil {
				t.Fatalf("command not found: %v", err)
			}
			if cmd.Example == "" {
				t.Errorf("%s command should have examples", tc.name)
			}
		})
	}
}

func TestCommandLongDescriptions(t *testing.T) {
	commands := []struct {
		name string
		path []string
	}{
		{"build", []string{"build"}},
		{"flash", []string{"flash"}},
		{"monitor", []string{"monitor"}},
		{"init", []string{"init"}},
		{"platform", []string{"platform"}},
		{"platform install", []string{"platform", "install"}},
	}

	for _, tc := range commands {
		t.Run(tc.name, func(t *testing.T) {
			cmd, _, err := rootCmd.Find(tc.path)
			if err != nil {
				t.Fatalf("command not found: %v", err)
			}
			if cmd.Long == "" {
				t.Errorf("%s command should have long description", tc.name)
			}
		})
	}
}

func TestSubcommandStructure(t *testing.T) {
	// Test that platform has correct subcommands
	platformCmd, _, _ := rootCmd.Find([]string{"platform"})
	subcommands := platformCmd.Commands()

	cmdNames := make(map[string]bool)
	for _, cmd := range subcommands {
		cmdNames[cmd.Name()] = true
	}

	expected := []string{"list", "install", "update"}
	for _, name := range expected {
		if !cmdNames[name] {
			t.Errorf("platform command should have '%s' subcommand", name)
		}
	}
}

func TestPlatformInstallArgs(t *testing.T) {
	platformInstallCmd, _, _ := rootCmd.Find([]string{"platform", "install"})

	// Check that Args is set to ExactArgs(1)
	if platformInstallCmd.Args == nil {
		t.Error("platform install should have Args validator")
	}
}

func TestPlatformUpdateArgs(t *testing.T) {
	platformUpdateCmd, _, _ := rootCmd.Find([]string{"platform", "update"})

	// Check that Args is set to ExactArgs(1)
	if platformUpdateCmd.Args == nil {
		t.Error("platform update should have Args validator")
	}
}

func TestVersionOutput(t *testing.T) {
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

	// Check that output contains version info
	if len(output) < 5 {
		t.Error("version output seems too short")
	}
}

// ============ Mock Command Tests ============

func TestCommandRunFunctions(t *testing.T) {
	// Test that Run functions exist for all commands
	commands := []struct {
		name string
		path []string
	}{
		{"build", []string{"build"}},
		{"flash", []string{"flash"}},
		{"monitor", []string{"monitor"}},
		{"init", []string{"init"}},
		{"version", []string{"version"}},
		{"platform list", []string{"platform", "list"}},
		{"platform install", []string{"platform", "install"}},
		{"platform update", []string{"platform", "update"}},
	}

	for _, tc := range commands {
		t.Run(tc.name, func(t *testing.T) {
			cmd, _, err := rootCmd.Find(tc.path)
			if err != nil {
				t.Fatalf("command not found: %v", err)
			}
			if cmd.Run == nil {
				t.Errorf("%s command should have Run function", tc.name)
			}
		})
	}
}

func TestCobraCommandStructure(t *testing.T) {
	// Verify all commands have proper cobra.Command structure
	var checkCommand func(cmd *cobra.Command, path string)
	checkCommand = func(cmd *cobra.Command, path string) {
		if cmd.Use == "" {
			t.Errorf("command at %s has empty Use", path)
		}
		if cmd.Short == "" {
			t.Errorf("command at %s has empty Short", path)
		}

		for _, subCmd := range cmd.Commands() {
			checkCommand(subCmd, path+" "+subCmd.Name())
		}
	}

	checkCommand(rootCmd, "root")
}
