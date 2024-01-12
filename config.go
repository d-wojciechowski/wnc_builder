package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
)

const CFG_FILE_CONTENT = `profile: prod
root: /opt
fail_on_error: false
commands:
  ootb:
    restart: echo "Restarting"
  custom:
    full: echo "Full"
input:
  build_order: ignored/compile.includes
  module_registry: ignored/moduleRegistry.xml
aliases:
  mpml: MPMLink
  mpmlc: MPMLinkCommon
  ppb: ProcessPlanBrowser
  ass: Associative
suites:
  current:
    restart: true
    build:
      MPMLink: cst
      MPMLinkCommon: cst
      ProcessPlanBrowser: cst
`

type Suite struct {
	Restart bool
	Build   *map[string]string
	Custom  []string
}

type OOTBCommands struct {
	Restart string
}

type Commands struct {
	OOTB   OOTBCommands
	Custom map[string]string
}

type Input struct {
	BuildOrder     string `yaml:"build_order"`
	ModuleRegistry string `yaml:"module_registry"`
}

type AppConfig struct {
	Profile     string
	Root        string
	FailOnError bool `yaml:"fail_on_error"`
	Commands    Commands
	Input       Input
	Aliases     map[string]string
	Suites      map[string]Suite
}

func createAppConfig() (*AppConfig, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("user HomeDirectory is not available. %w", err)
	}
	var appConfigDir = fmt.Sprintf("%s/.wc_builder", dir)
	var configPath = fmt.Sprintf("%s/cfg.yml", appConfigDir)

	_, err = os.Stat(configPath)
	if errors.Is(err, fs.ErrNotExist) {
		err := createFileWhenConfigMissing(appConfigDir, configPath)
		if err != nil {
			return nil, err
		}
	}

	bytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not read configuration content. %w", err)
	}

	c := &AppConfig{}
	err = yaml.Unmarshal(bytes, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", configPath, err)
	}

	return c, err
}

func createFileWhenConfigMissing(appConfigDir string, configPath string) error {
	err := os.Mkdir(appConfigDir, os.ModePerm)
	if errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("could not create configuration directory in user home dir. %w", err)
	}
	configFile, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("could not create configuration file in config dir. %w", err)
	}
	_, err = configFile.WriteString(CFG_FILE_CONTENT)
	if err != nil {
		return fmt.Errorf("could not write configuration to config file. %w", err)
	}
	return nil
}
