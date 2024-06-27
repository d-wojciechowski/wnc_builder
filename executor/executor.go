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

type Executor interface {
	RunTasks(tasks []*Task) error
	RunCommands(tasks *Task) error
	PrintSummary(tasks []*Task)
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

func (e *executor) RunTasks(tasks []*Task) error {
	for _, task := range tasks {
		err := e.RunCommands(task)
		if e.appConfig.FailOnError && err != nil {
			return err
		}
	}
	return nil
}

func (e *executor) RunCommands(tasks *Task) error {
	for _, command := range tasks.Commands {
		err := e.runCommand(command)
		if e.appConfig.FailOnError && err != nil {
			return err
		}
	}
	return nil
}

func (e *executor) PrintSummary(tasks []*Task) {
	fmt.Println(strings.Repeat("-", config.CommandSize))
	fmt.Println("Application finished successfully")
	for _, task := range tasks {
		for _, command := range task.Commands {
			roundedDuration := roundDuration(command.Duration, time.Millisecond*10)
			fmt.Printf("%s %s %s in %s - %s\n", command.Status.Color(), command.Status, config.NoColor, roundedDuration, strings.Replace(command.Command, "\n", " \\n ", -1))
		}
	}
}

func roundDuration(d time.Duration, precision time.Duration) time.Duration {
	if precision <= 0 {
		return d
	}
	rounding := precision / 2
	return (d + rounding) / precision * precision
}

func (e *executor) runCommand(command *Command) error {
	e.printHeader(command)
	if command.Status != config.Prepared {
		return nil
	}
	command.Status = config.Running
	start := time.Now()

	toBeRun := e.prepareCommand(command)
	toBeRun.Stdout = os.Stdout
	toBeRun.Stderr = os.Stderr
	err := toBeRun.Run()

	command.Duration = time.Since(start)
	if err != nil {
		command.Status = config.Failed
		fmt.Printf("Command %s failed with code %s.\n", strings.Replace(command.Command, "\n", " \\n ", -1), toBeRun.Err)
		if e.appConfig.FailOnError {
			return err
		}
	} else {
		command.Status = config.Completed
		fmt.Printf("Command %s completed successfully.\n", strings.Replace(command.Command, "\n", "\\n", -1))
	}
	e.printFooter(command)
	return nil
}

func (e *executor) printHeader(command *Command) {
	fmt.Println(strings.Repeat(config.CmdFiller, config.CommandSize))
	message := "Executing command " + command.Command
	dashCombo := e.calculateFiller(len(message))
	fmt.Println(strings.Join([]string{dashCombo, message, dashCombo}, " "))
	fmt.Println(strings.Repeat(config.CmdFiller, config.CommandSize))
}

func (e *executor) calculateFiller(messageLen int) string {
	dashCount := ((config.CommandSize - messageLen) / 2) - 1
	if dashCount < 0 {
		dashCount = 0
	}
	return strings.Repeat(config.CmdFiller, dashCount)
}

func (e *executor) printFooter(command *Command) {}

func (e *executor) prepareCommand(cmd *Command) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/U", "/c", cmd.Command)
	} else {
		return exec.Command("sh", "-c", cmd.Command)
	}
}
