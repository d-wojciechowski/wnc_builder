package config

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
func (t Target) ModuleDependent() []string {
	return []string{"build", "test_unit", "test_integration", "suite"}
}
func (t Target) ModuleAgnostic() []string {
	return []string{"build", "test_unit", "test_integration", "suite"}
}
func (t Target) EnumIndex() int {
	return int(t)
}

const COMMAND_SIZE = 128

const SRC = "src"
const SRC_TEST = "src_test"
const SRC_WEB = "src_web"

const SRC_SYMBOL = "s"
const SRC_WEB_SYMBOL = "w"
const SRC_TEST_SYMBOL = "t"
const CLOBBER_SYMBOL = "c"

var SRC_ALIASES = map[string]string{"s": SRC, "t": SRC_TEST, "w": SRC_WEB}
var TEST_SOURCES = []string{SRC, SRC_TEST, SRC_WEB}

const BUILD_COMMAND_FORMAT = "ant -f %s/%s/build.xml"
const CLOBBER_COMMAND_FORMAT = "ant clobber -f %s/%s/build.xml"
const TEST_COMMAND_FORMAT = "ant %s -f %s/%s/build.xml"
const SPECIFIC_TEST_COMMAND_FORMAT = " -Dtest.includes=**/%s"
