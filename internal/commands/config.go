package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// PlatformConfig platform.yaml 结构
type PlatformConfig struct {
	Platform string                `yaml:"platform"`
	Compiler PlatformCompiler      `yaml:"compiler"`
	Flash    PlatformFlash         `yaml:"flash"`
	Chips    map[string]ChipConfig `yaml:"chips"`
}

// PlatformCompiler 平台级编译器配置
type PlatformCompiler struct {
	DefaultCflags  []string `yaml:"default_cflags"`
	DefaultLdflags []string `yaml:"default_ldflags"`
}

// PlatformFlash 平台级烧录配置
type PlatformFlash struct {
	DefaultBaud int `yaml:"default_baud"`
	MaxBaud     int `yaml:"max_baud"`
}

// ChipConfig 芯片配置
type ChipConfig struct {
	Core      string       `yaml:"core"`
	FlashSize string       `yaml:"flash_size"`
	RamSize   string       `yaml:"ram_size"`
	Psram     bool         `yaml:"psram"`
	Compiler  ChipCompiler `yaml:"compiler"`
	Flash     ChipFlash    `yaml:"flash"`
}

// ChipCompiler 芯片级编译器配置
type ChipCompiler struct {
	Prefix string   `yaml:"prefix"`
	Cflags []string `yaml:"cflags"`
}

// ChipFlash 芯片级烧录配置
type ChipFlash struct {
	Plugin    string `yaml:"plugin"`
	EntryAddr string `yaml:"entry_addr"`
}

// ProjectConfig 生成的 zyrthi.yaml 结构
type ProjectConfig struct {
	Platform string           `yaml:"platform"`
	Chip     string           `yaml:"chip"`
	Compiler ProjectCompiler `yaml:"compiler"`
	Flash    ProjectFlash    `yaml:"flash"`
	Monitor  ProjectMonitor  `yaml:"monitor"`
	Project  ProjectMeta     `yaml:"project"`
}

type ProjectCompiler struct {
	Prefix   string   `yaml:"prefix"`
	Cflags   []string `yaml:"cflags"`
	Ldflags  []string `yaml:"ldflags"`
	Includes []string `yaml:"includes"`
}

type ProjectFlash struct {
	Plugin      string `yaml:"plugin"`
	EntryAddr   string `yaml:"entry_addr"`
	FlashSize   string `yaml:"flash_size"`
	DefaultBaud int    `yaml:"default_baud"`
}

type ProjectMonitor struct {
	Baud int `yaml:"baud"`
}

type ProjectMeta struct {
	Name    string   `yaml:"name"`
	Sources []string `yaml:"sources"`
}

// loadPlatformConfig 加载 platform.yaml
func loadPlatformConfig(platform string) (*PlatformConfig, error) {
	searchPaths := []string{
		filepath.Join(os.Getenv("HOME"), ".zyrthi", "platforms", platform, "platform.yaml"),
		filepath.Join("platforms", platform, "platform.yaml"),
		filepath.Join("..", "examples", "platform-"+platform+".yaml"),
	}

	var platformPath string
	for _, p := range searchPaths {
		if _, err := os.Stat(p); err == nil {
			platformPath = p
			break
		}
	}

	if platformPath == "" {
		return nil, fmt.Errorf("找不到平台配置: %s\n请先运行: zyrthi platform install %s", platform, platform)
	}

	data, err := os.ReadFile(platformPath)
	if err != nil {
		return nil, err
	}

	var cfg PlatformConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// generateProjectConfig 从平台配置生成项目配置
func generateProjectConfig(platform, chip string) (*ProjectConfig, error) {
	platformCfg, err := loadPlatformConfig(platform)
	if err != nil {
		return nil, err
	}

	chipCfg, ok := platformCfg.Chips[chip]
	if !ok {
		return nil, fmt.Errorf("平台 %s 不支持芯片 %s", platform, chip)
	}

	cfg := &ProjectConfig{
		Platform: platform,
		Chip:     chip,
		Compiler: ProjectCompiler{
			Prefix:   chipCfg.Compiler.Prefix,
			Cflags:   append(platformCfg.Compiler.DefaultCflags, chipCfg.Compiler.Cflags...),
			Ldflags:  platformCfg.Compiler.DefaultLdflags,
			Includes: []string{"src/"},
		},
		Flash: ProjectFlash{
			Plugin:      chipCfg.Flash.Plugin,
			EntryAddr:   chipCfg.Flash.EntryAddr,
			FlashSize:   chipCfg.FlashSize,
			DefaultBaud: platformCfg.Flash.DefaultBaud,
		},
		Monitor: ProjectMonitor{
			Baud: platformCfg.Flash.DefaultBaud,
		},
		Project: ProjectMeta{
			Name:    "my-project",
			Sources: []string{"src/"},
		},
	}

	return cfg, nil
}

// writeProjectConfig 写入 zyrthi.yaml
func writeProjectConfig(cfg *ProjectConfig, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	header := `# zyrthi.yaml - 项目配置
# 由 "zyrthi init" 自动生成

`

	return os.WriteFile(path, append([]byte(header), data...), 0644)
}
