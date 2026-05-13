package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestInitConfig(t *testing.T) {
	if err := LoadConfig(); err != nil {
		t.Error(err)
	}
	t.Logf("Config loaded successfully%s,%s,%s", viper.GetString("llm.baseurl"), viper.GetString("llm.apikey"), viper.GetString("llm.model"))
}
