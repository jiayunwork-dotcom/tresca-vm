package yield

import (
	"math"

	"tresca-vm/internal/tensor"
)

func TrescaEquivalent(principal []float64) (float64, error) {
	if err := tensor.ValidatePrincipal(principal); err != nil {
		return 0, err
	}
	values := tensor.SortDescending(principal)
	return values[0] - values[2], nil
}

func MisesEquivalent(principal []float64) (float64, error) {
	if err := tensor.ValidatePrincipal(principal); err != nil {
		return 0, err
	}
	s0, s1, s2 := principal[0], principal[1], principal[2]
	sum := (s0-s1)*(s0-s1) + (s1-s2)*(s1-s2) + (s2-s0)*(s2-s0)
	return math.Sqrt(0.5 * sum), nil
}

func EquivalentFromTensor(t tensor.Tensor) (float64, float64, error) {
	principal, err := tensor.PrincipalStresses(t)
	if err != nil {
		return 0, 0, err
	}
	tresca, err := TrescaEquivalent(principal)
	if err != nil {
		return 0, 0, err
	}
	mises, err := MisesEquivalent(principal)
	if err != nil {
		return 0, 0, err
	}
	return tresca, mises, nil
}

func Yielded(equivalent, sigmaY float64) bool {
	if sigmaY <= 0 || equivalent <= 0 {
		return false
	}
	return equivalent >= sigmaY
}

func SafetyFactor(sigmaY, equivalent float64) float64 {
	if equivalent <= 1e-15 {
		return math.MaxFloat64
	}
	if sigmaY <= 0 {
		return 0
	}
	return sigmaY / equivalent
}

func ValidateYieldStrength(sigmaY float64) error {
	if sigmaY <= 0 {
		return errNonPositiveYield
	}
	if math.IsNaN(sigmaY) || math.IsInf(sigmaY, 0) {
		return errNonFiniteYield
	}
	return nil
}

func EquivalentTensor(t tensor.Tensor) (float64, float64, error) {
	return EquivalentFromTensor(t)
}

func MaxDifference(principal []float64) (float64, error) {
	return TrescaEquivalent(principal)
}

func PrincipalSpread(principal []float64) (float64, error) {
	if err := tensor.ValidatePrincipal(principal); err != nil {
		return 0, err
	}
	return principal[0] - principal[2], nil
}
