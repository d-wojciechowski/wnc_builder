package executor

import (
	"errors"
	"fmt"
	"strings"
	"wnc_builder/config"
	"wnc_builder/module"
)

type taskBuilder struct {
	appConfig     *config.AppConfig
	modulesConfig map[string]*module.ModuleInfo
}

type TaskBuilder interface {
	BuildTasks(arguments *config.ProgramArguments) ([]*Task, error)
}

func NewTaskBuilder(appConfig *config.AppConfig, modulesConfig map[string]*module.ModuleInfo) TaskBuilder {
	builder := taskBuilder{
		appConfig:     appConfig,
		modulesConfig: modulesConfig,
	}
	return &builder
}

func (tb *taskBuilder) buildExplicitTasks(arguments *config.ProgramArguments) ([]*Task, error) {
	tasks := make([]*Task, 0, 1)
	if arguments.Build != nil && len(arguments.Build) > 0 {
		for _, moduleSpec := range arguments.Build {
			moduleInfo, targets, err := tb.getTaskSpec(moduleSpec)
			if err != nil {
				return nil, err
			}
			task := Task{
				Target:  config.Build,
				Module:  moduleInfo,
				targets: targets,
			}
			task.Commands = tb.createBuildCommands(task)
			tasks = append(tasks, &task)
		}
	}
	if arguments.TestUnit != nil && len(arguments.TestUnit) > 0 {
		for _, moduleSpec := range arguments.TestUnit {
			moduleInfo, targets, err := tb.getTaskSpec(moduleSpec)
			if err != nil {
				return nil, err
			}
			task := Task{
				Target:  config.TestUnit,
				Module:  moduleInfo,
				targets: targets,
			}
			task.Commands = []*Command{tb.createTestCommands(task)}
			tasks = append(tasks, &task)

		}
	}
	if arguments.TestIntegration != nil && len(arguments.TestIntegration) > 0 {
		for _, moduleSpec := range arguments.TestIntegration {
			moduleInfo, targets, err := tb.getTaskSpec(moduleSpec)
			if err != nil {
				return nil, err
			}
			task := Task{
				Target:  config.TestIntegration,
				Module:  moduleInfo,
				targets: targets,
			}
			task.Commands = []*Command{tb.createTestCommands(task)}
			tasks = append(tasks, &task)

		}
	}
	if arguments.Custom != nil && len(arguments.Custom) > 0 {
		for _, task := range arguments.Custom {
			task := Task{
				Target:  config.Custom,
				targets: task,
			}
			command, err := tb.createCustomCommands(task)
			if err != nil {
				return nil, err
			}
			task.Commands = []*Command{command}
			tasks = append(tasks, &task)
		}
	}
	if arguments.NumKey != nil && len(arguments.NumKey) > 0 {
		for _, task := range arguments.NumKey {
			moduleInfo, _, err := tb.getTaskSpec(task)
			if err != nil {
				return nil, err
			}
			task := Task{
				Target:  config.NumKey,
				targets: task,
				Module:  moduleInfo,
			}
			command := tb.createNumKeyCommand(task)
			task.Commands = []*Command{command}
			tasks = append(tasks, &task)
		}
	}
	if arguments.Restart {
		task := Task{
			Target:   config.Restart,
			Commands: []*Command{{Command: tb.appConfig.Commands.OOTB.Restart}},
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}

func (tb *taskBuilder) getTaskSpec(moduleSpec string) (*module.ModuleInfo, string, error) {
	spec := strings.Split(moduleSpec, "_")
	moduleId := spec[0]
	var targets string
	if len(spec) > 1 {
		targets = spec[1]
	}
	definedModule, err := tb.findModuleById(moduleId)
	if err != nil {
		return nil, "", err
	}
	return definedModule, targets, nil
}

func (tb *taskBuilder) findModuleById(id string) (*module.ModuleInfo, error) {
	moduleName := tb.appConfig.Aliases[id]
	if moduleName == "" {
		moduleByDirectCall := tb.modulesConfig[id]
		if moduleByDirectCall != nil {
			return moduleByDirectCall, nil
		} else {
			return nil, errors.New(fmt.Sprintf("Module alias %s not found.", id))
		}
	}
	return tb.modulesConfig[moduleName], nil
}

func (tb *taskBuilder) BuildTasks(arguments *config.ProgramArguments) ([]*Task, error) {
	return tb.buildExplicitTasks(arguments)
}

func (tb *taskBuilder) createBuildCommands(task Task) []*Command {
	commands := make([]*Command, 0, 5)
	for key, value := range config.ClobberableSources {
		if strings.Contains(task.targets, key) {
			if strings.Contains(task.targets, config.ClobberSymbol) {
				command := Command{
					Command: fmt.Sprintf(config.ClobberCommandFormat, task.Module.Location, value),
				}
				commands = append(commands, &command)
			}
			command := Command{
				Command: fmt.Sprintf(config.BuildCommandFormat, task.Module.Location, value),
			}
			commands = append(commands, &command)
		}
	}

	for key, value := range config.NonClobberableSources {
		if strings.Contains(task.targets, key) {
			command := Command{
				Command: fmt.Sprintf(config.BuildCommandFormat, task.Module.Location, value),
			}
			commands = append(commands, &command)
		}
	}
	return commands
}

func (tb *taskBuilder) createTestCommands(task Task) *Command {
	testCommand := fmt.Sprintf(config.TestCommandFormat, strings.ReplaceAll(task.Target.String(), "_", "."), task.Module.Location, config.SrcAliases[config.SrcTestSymbol])
	if task.targets != "" {
		testCommand = testCommand + fmt.Sprintf(config.SpecificTestCommandFormat, task.targets)
	}
	return &Command{Command: testCommand}
}

func (tb *taskBuilder) createNumKeyCommand(task Task) *Command {
	testCommand := fmt.Sprintf(config.NumKeyBuildCommandFormat, task.Module.Location, config.SrcAliases[config.SrcSymbol])
	return &Command{Command: testCommand}
}

func (tb *taskBuilder) createCustomCommands(task Task) (*Command, error) {
	command := tb.appConfig.Commands.Custom[task.targets]
	if command != "" {
		return &Command{Command: command}, nil
	}
	if tb.appConfig.FailOnError {
		return nil, fmt.Errorf("command %s not found in custom commands", task.targets)
	}
	return &Command{Command: task.targets, Status: config.Failed}, nil
}
