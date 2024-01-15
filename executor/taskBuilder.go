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

func (tb *taskBuilder) buildSuiteTasks(arguments *config.ProgramArguments) ([]*Task, error) {
	tasks := make([]*Task, 0, 1)
	if arguments.Suite != nil && len(arguments.Suite) > 0 && arguments.Suite[0] != "" {
		suite := tb.appConfig.Suites[arguments.Suite[0]]
		for moduleSpec, targets := range suite.Build {
			moduleInfo, err := tb.findModuleById(moduleSpec)
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
		if suite.Custom != nil {
			for _, task := range suite.Custom {
				task := Task{
					Target:  config.Custom,
					targets: task,
				}
				task.Commands = []*Command{tb.createTestCommands(task)}
				tasks = append(tasks, &task)
			}
		}
		if suite.Restart {
			task := Task{
				Target:   config.Restart,
				Commands: []*Command{{Command: tb.appConfig.Commands.OOTB.Restart}},
			}
			tasks = append(tasks, &task)
		}
		return tasks, nil
	}
	return tasks, fmt.Errorf("could not find suite with name %s", arguments.Suite[0])
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
		for _, moduleSpec := range arguments.TestUnit {
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
	targets := spec[1]
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
	if len(arguments.Suite) > 0 {
		return tb.buildSuiteTasks(arguments)
	} else {
		return tb.buildExplicitTasks(arguments)
	}
}

func (tb *taskBuilder) createBuildCommands(task Task) []*Command {
	commands := make([]*Command, 0, 5)
	if strings.Contains(task.targets, config.SrcSymbol) {
		if strings.Contains(task.targets, config.ClobberSymbol) {
			command := Command{
				Command: fmt.Sprintf(config.ClobberCommandFormat, task.Module.Location, config.SrcAliases[config.SrcSymbol]),
			}
			commands = append(commands, &command)
		}
		command := Command{
			Command: fmt.Sprintf(config.BuildCommandFormat, task.Module.Location, config.SrcAliases[config.SrcSymbol]),
		}
		commands = append(commands, &command)
	}
	if strings.Contains(task.targets, config.SrcTestSymbol) {
		if strings.Contains(task.targets, config.ClobberSymbol) {
			command := Command{
				Command: fmt.Sprintf(config.ClobberCommandFormat, task.Module.Location, config.SrcAliases[config.SrcTestSymbol]),
			}
			commands = append(commands, &command)
		}
		command := Command{
			Command: fmt.Sprintf(config.BuildCommandFormat, task.Module.Location, config.SrcAliases[config.SrcTestSymbol]),
		}
		commands = append(commands, &command)
	}
	if strings.Contains(task.targets, config.SrcWebSymbol) {
		command := Command{
			Command: fmt.Sprintf(config.BuildCommandFormat, task.Module.Location, config.SrcAliases[config.SrcWebSymbol]),
		}
		commands = append(commands, &command)
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
