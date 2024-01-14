package executor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
	"wnc_builder/config"
	"wnc_builder/module"
)

type Command struct {
	Command  string
	Status   config.ExecutionStatus
	Duration time.Duration
}

type Task struct {
	Commands []*Command
	Target   config.Target
	Module   *module.ModuleInfo
	targets  string
}

type executor struct {
	appConfig     *config.AppConfig
	modulesConfig map[string]*module.ModuleInfo
}

func NewTaskExecutor(appConfig *config.AppConfig, modulesConfig map[string]*module.ModuleInfo) Executor {
	executor := executor{
		appConfig:     appConfig,
		modulesConfig: modulesConfig,
	}
	return &executor
}

type Executor interface {
	RunTasks(tasks []*Task) error
	RunCommands(tasks *Task) error
	PrintSummary(tasks []*Task)
}

func (e *executor) RunTasks(tasks []*Task) error {
	for _, task := range tasks {
		return e.RunCommands(task)
	}
	return nil
}

func (e *executor) RunCommands(tasks *Task) error {
	for _, command := range tasks.Commands {
		return e.runCommand(command)
	}
	return nil
}

func (e *executor) PrintSummary(tasks []*Task) {
	fmt.Println(strings.Repeat("-", config.CommandSize))
	fmt.Println("Application finished successfully")
	for _, task := range tasks {
		for _, command := range task.Commands {
			fmt.Printf("%s %s %s in %s - %s\n", command.Status.Color(), command.Status, config.NoColor, command.Duration, command.Command)
		}
	}
}

func (e *executor) runCommand(command *Command) error {
	e.printHeader(command)
	if command.Status != config.Prepared {
		return nil
	}
	command.Status = config.Running
	start := time.Now()
	toBeRun := e.preparecommand(command)
	toBeRun.Stdout = os.Stdout
	toBeRun.Stderr = os.Stderr
	err := toBeRun.Run()
	command.Duration = time.Since(start)
	if err != nil {
		command.Status = config.Failed
		fmt.Printf("Command %s failed with code %s.\n", command.Command, toBeRun.Err)
		if e.appConfig.FailOnError {
			return err
		}
	} else {
		command.Status = config.Completed
		fmt.Printf("Command %s completed successfully.\n", command.Command)
	}
	e.printFooter(command)
	return nil
}

func (e *executor) printHeader(command *Command) {
	println(strings.Repeat("-", config.CommandSize))
	message := "Executing command " + command.Command
	dashCount := ((config.CommandSize - len(message)) / 2) - 1
	println(strings.Repeat("-", dashCount) + " " + message + " " + strings.Repeat("-", dashCount))
	println(strings.Repeat("-", config.CommandSize))
}

func (e *executor) printFooter(command *Command) {}

func (e *executor) preparecommand(cmd *Command) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/U", "/c", cmd.Command)
	} else {
		return exec.Command("sh", "-c", cmd.Command)
	}
}
