package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`

	Domains struct {
		Transfer360 struct {
			Sandbox string `yaml:"sandbox"`
		} `yaml:"transfer360"`
	} `yaml:"domains"`

	Endpoints struct {
		TestSearch struct {
			AcmeLease    string `yaml:"acmelease"`
			LeaseCompany string `yaml:"leasecompany"`
			FleetCompany string `yaml:"fleetcompany"`
			HireCompany  string `yaml:"hirecompany"`
		} `yaml:"test_search"`
	} `yaml:"endpoints"`

	Google struct {
		UsingCloud         bool   `yaml:"usingCloud"`
		ServiceAccountPath string `yaml:"serviceAccountPath"`
		PubSubTopic        string `yaml:"pubSubTopic"`
	} `yaml:"google"`

	Emulator struct {
		Host        string `yaml:"host"`
		Port        int    `yaml:"port"`
		ProjectId   string `yaml:"projectId"`
		PubSubTopic string `yaml:"pubSubTopic"`
	} `yaml:"local_emulator"`
}

// LoadConfig loads the configuration from the YAML file and replaces placeholders with actual values
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	cfg.Endpoints.TestSearch.AcmeLease = replacePlaceholder(cfg.Endpoints.TestSearch.AcmeLease, cfg.Domains.Transfer360.Sandbox)
	cfg.Endpoints.TestSearch.LeaseCompany = replacePlaceholder(cfg.Endpoints.TestSearch.LeaseCompany, cfg.Domains.Transfer360.Sandbox)
	cfg.Endpoints.TestSearch.FleetCompany = replacePlaceholder(cfg.Endpoints.TestSearch.FleetCompany, cfg.Domains.Transfer360.Sandbox)
	cfg.Endpoints.TestSearch.HireCompany = replacePlaceholder(cfg.Endpoints.TestSearch.HireCompany, cfg.Domains.Transfer360.Sandbox)

	return &cfg, nil
}

func replacePlaceholder(endpoint string, sandboxURL string) string {
	return strings.Replace(endpoint, "${domains.transfer360.sandbox}", sandboxURL, -1)
}
