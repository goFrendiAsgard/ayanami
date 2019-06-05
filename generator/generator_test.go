package generator

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type testConfig struct {
	Name          string
	Valid         bool
	State         map[string]bool
	ScaffoldError bool
	BuildError    bool
}

func (config testConfig) Validate() bool {
	return config.Valid
}

func (config testConfig) Scaffold() error {
	if config.ScaffoldError {
		config.State[fmt.Sprintf("config%sScaffoldError", config.Name)] = true
		return errors.New("error")
	}
	config.State[fmt.Sprintf("config%sScaffoldDone", config.Name)] = true
	return nil
}

func (config testConfig) Build() error {
	if config.BuildError {
		config.State[fmt.Sprintf("config%sBuildError", config.Name)] = true
		return errors.New("error")
	}
	config.State[fmt.Sprintf("config%sBuildDone", config.Name)] = true
	return nil
}

type testProcedure struct {
	Name          string
	Valid         bool
	State         map[string]bool
	ScaffoldError bool
	BuildError    bool
}

func (procedure testProcedure) Validate(config Configs) bool {
	return procedure.Valid
}

func (procedure testProcedure) Scaffold(config Configs) error {
	if procedure.ScaffoldError {
		procedure.State[fmt.Sprintf("procedure%sScaffoldError", procedure.Name)] = true
		return errors.New("error")
	}
	procedure.State[fmt.Sprintf("procedure%sScaffoldDone", procedure.Name)] = true
	return nil
}

func (procedure testProcedure) Build(config Configs) error {
	if procedure.BuildError {
		procedure.State[fmt.Sprintf("procedure%sBuildError", procedure.Name)] = true
		return errors.New("error")
	}
	procedure.State[fmt.Sprintf("procedure%sBuildDone", procedure.Name)] = true
	return nil
}

func TestNormal(t *testing.T) {
	actualState := map[string]bool{}
	// set generator
	generator := NewGenerator()
	generator.AddConfigs(
		Configs{
			testConfig{Name: "A", Valid: true, State: actualState},
			testConfig{Name: "B", Valid: true, State: actualState},
		},
	)
	generator.AddProcedures(
		Procedures{
			testProcedure{Name: "A", Valid: true, State: actualState},
			testProcedure{Name: "B", Valid: true, State: actualState},
		},
	)
	// scaffold
	err := generator.Scaffold()
	if err != nil {
		t.Errorf("Getting error: %s", err)
	}
	// build
	err = generator.Build()
	if err != nil {
		t.Errorf("Getting error: %s", err)
	}
	// evaluate
	expectedState := map[string]bool{
		"configAScaffoldDone":    true,
		"configBScaffoldDone":    true,
		"procedureAScaffoldDone": true,
		"procedureBScaffoldDone": true,
		"configABuildDone":       true,
		"configBBuildDone":       true,
		"procedureABuildDone":    true,
		"procedureBBuildDone":    true,
	}
	if !reflect.DeepEqual(expectedState, actualState) {
		t.Errorf("expected : %#v\nget: %#v", expectedState, actualState)
	}
}

func TestInvalidConfig(t *testing.T) {
	actualState := map[string]bool{}
	// set generator
	generator := NewGenerator()
	generator.AddConfigs(
		Configs{
			testConfig{Name: "A", Valid: true, State: actualState},
			testConfig{Name: "B", Valid: false, State: actualState},
		},
	)
	generator.AddProcedures(
		Procedures{
			testProcedure{Name: "A", Valid: true, State: actualState},
			testProcedure{Name: "B", Valid: true, State: actualState},
		},
	)
	// scaffold
	err := generator.Scaffold()
	if err == nil {
		t.Error("Error expected")
	}
	// build
	err = generator.Build()
	if err == nil {
		t.Error("Error expected")
	}
	// evaluate
	expectedState := map[string]bool{}
	if !reflect.DeepEqual(expectedState, actualState) {
		t.Errorf("expected : %#v\nget: %#v", expectedState, actualState)
	}
}

func TestInvalidProcedure(t *testing.T) {
	actualState := map[string]bool{}
	// set generator
	generator := NewGenerator()
	generator.AddConfigs(
		Configs{
			testConfig{Name: "A", Valid: true, State: actualState},
			testConfig{Name: "B", Valid: true, State: actualState},
		},
	)
	generator.AddProcedures(
		Procedures{
			testProcedure{Name: "A", Valid: true, State: actualState},
			testProcedure{Name: "B", Valid: false, State: actualState},
		},
	)
	// scaffold
	err := generator.Scaffold()
	if err == nil {
		t.Error("Error expected")
	}
	// build
	err = generator.Build()
	if err == nil {
		t.Error("Error expected")
	}
	// evaluate
	expectedState := map[string]bool{}
	if !reflect.DeepEqual(expectedState, actualState) {
		t.Errorf("expected : %#v\nget: %#v", expectedState, actualState)
	}
}

func TestErrorConfig(t *testing.T) {
	actualState := map[string]bool{}
	// set generator
	generator := NewGenerator()
	generator.AddConfigs(
		Configs{
			testConfig{Name: "A", Valid: true, State: actualState},
			testConfig{Name: "B", Valid: true, State: actualState, ScaffoldError: true, BuildError: true},
		},
	)
	generator.AddProcedures(
		Procedures{
			testProcedure{Name: "A", Valid: true, State: actualState},
			testProcedure{Name: "B", Valid: true, State: actualState},
		},
	)
	// scaffold
	err := generator.Scaffold()
	if err == nil {
		t.Error("Error expected")
	}
	// build
	err = generator.Build()
	if err == nil {
		t.Error("Error expected")
	}
	// evaluate
	expectedState := map[string]bool{
		"configAScaffoldDone":  true,
		"configABuildDone":     true,
		"configBScaffoldError": true,
		"configBBuildError":    true,
	}
	if !reflect.DeepEqual(expectedState, actualState) {
		t.Errorf("expected : %#v\nget: %#v", expectedState, actualState)
	}
}

func TestErrorProcedure(t *testing.T) {
	actualState := map[string]bool{}
	// set generator
	generator := NewGenerator()
	generator.AddConfigs(
		Configs{
			testConfig{Name: "A", Valid: true, State: actualState},
			testConfig{Name: "B", Valid: true, State: actualState},
		},
	)
	generator.AddProcedures(
		Procedures{
			testProcedure{Name: "A", Valid: true, State: actualState},
			testProcedure{Name: "B", Valid: true, State: actualState, ScaffoldError: true, BuildError: true},
		},
	)
	// scaffold
	err := generator.Scaffold()
	if err == nil {
		t.Error("Error expected")
	}
	// build
	err = generator.Build()
	if err == nil {
		t.Error("Error expected")
	}
	// evaluate
	expectedState := map[string]bool{
		"configAScaffoldDone":     true,
		"configABuildDone":        true,
		"configBScaffoldDone":     true,
		"configBBuildDone":        true,
		"procedureAScaffoldDone":  true,
		"procedureABuildDone":     true,
		"procedureBScaffoldError": true,
		"procedureBBuildError":    true,
	}
	if !reflect.DeepEqual(expectedState, actualState) {
		t.Errorf("expected : %#v\nget: %#v", expectedState, actualState)
	}
}
