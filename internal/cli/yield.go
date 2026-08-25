package cli

import (
	"fmt"
	"strconv"
	"strings"

	"tresca-vm/internal/yield"
)

func runYield(args []string) int {
	fs := flagSet("yield")
	principalText := fs.String("principal", "", "comma-separated principal stresses, e.g. 100,0,-50")
	sigmaY := fs.Float64("sigmaY", 0, "yield strength")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *principalText == "" {
		fmt.Fprintln(fs.Output(), "yield requires -principal")
		return 2
	}
	principal, err := parseList(*principalText)
	if err != nil {
		return fail(err)
	}
	result, err := yield.Evaluate(principal, *sigmaY)
	if err != nil {
		return fail(err)
	}
	return printJSON(result)
}

func runTensor(args []string) int {
	fs := flagSet("tensor")
	sxx := fs.Float64("sxx", 0, "normal stress x")
	syy := fs.Float64("syy", 0, "normal stress y")
	szz := fs.Float64("szz", 0, "normal stress z")
	sxy := fs.Float64("sxy", 0, "shear stress xy")
	sxz := fs.Float64("sxz", 0, "shear stress xz")
	syz := fs.Float64("syz", 0, "shear stress yz")
	sigmaY := fs.Float64("sigmaY", 0, "yield strength")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	t := tensorSymmetric(*sxx, *syy, *szz, *sxy, *sxz, *syz)
	result, err := yield.EvaluateTensor(t, *sigmaY)
	if err != nil {
		return fail(err)
	}
	return printJSON(result)
}

func parseList(text string) ([]float64, error) {
	parts := strings.Split(text, ",")
	out := make([]float64, 0, len(parts))
	for _, part := range parts {
		value, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number %q", part)
		}
		out = append(out, value)
	}
	if len(out) != 3 {
		return nil, fmt.Errorf("principal list must have 3 values, got %d", len(out))
	}
	return out, nil
}
