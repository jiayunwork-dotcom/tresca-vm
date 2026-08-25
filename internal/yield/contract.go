package yield

import (
	"fmt"
	"math"

	"tresca-vm/internal/tensor"
)

type ContractCheck struct {
	Name    string `json:"name"`
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func VerifyUniaxial(sigmaY float64) ContractCheck {
	result, err := Evaluate([]float64{sigmaY, 0, 0}, sigmaY)
	if err != nil {
		return ContractCheck{Name: "uniaxial", OK: false, Message: err.Error()}
	}
	ok := result.TrescaYielded && result.MisesYielded &&
		math.Abs(result.TrescaStress-sigmaY) < 1e-8 &&
		math.Abs(result.MisesStress-sigmaY) < 1e-8
	return ContractCheck{
		Name:    "uniaxial",
		OK:      ok,
		Message: fmt.Sprintf("tresca=%.6g mises=%.6g sigma_y=%.6g", result.TrescaStress, result.MisesStress, sigmaY),
	}
}

func VerifyPureShearBoundary(sigmaY float64) ContractCheck {
	result, err := PureShear(sigmaY/2, sigmaY)
	if err != nil {
		return ContractCheck{Name: "pure-shear", OK: false, Message: err.Error()}
	}
	ok := result.TrescaYielded && !result.MisesYielded
	return ContractCheck{
		Name:    "pure-shear",
		OK:      ok,
		Message: fmt.Sprintf("tau=%.6g tresca=%.6g mises=%.6g", sigmaY/2, result.TrescaStress, result.MisesStress),
	}
}

func VerifyHydrostatic(mean float64) ContractCheck {
	result, err := Evaluate([]float64{mean, mean, mean}, 200)
	if err != nil {
		return ContractCheck{Name: "hydrostatic", OK: false, Message: err.Error()}
	}
	ok := result.TrescaStress == 0 && result.MisesStress == 0 && !result.AnyYielded()
	return ContractCheck{
		Name:    "hydrostatic",
		OK:      ok,
		Message: fmt.Sprintf("tresca=%.6g mises=%.6g", result.TrescaStress, result.MisesStress),
	}
}

func VerifyScaleRule(principal []float64, factor float64) ContractCheck {
	base, err := Evaluate(principal, 300)
	if err != nil {
		return ContractCheck{Name: "scale", OK: false, Message: err.Error()}
	}
	scaledPrincipal := make([]float64, 3)
	for i, value := range principal {
		scaledPrincipal[i] = value * factor
	}
	scaled, err := Evaluate(scaledPrincipal, 300)
	if err != nil {
		return ContractCheck{Name: "scale", OK: false, Message: err.Error()}
	}
	ok := math.Abs(scaled.TrescaStress-base.TrescaStress*factor) < 1e-8 &&
		math.Abs(scaled.MisesStress-base.MisesStress*factor) < 1e-8
	return ContractCheck{
		Name:    "scale",
		OK:      ok,
		Message: fmt.Sprintf("factor=%.3g base=%.6g scaled=%.6g", factor, base.TrescaStress, scaled.TrescaStress),
	}
}

func VerifySignReversal(principal []float64) ContractCheck {
	base, err := Evaluate(principal, 300)
	if err != nil {
		return ContractCheck{Name: "sign", OK: false, Message: err.Error()}
	}
	negated := make([]float64, 3)
	for i, value := range principal {
		negated[i] = -value
	}
	reversed, err := Evaluate(negated, 300)
	if err != nil {
		return ContractCheck{Name: "sign", OK: false, Message: err.Error()}
	}
	ok := math.Abs(base.TrescaStress-reversed.TrescaStress) < 1e-8 &&
		math.Abs(base.MisesStress-reversed.MisesStress) < 1e-8
	return ContractCheck{
		Name:    "sign",
		OK:      ok,
		Message: fmt.Sprintf("tresca=%.6g/%.6g mises=%.6g/%.6g", base.TrescaStress, reversed.TrescaStress, base.MisesStress, reversed.MisesStress),
	}
}

func VerifyPlaneZeroThird(sxx, syy, txy float64) ContractCheck {
	result, err := PlaneStress(sxx, syy, txy, 300)
	if err != nil {
		return ContractCheck{Name: "plane-zero", OK: false, Message: err.Error()}
	}
	hasZero := false
	for _, value := range result.Principal {
		if math.Abs(value) < 1e-9 {
			hasZero = true
		}
	}
	return ContractCheck{
		Name:    "plane-zero",
		OK:      hasZero,
		Message: fmt.Sprintf("principal=%v", result.Principal),
	}
}

func RunContractChecks() []ContractCheck {
	checks := []ContractCheck{
		VerifyUniaxial(200),
		VerifyPureShearBoundary(200),
		VerifyHydrostatic(80),
		VerifyScaleRule([]float64{100, 40, -20}, 2),
		VerifySignReversal([]float64{100, 50, -30}),
		VerifyPlaneZeroThird(120, -40, 30),
	}
	return checks
}

func ContractsPass(checks []ContractCheck) bool {
	for _, check := range checks {
		if !check.OK {
			return false
		}
	}
	return true
}

func EquivalentFromTensorSafe(t tensor.Tensor, sigmaY float64) Result {
	result, err := EvaluateTensor(t, sigmaY)
	if err != nil {
		return Result{}
	}
	return result
}
