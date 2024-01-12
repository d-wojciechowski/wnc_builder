package main

type Target int

const (
	BUILD Target = iota
	TEST_UNIT
	TEST_INTEGRATION
	SUITE
	RESTART
	CUSTOM
)

func (t Target) String() string {
	return [...]string{"build", "test_unit", "test_integration", "suite", "restart", "custom"}[t]
}
func (t Target) EnumIndex() int {
	return int(t)
}

const COMMAND_SIZE = 128

const SRC = "src"
const SRC_TEST = "src_test"
const SRC_WEB = "src_web"

var SRC_ALIASES = map[string]string{"s": SRC, "t": SRC_TEST, "w": SRC_WEB}
var TEST_SOURCES = []string{SRC, SRC_TEST, SRC_WEB}
