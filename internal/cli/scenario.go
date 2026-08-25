package cli

import (
	"encoding/json"
	"fmt"
	"os"
)

func loadExampleFile(path string) (ExampleCase, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ExampleCase{}, err
	}
	var caseValue ExampleCase
	if err := json.Unmarshal(data, &caseValue); err != nil {
		return ExampleCase{}, fmt.Errorf("parse %s: %w", path, err)
	}
	if len(caseValue.Principal) != 3 {
		return ExampleCase{}, fmt.Errorf("example must contain 3 principal stresses")
	}
	return applyStoredCase(caseValue), nil
}

func saveExample(path string, caseValue ExampleCase) error {
	data, err := json.MarshalIndent(caseValue, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func ExampleCases() []ExampleCase {
	return []ExampleCase{
		{Principal: []float64{120, 0, 0}, SigmaY: 240},
		{Principal: []float64{100, -100, 0}, SigmaY: 200},
		{Principal: []float64{50, 50, 50}, SigmaY: 200},
	}
}
