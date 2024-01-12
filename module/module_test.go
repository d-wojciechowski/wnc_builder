package module

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"wnc_builder/config"
)

func buildModulePath(part string) string {
	return "/" + strings.Join([]string{"opt", "a", "path", "to", part}, "/")
}

func Test_buildModuleInfos(t *testing.T) {
	type args struct {
		cfg         *config.AppConfig
		calculators []func(info *ModuleInfo) error
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*ModuleInfo
		wantErr bool
	}{
		{
			name: "Should parse order correctly with fixture",
			args: args{cfg: &config.AppConfig{
				Root: "/opt",
				Input: config.Input{
					ModuleRegistry: "../testFixtures/moduleRegistry.xml",
				},
			}, calculators: []func(info *ModuleInfo) error{},
			},
			want: map[string]*ModuleInfo{"ModuleA": {
				Name:     "ModuleA",
				Location: buildModulePath("ModuleA"),
				Order:    0,
				Sources:  nil,
			}, "ModuleB": {
				Name:     "ModuleB",
				Location: buildModulePath("ModuleB"),
				Order:    0,
				Sources:  nil,
			}, "ModuleC": {
				Name:     "ModuleC",
				Location: buildModulePath("ModuleC"),
				Order:    0,
				Sources:  nil,
			}},
			wantErr: false,
		},
		{
			name: "Should parse apply calculator when available",
			args: args{cfg: &config.AppConfig{
				Root: "/opt",
				Input: config.Input{
					ModuleRegistry: "../testFixtures/moduleRegistry.xml",
				},
			}, calculators: []func(info *ModuleInfo) error{func(info *ModuleInfo) error {
				info.Order = 1
				return nil
			}},
			},
			want: map[string]*ModuleInfo{"ModuleA": {
				Name:     "ModuleA",
				Location: buildModulePath("ModuleA"),
				Order:    1,
				Sources:  nil,
			}, "ModuleB": {
				Name:     "ModuleB",
				Location: buildModulePath("ModuleB"),
				Order:    1,
				Sources:  nil,
			}, "ModuleC": {
				Name:     "ModuleC",
				Location: buildModulePath("ModuleC"),
				Order:    1,
				Sources:  nil,
			}},
			wantErr: false,
		},
		{
			name: "Should not parse when incorrect file specified",
			args: args{cfg: &config.AppConfig{
				Input: config.Input{
					BuildOrder: "../testFixtures/notExisting.xml",
				},
			}, calculators: []func(info *ModuleInfo) error{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildModuleInfos(tt.args.cfg, tt.args.calculators)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildModuleInfos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildModuleInfos() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBuildOrderMap(t *testing.T) {
	type args struct {
		cfg *config.AppConfig
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]int
		wantErr bool
	}{
		{
			name: "Should parse order correctly with fixture",
			args: args{&config.AppConfig{
				Input: config.Input{
					BuildOrder: "../testFixtures/orderFile.includes",
				},
			}},
			want:    map[string]int{"ModuleA": 0, "ModuleB": 1},
			wantErr: false,
		},
		{
			name: "Should not parse when incorrect file specified",
			args: args{&config.AppConfig{
				Input: config.Input{
					BuildOrder: "../testFixtures/notExisting.includes",
				},
			}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildOrderMap(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildOrderCalculator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildOrderCalculator() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildSourceCalculator(t *testing.T) {
	type args struct {
		cfg *config.AppConfig
	}
	type preparation struct {
		root     string
		children []string
	}
	tests := []struct {
		name    string
		prep    preparation
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Should create expected sources for test profile",
			args: args{
				cfg: &config.AppConfig{Profile: "test"},
			},
			prep: preparation{
				root:     "",
				children: nil,
			},
			want:    config.TestSources,
			wantErr: false,
		},
		{
			name: "Should create a list of sources",
			args: args{
				cfg: &config.AppConfig{Profile: "prod"},
			},
			prep: preparation{
				root:     ".",
				children: []string{"src", "NOT"},
			},
			want:    []string{"src"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toRemove := make([]string, 0, len(tt.prep.children)+1)
			rootPath := tt.prep.root + string(os.PathSeparator) + "root"
			toRemove = append(toRemove, rootPath)
			err := os.Mkdir(rootPath, os.ModePerm)
			for _, child := range tt.prep.children {
				path := rootPath + string(os.PathSeparator) + child
				toRemove = append(toRemove, path)
				err = os.Mkdir(path, os.ModePerm)
			}

			calculator := buildSourceCalculator(tt.args.cfg)
			moduleInfo := ModuleInfo{Location: rootPath}
			err = calculator(&moduleInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildSourceCalculator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(moduleInfo.Sources, tt.want) {
				t.Errorf("buildSourceCalculator() = %v, want %v", moduleInfo.Sources, tt.want)
			}
			for _, path := range toRemove {
				err := os.Remove(path)
				if err != nil {
					return
				}
			}
		})
	}
}
