package cli

import (
	"tresca-vm/internal/yield"
)

type ExampleCase struct {
	Principal []float64 `json:"principal"`
	SigmaY    float64   `json:"sigma_y"`
}

func runExample(args []string) int {
	fs := flagSet("example")
	file := fs.String("file", "example/uniaxial.json", "scenario JSON file")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	caseValue, err := loadExample(*file)
	if err != nil {
		return fail(err)
	}
	result, err := yield.Evaluate(caseValue.Principal, caseValue.SigmaY)
	if err != nil {
		return fail(err)
	}
	return printJSON(result)
}

func loadExample(path string) (ExampleCase, error) {
	return loadExampleFile(path)
}
