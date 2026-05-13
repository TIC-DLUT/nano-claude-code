package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	customErrors "github.com/TIC-DLUT/nano-claude-code/errors"

	"github.com/spf13/viper"
)

func LoadConfig() error {
	homePath, _ := os.UserHomeDir()

	viper.SetConfigName("config")
	viper.AddConfigPath(filepath.Join(homePath, ".nano-claude-code"))

	BindEnv()

	err := viper.ReadInConfig()

	if err != nil {
		var configFileNotExit viper.ConfigFileNotFoundError
		// 配置文件未找到，引导新建
		if errors.As(err, &configFileNotExit) {
			if err := newConfigFile(filepath.Join(homePath, ".nano-claude-code")); err != nil {
				return err
			} else {
				// 配置文件新建成功，重新读取配置
				return viper.ReadInConfig()
			}
		}
		return customErrors.ReadInConfigError
	}
	return nil
}

func newConfigFile(configDir string) error {
	fmt.Println("Welcome to nano-claude-code!")
	fmt.Println("First time use requires configuring LLM parameters")
	fmt.Println()

	var config Config

	newConfigInputInit(&config)

	//构建文本内容
	configContent, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return customErrors.ConfigFileWriteError
	}

	// 创建配置目录
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return customErrors.ConfigDirCreateError
	}

	// 写入配置文件
	configPath := filepath.Join(configDir, "config.json")

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return customErrors.ConfigFileWriteError
	}

	fmt.Println()
	fmt.Println("Configuration file saved to:", configPath)
	return nil
}
