package internal

import (
	"fmt"
	"io/ioutil"
	"os"
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
		Sender string `yaml:"sender"`
	}

	SenderSet map[string]map[string]interface{}
)

func LoadConfigFile(configPath string) *ConfigTemplate {
	cp := strings.Replace(configPath, "yaml", "yml", 1)
	f, err := os.Open(cp)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

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
