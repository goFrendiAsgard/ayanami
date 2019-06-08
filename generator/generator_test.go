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

func (c testConfig) Validate() bool {
	return c.Valid
}

func (c testConfig) Scaffold() error {
	if c.ScaffoldError {
		c.State[fmt.Sprintf("config%sScaffoldError", c.Name)] = true
		return errors.New("error")
	}
	c.State[fmt.Sprintf("config%sScaffoldDone", c.Name)] = true
	return nil
}

func (c testConfig) Build() error {
	if c.BuildError {
		c.State[fmt.Sprintf("config%sBuildError", c.Name)] = true
		return errors.New("error")
	}
	c.State[fmt.Sprintf("config%sBuildDone", c.Name)] = true
	return nil
}

type testProcedure struct {
	Name          string
	Valid         bool
	State         map[string]bool
	ScaffoldError bool
	BuildError    bool
}

func (p testProcedure) Validate(config Configs) bool {
	return p.Valid
}

func (p testProcedure) Scaffold(config Configs) error {
	if p.ScaffoldError {
		p.State[fmt.Sprintf("procedure%sScaffoldError", p.Name)] = true
		return errors.New("error")
	}
	p.State[fmt.Sprintf("procedure%sScaffoldDone", p.Name)] = true
	return nil
}

func (p testProcedure) Build(config Configs) error {
	if p.BuildError {
		p.State[fmt.Sprintf("procedure%sBuildError", p.Name)] = true
		return errors.New("error")
	}
	p.State[fmt.Sprintf("procedure%sBuildDone", p.Name)] = true
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
