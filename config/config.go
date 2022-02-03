package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Runtime struct {
	MaxCpus uint `yaml:"maxCpus"`
}

type Core struct {
	JobCleanThreshold     uint   `yaml:"jobCleanThreshold"`
	FileTransferSizeLimit string `yaml:"fileTransferSizeLimit"`
}

type Auth struct {
	Method     string                 `yaml:"method"`
	Credential map[string]interface{} `yaml:"credential"`
}

type Http struct {
	Ip   string `yaml:"ip"`
	Port uint   `yaml:"port"`
	Auth Auth   `yaml:"auth"`
}

type Log struct {
	Level         string `yaml:"level"`
	FileSizeLimit string `yaml:"fileSizeLimit"`
}

type Config struct {
	Runtime Runtime `yaml:"runtime"`
	Core    Core    `yaml:"core"`
	Http    Http    `yaml:"http"`
	Log     Log     `yaml:"log"`
}

func LoadConfig(filePath string) (config Config, err error) {
	// Read file
	var content []byte
	content, err = ioutil.ReadFile(filePath)
	if nil != err {
		return
	}

	// Unmarshal
	err = yaml.UnmarshalStrict(content, &config)

	return
}
