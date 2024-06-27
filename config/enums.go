package config

type ExecutionStatus int

const (
	Prepared ExecutionStatus = iota
	Running
	Completed
	Failed
)

func (t ExecutionStatus) String() string {
	return [...]string{"PREPARED", "RUNNING", "COMPLETED", "FAILED"}[t]
}
func (t ExecutionStatus) Color() string {
	return [...]string{WarningColor, WarningColor, OkColor, ErrColor}[t]
}
func (t ExecutionStatus) EnumIndex() int {
	return int(t)
}

type Target int

const (
	Build Target = iota
	TestUnit
	TestIntegration
	SuiteTarget
	Restart
	Custom
	NumKey
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
