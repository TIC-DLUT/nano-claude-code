package config

import (
	"fmt"
)

type Config struct {
	LLM struct {
		BaseURL string `json:"baseurl"`
		APIKey  string `json:"apikey"`
		Model   string `json:"model"`
	} `json:"llm"`
}

type configInput struct {
	Name         string
	Prompt       string  // 提示语
	DefaultValue string  // 默认值
	Target       *string // 要赋值的变量指针
	Required     bool
}

func newConfigInputInit(c *Config) error {
	inputs := []configInput{
		{
			Name:         "Base URL",
			Prompt:       "Please enter API Base URL (press Enter for official API)",
			DefaultValue: "https://api.anthropic.com",
			Target:       &c.LLM.BaseURL,
			Required:     false,
		},
		{
			Name:     "API Key",
			Prompt:   "Please enter API Key",
			Target:   &c.LLM.APIKey,
			Required: true,
		},
		{
			Name:     "Model ID",
			Prompt:   "Please enter model ID",
			Target:   &c.LLM.Model,
			Required: true,
		},
	}

	for _, input := range inputs {
		fmt.Printf("%s: ", input.Prompt)
		fmt.Scanf("%s", input.Target)
		if input.Required {
			// 必填参数，用户输入为空时重复让用户输入
			for {
				if *input.Target != "" {
					break
				}
				fmt.Printf("%s cannot be empty, please try again\n", input.Name)
				fmt.Printf("%s: ", input.Prompt)
				fmt.Scanf("%s", input.Target)
			}
		} else {
			// 非必填参数，使用默认值
			//使用默认值
			if *input.Target == "" && input.DefaultValue != "" {
				*input.Target = input.DefaultValue
				fmt.Printf("Using default value: %s\n", input.DefaultValue)
			}
		}
	}
	return nil
}
