package executor

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"wnc_builder/config"
	"wnc_builder/module"
)

type ExecutionStatus int

const (
	PREPARED ExecutionStatus = iota
	RUNNING
	COMPLETED
	FAILED
)

func (t ExecutionStatus) String() string {
	return [...]string{"PREPARED", "RUNNING", "COMPLETED", "FAILED"}[t]
}
func (t ExecutionStatus) EnumIndex() int {
	return int(t)
}

type Command struct {
	Command  string
	Status   ExecutionStatus
	Duration time.Duration
}

type Task struct {
	Commands []*Command
	Target   config.Target
	Module   *module.ModuleInfo
	targets  string
}

type taskBuilder struct {
	appConfig     config.AppConfig
	modulesConfig map[string]*module.ModuleInfo
}

type TaskBuilder interface {
	BuildTasks(arguments *config.ProgramArguments) []*Task
}

func NewTaskBuilder(appConfig config.AppConfig, modulesConfig map[string]*module.ModuleInfo) TaskBuilder {
	builder := taskBuilder{
		appConfig:     appConfig,
		modulesConfig: modulesConfig,
	}
	return &builder
}

func (tb *taskBuilder) buildSuiteTasks(arguments *config.ProgramArguments) []*Task {
	tasks := make([]*Task, 0, 1)

	return tasks
}

func (tb *taskBuilder) buildExplicitTasks(arguments *config.ProgramArguments) ([]*Task, error) {
	tasks := make([]*Task, 0, 1)
	if len(arguments.Build) > 0 {
		for _, moduleSpec := range arguments.Build {
			moduleInfo, targets, err := tb.getTaskSpec(moduleSpec)
			if err != nil {
				return nil, err
			}
			task := Task{
				Target:  config.BUILD,
				Module:  moduleInfo,
				targets: targets,
			}
			task.Commands = tb.createBuildCommands(task)
			tasks = append(tasks, &task)
		}
	}
	if len(arguments.TestUnit) > 0 {
		for _, moduleSpec := range arguments.TestUnit {
			moduleInfo, targets, err := tb.getTaskSpec(moduleSpec)
			if err != nil {
				return nil, err
			}
			task := Task{
				Target:  config.TEST_UNIT,
				Module:  moduleInfo,
				targets: targets,
			}
			task.Commands = []*Command{tb.createTestCommands(task)}
			tasks = append(tasks, &task)

		}
	}
	if len(arguments.TestIntegration) > 0 {
		for _, moduleSpec := range arguments.TestUnit {
			moduleInfo, targets, err := tb.getTaskSpec(moduleSpec)
			if err != nil {
				return nil, err
			}
			task := Task{
				Target:  config.TEST_INTEGRATION,
				Module:  moduleInfo,
				targets: targets,
			}
			task.Commands = []*Command{tb.createTestCommands(task)}
			tasks = append(tasks, &task)

		}
	}
	if len(arguments.Custom) > 0 {
		for _, task := range arguments.Custom {
			task := Task{
				Target:  config.CUSTOM,
				targets: task,
			}
			task.Commands = []*Command{tb.createTestCommands(task)}
			tasks = append(tasks, &task)
		}
	}
	if arguments.Restart {
		task := Task{
			Target:   config.RESTART,
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
		return tb.buildSuiteTasks(arguments), nil
	} else {
		return tb.buildExplicitTasks(arguments)
	}
}

func (tb *taskBuilder) createBuildCommands(task Task) []*Command {
	commands := make([]*Command, 0, 5)
	if strings.Contains(task.targets, config.SRC_SYMBOL) {
		if strings.Contains(task.targets, config.CLOBBER_SYMBOL) {
			command := Command{
				Command: fmt.Sprintf(config.CLOBBER_COMMAND_FORMAT, task.Module.Location, config.SRC_ALIASES[config.SRC_SYMBOL]),
			}
			commands = append(commands, &command)
		}
		command := Command{
			Command: fmt.Sprintf(config.BUILD_COMMAND_FORMAT, task.Module.Location, config.SRC_ALIASES[config.SRC_SYMBOL]),
		}
		commands = append(commands, &command)
	}
	if strings.Contains(task.targets, config.SRC_TEST_SYMBOL) {
		if strings.Contains(task.targets, config.CLOBBER_SYMBOL) {
			command := Command{
				Command: fmt.Sprintf(config.CLOBBER_COMMAND_FORMAT, task.Module.Location, config.SRC_ALIASES[config.SRC_TEST_SYMBOL]),
			}
			commands = append(commands, &command)
		}
		command := Command{
			Command: fmt.Sprintf(config.BUILD_COMMAND_FORMAT, task.Module.Location, config.SRC_ALIASES[config.SRC_TEST_SYMBOL]),
		}
		commands = append(commands, &command)
	}
	if strings.Contains(task.targets, config.SRC_WEB_SYMBOL) {
		command := Command{
			Command: fmt.Sprintf(config.BUILD_COMMAND_FORMAT, task.Module.Location, config.SRC_ALIASES[config.SRC_WEB_SYMBOL]),
		}
		commands = append(commands, &command)
	}
	return commands
}

func (tb *taskBuilder) createTestCommands(task Task) *Command {
	testCommand := fmt.Sprintf(config.TEST_COMMAND_FORMAT, strings.ReplaceAll(task.Target.String(), "_", "."), task.Module.Location, config.SRC_ALIASES[config.SRC_TEST_SYMBOL])
	if task.targets != "" {
		testCommand = testCommand + fmt.Sprintf(config.SPECIFIC_TEST_COMMAND_FORMAT, task.targets)
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
	return &Command{Command: task.targets, Status: FAILED}, nil
}
