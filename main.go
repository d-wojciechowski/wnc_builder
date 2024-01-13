package main

import (
	"fmt"
	"wnc_builder/config"
	"wnc_builder/module"
)

func main() {

	appConfig, err := config.CreateAppConfig()
	if err != nil {
		fmt.Println(err.Error())
	}
	_ = config.ParseCmdArgs()

	_, err = module.CalculateModuleInfo(appConfig)
	if err != nil {
		fmt.Println(err.Error())
	}

}
