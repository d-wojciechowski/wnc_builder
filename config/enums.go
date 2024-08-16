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
	TestSelenium
	Restart
	Custom
	NumKey
)

func (t Target) String() string {
	return [...]string{"build", "test_unit", "test_integration", "test_selenium", "restart", "custom"}[t]
}
func (t Target) ModuleDependent() []string {
	return []string{"build", "test_unit", "test_integration", "test_selenium"}
}
func (t Target) ModuleAgnostic() []string {
	return []string{"build", "test_unit", "test_integration", "test_selenium"}
}
func (t Target) EnumIndex() int {
	return int(t)
}
