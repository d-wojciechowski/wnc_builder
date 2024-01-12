package main

import (
	"fmt"
	"os"
	"wnc_builder/config"
	"wnc_builder/executor"
	"wnc_builder/module"
)

func main() {
	appConfig, err := config.CreateAppConfig()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	cmdArgs := config.ParseCmdArgs()

	moduleInfos, err := module.CalculateModuleInfo(appConfig)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	taskBuilder := executor.NewTaskBuilder(appConfig, moduleInfos)
	tasks, err := taskBuilder.BuildTasks(cmdArgs)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	taskExecutor := executor.NewTaskExecutor(appConfig, moduleInfos)
	err = taskExecutor.RunTasks(tasks)
	taskExecutor.PrintSummary(tasks)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
