package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type (
	ConfigTemplate struct {
		AWSAccountSet AWSAccountSet `yaml:"accounts"`
		SenderSet     SenderSet     `yaml:"senders"`
	}

	AWSAccountSet map[string]AWSAccount

	AWSAccount struct {
		Name       string   `yaml:"name"`
		Exclusions []string `yaml:"exclusions"`
		Sender     string   `yaml:"sender"`
		Profile    string   `yaml:"profile"`
		Region     string   `yaml:"region"`
	}

	SenderSet map[string]map[string]interface{}
)

func loadConfigFile(configPath string) *ConfigTemplate {
	cp := strings.Replace(configPath, "yaml", "yml", 1)
	f, err := os.Open(filepath.Clean(cp))
	if err != nil {
		fmt.Println(err)
	}
	defer func() { _ = f.Close() }()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
	}

	var ct ConfigTemplate
	err = yaml.Unmarshal(b, &ct)
	if err != nil {
		fmt.Println(err)
	}
	return &ct
}
