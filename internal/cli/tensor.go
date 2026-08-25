package cli

import (
	"tresca-vm/internal/tensor"
	"tresca-vm/internal/yield"
)

func tensorSymmetric(sxx, syy, szz, sxy, sxz, syz float64) tensor.Tensor {
	return tensor.NewSymmetric(sxx, syy, szz, sxy, sxz, syz)
}

func runPureShear(args []string) int {
	fs := flagSet("pure-shear")
	tau := fs.Float64("tau", 0, "shear stress")
	sigmaY := fs.Float64("sigmaY", 0, "yield strength")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	result, err := yield.PureShear(*tau, *sigmaY)
	if err != nil {
		return fail(err)
	}
	return printJSON(result)
}

func runPlane(args []string) int {
	fs := flagSet("plane")
	sxx := fs.Float64("sxx", 0, "normal stress x")
	syy := fs.Float64("syy", 0, "normal stress y")
	txy := fs.Float64("txy", 0, "shear stress xy")
	sigmaY := fs.Float64("sigmaY", 0, "yield strength")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	result, err := yield.PlaneStress(*sxx, *syy, *txy, *sigmaY)
	if err != nil {
		return fail(err)
	}
	return printJSON(result)
}
