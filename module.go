package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type ModuleInfo struct {
	Name     string
	Location string
	Order    int
	Sources  []string
}

type moduleXML struct {
	Name     string `xml:"name,attr"`
	Location string `xml:"location,attr"`
}

type modules struct {
	XMLName xml.Name    `xml:"ModuleRegistry"`
	Modules []moduleXML `xml:"Module"`
}

var orderPattern, _ = regexp.Compile("^#? ?(\\w+)/(\\w+)\n?")

func CalculateModuleInfo(cfg *AppConfig) (map[string]ModuleInfo, error) {
	orderCalculator, err := buildOrderCalculator(cfg)
	if err != nil {
		return nil, err
	}
	sourceCalculator := buildSourceCalculator(cfg)
	calculators := []func(info *ModuleInfo) error{orderCalculator, sourceCalculator}
	infos, err := buildModuleInfos(cfg, calculators)
	if err != nil {
		return nil, err
	}
	return infos, nil
}

func buildOrderCalculator(cfg *AppConfig) (func(info *ModuleInfo) error, error) {
	result, err := buildOrderMap(cfg)
	if err != nil {
		return nil, err
	}
	return func(info *ModuleInfo) error {
		info.Order = result[info.Name]
		return nil
	}, nil
}

func buildOrderMap(cfg *AppConfig) (map[string]int, error) {
	buildOrderPath := cfg.Input.BuildOrder
	idx := 0

	orderFile, err := os.Open(buildOrderPath)
	defer orderFile.Close()

	if err != nil {
		return nil, fmt.Errorf("the build order file is not available. Please check. %w", err)
	}
	scanner := bufio.NewScanner(orderFile)
	scanner.Split(bufio.ScanLines)

	result := make(map[string]int)
	for scanner.Scan() {
		line := scanner.Text()
		if orderPattern.MatchString(line) {
			result[strings.Split(strings.Trim(line, " \t\r\n#"), "/")[1]] = idx
			idx = idx + 1
		}
	}
	return result, nil
}

func buildSourceCalculator(cfg *AppConfig) func(info *ModuleInfo) error {
	return func(info *ModuleInfo) error {
		if cfg.Profile == "test" {
			info.Sources = TEST_SOURCES
			return nil
		}
		entries, err := os.ReadDir(info.Location)
		if err != nil {
			return fmt.Errorf("could not navigate through folder location: %s. %w", info.Location, err)
		}
		sources := make([]string, 0, 3)
		for _, entry := range entries {
			if entry.IsDir() && strings.Contains(entry.Name(), "src") {
				sources = append(sources, entry.Name())
			}
		}
		info.Sources = sources
		return nil
	}
}

func buildModuleInfos(cfg *AppConfig, calculators []func(info *ModuleInfo) error) (map[string]ModuleInfo, error) {
	moduleRegistryPath := cfg.Input.ModuleRegistry
	fileByteContent, _ := os.ReadFile(moduleRegistryPath)

	modulesFromXml := modules{}
	err := xml.Unmarshal(fileByteContent, &modulesFromXml)
	if err != nil {
		return nil, fmt.Errorf("module info could not be unmarshalled. %w", err)
	}

	result := make(map[string]ModuleInfo, len(modulesFromXml.Modules))
	for _, xmlModule := range modulesFromXml.Modules {
		name := strings.Split(xmlModule.Name, "/")[1]
		absLocation := strings.Join([]string{cfg.Root, xmlModule.Location}, "/")
		module := ModuleInfo{
			Name:     name,
			Location: absLocation,
			Order:    0,
			Sources:  nil,
		}
		for _, calculator := range calculators {
			err := calculator(&module)
			if err != nil {
				return nil, err
			}
		}
		result[name] = module
	}
	return result, nil
}
