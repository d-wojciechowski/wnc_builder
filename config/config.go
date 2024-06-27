package config

import (
	"errors"
	"fmt"
	"github.com/alexflint/go-arg"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
)

const CfgFileContent = `profile: prod
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
`

func ParseCmdArgs() *ProgramArguments {
	args := &ProgramArguments{}
	arg.MustParse(args)
	return args
}

type ProgramArguments struct {
	Build           []string `arg:"-b,--build" help:"Execute build "`
	TestUnit        []string `arg:"-u,--test-unit" help:"Execute [unit tests] / [unit test by name]"`
	TestIntegration []string `arg:"-i,--test-integration" help:"Execute [integ tests] / [integ test by name]"`
	Custom          []string `arg:"-c,--custom" help:"Execute custom command defined in CFG"`
	NumKey          []string `arg:"-n,--num-key" help:"Execute numkey build"`
	Restart         bool     `arg:"-r,--restart" help:"Execute restart"`
	Dry             bool     `arg:"-d,--dry" help:"Just generate commands."`
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
}

func CreateAppConfig() (*AppConfig, error) {
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
	_, err = configFile.WriteString(CfgFileContent)
	if err != nil {
		return fmt.Errorf("could not write configuration to config file. %w", err)
	}
	return nil
}
