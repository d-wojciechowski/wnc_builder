package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
)

type ProgramArguments struct {
	Build           []string `arg:"-b,--build" help:"Execute build "`
	TestUnit        []string `arg:"-u,--test-unit" help:"Execute [unit tests] / [unit test by name]"`
	TestIntegration []string `arg:"-i,--test-integration" help:"Execute [integ tests] / [integ test by name]"`
	Suite           []string `arg:"-s,--suite" help:"Execute suite defined in CFG"`
	Custom          []string `arg:"-c,--custom" help:"Execute custom command defined in CFG"`
	Restart         bool     `arg:"-r,--restart" help:"Execute restart"`
}

func main() {

	appConfig, err := createAppConfig()
	if err != nil {
		fmt.Println(err.Error())
	}
	_ = parseCmdArgs()
	println()

	buildModuleInfos(appConfig, nil)

}

func parseCmdArgs() *ProgramArguments {
	args := &ProgramArguments{}
	arg.MustParse(args)
	return args
}
