package executor

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"
	"wnc_builder/config"
	"wnc_builder/module"
)

func Test_executor_calculateFiller(t *testing.T) {
	type fields struct {
		appConfig     *config.AppConfig
		modulesConfig map[string]*module.ModuleInfo
	}
	type args struct {
		messageLen int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "Calculate with even msg len",
			fields: fields{},
			args:   args{config.CommandSize - 8},
			want:   "---",
		},
		{
			name:   "Calculate with uneven msg len",
			fields: fields{},
			args:   args{config.CommandSize - 7},
			want:   "--",
		},
		{
			name:   "Calculate with message equal with command size",
			fields: fields{},
			args:   args{config.CommandSize},
			want:   "",
		},
		{
			name:   "Calculate with message less than command size",
			fields: fields{},
			args:   args{config.CommandSize + 1},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{
				appConfig:     tt.fields.appConfig,
				modulesConfig: tt.fields.modulesConfig,
			}
			if got := e.calculateFiller(tt.args.messageLen); got != tt.want {
				t.Errorf("calculateFiller() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_executor_printHeader(t *testing.T) {
	type fields struct {
		appConfig     *config.AppConfig
		modulesConfig map[string]*module.ModuleInfo
	}
	type args struct {
		command *Command
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "A test with example input",
			fields: fields{},
			args: args{
				&Command{
					Command:  "echo \"Test\"",
					Status:   config.Completed,
					Duration: time.Second,
				},
			},
			want: wrapExpectedMessage(" Executing command echo \"Test\" "),
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			print()

			outC := make(chan string)
			// copy the output in a separate goroutine so printing can't block indefinitely
			go func() {
				var buf bytes.Buffer
				io.Copy(&buf, r)
				outC <- buf.String()
			}()

			e := &executor{
				appConfig:     tt.fields.appConfig,
				modulesConfig: tt.fields.modulesConfig,
			}
			e.printHeader(tt.args.command)
			w.Close()
			os.Stdout = old // restoring the real stdout
			got := <-outC
			if got != tt.want {
				t.Errorf("printHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func wrapExpectedMessage(message string) string {
	return strings.Repeat(config.CmdFiller, config.CommandSize) + "\n" + strings.Repeat(config.CmdFiller, 48) + message + strings.Repeat(config.CmdFiller, 48) + "\n" + strings.Repeat(config.CmdFiller, config.CommandSize) + "\n"
}

func Test_executor_runCommand(t *testing.T) {
	type fields struct {
		appConfig     *config.AppConfig
		modulesConfig map[string]*module.ModuleInfo
	}
	type args struct {
		command *Command
	}
	type want struct {
		header  string
		suffix  string
		command *Command
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    want
		wantErr bool
	}{
		{
			name:   "Should run green command",
			fields: fields{},
			args: args{
				command: &Command{
					Command: "echo \"Test\"",
					Status:  config.Prepared,
				},
			},
			want: want{
				header: " Executing command echo \"Test\"",
				suffix: "Command echo \"Test\" completed successfully.\n",
				command: &Command{
					Command: "echo \"Test\"",
					Status:  config.Completed,
				},
			},
			wantErr: false,
		},
		{
			name:   "Should run error command",
			fields: fields{appConfig: &config.AppConfig{FailOnError: true}},
			args: args{
				command: &Command{
					Command: "ant build",
					Status:  config.Prepared,
				},
			},
			want: want{
				header: " Executing command ant build ",
				suffix: "Command ant build failed with code %!s(<nil>).\n",
				command: &Command{
					Command: "echo \"Test\"",
					Status:  config.Failed,
				},
			},
			wantErr: true,
		},
		{
			name:   "Should run error command",
			fields: fields{appConfig: &config.AppConfig{FailOnError: false}},
			args: args{
				command: &Command{
					Command: "ant build",
					Status:  config.Prepared,
				},
			},
			want: want{
				header: " Executing command ant build ",
				suffix: "Command ant build failed with code %!s(<nil>).\n",
				command: &Command{
					Command: "echo \"Test\"",
					Status:  config.Failed,
				},
			},
			wantErr: false,
		},
		{
			name:   "Should not run command when status not Prepared",
			fields: fields{appConfig: &config.AppConfig{FailOnError: false}},
			args: args{
				command: &Command{
					Command: "ant build",
					Status:  config.Failed,
				},
			},
			want: want{
				header: "",
				suffix: "",
				command: &Command{
					Command: "echo \"Test\"",
					Status:  config.Failed,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			print()

			outC := make(chan string)
			// copy the output in a separate goroutine so printing can't block indefinitely
			go func() {
				var buf bytes.Buffer
				io.Copy(&buf, r)
				outC <- buf.String()
			}()

			e := &executor{
				appConfig:     tt.fields.appConfig,
				modulesConfig: tt.fields.modulesConfig,
			}
			err := e.runCommand(tt.args.command)
			w.Close()
			os.Stdout = old // restoring the real stdout
			got := <-outC

			if tt.args.command.Status != tt.want.command.Status {
				t.Errorf("runCommand() error = status does not match got %v != want %v", tt.args.command.Status, tt.want.command.Status)
			}

			if err != nil {
				if tt.wantErr {
					if !strings.Contains(got, tt.want.header) || !strings.HasSuffix(got, tt.want.suffix) {
						t.Errorf("runCommand() error = %v, wantErr %v", err, tt.wantErr)
					}
					return
				} else {
					t.Errorf("runCommand() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if !strings.Contains(got, tt.want.header) || !strings.HasSuffix(got, tt.want.suffix) {
				t.Errorf("printHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_roundDuration(t *testing.T) {
	type args struct {
		d         time.Duration
		precision time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "Should limit duration to 2 significant digits",
			args: args{
				d:         time.Microsecond * 1234567,
				precision: time.Millisecond * 10,
			},
			want: time.Millisecond * 1230,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := roundDuration(tt.args.d, tt.args.precision); got != tt.want {
				t.Errorf("roundDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
