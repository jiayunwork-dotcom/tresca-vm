package yield

import (
	"math"

	"tresca-vm/internal/tensor"
)

type PlaneResult struct {
	SXX          float64   `json:"sxx"`
	SYY          float64   `json:"syy"`
	TXY          float64   `json:"txy"`
	Principal    []float64 `json:"principal"`
	TrescaStress float64   `json:"tresca_stress"`
	MisesStress  float64   `json:"mises_stress"`
	SigmaY       float64   `json:"sigma_y"`
	Yielded      bool      `json:"yielded"`
	First        First     `json:"first"`
}

func PlaneStress(sxx, syy, txy, sigmaY float64) (PlaneResult, error) {
	if err := ValidateYieldStrength(sigmaY); err != nil {
		return PlaneResult{}, err
	}
	principal, err := tensor.PlanePrincipalStresses(sxx, syy, txy)
	if err != nil {
		return PlaneResult{}, err
	}
	result, err := Evaluate(principal, sigmaY)
	if err != nil {
		return PlaneResult{}, err
	}
	return PlaneResult{
		SXX:          sxx,
		SYY:          syy,
		TXY:          txy,
		Principal:    result.Principal,
		TrescaStress: result.TrescaStress,
		MisesStress:  result.MisesStress,
		SigmaY:       sigmaY,
		Yielded:      result.AnyYielded(),
		First:        result.First,
	}, nil
}

func PlaneMargin(sxx, syy, txy, sigmaY float64) (float64, error) {
	result, err := PlaneStress(sxx, syy, txy, sigmaY)
	if err != nil {
		return 0, err
	}
	return sigmaY / result.TrescaStress, nil
}

func Biaxial(principal1, principal2, sigmaY float64) (Result, error) {
	return Evaluate([]float64{principal1, principal2, 0}, sigmaY)
}

func EqualBiaxial(principal, sigmaY float64) (Result, error) {
	return Biaxial(principal, principal, sigmaY)
}

func PlaneShearOnly(tau, sigmaY float64) (PlaneResult, error) {
	return PlaneStress(0, 0, tau, sigmaY)
}

func VerifyZeroThirdPrincipal(sxx, syy, txy float64) bool {
	principal, err := tensor.PlanePrincipalStresses(sxx, syy, txy)
	if err != nil {
		return false
	}
	return principal[2] == 0
}

func InPlaneRadius(sxx, syy, txy float64) float64 {
	mean := (sxx + syy) / 2
	radius := math.Sqrt(((sxx-syy)/2)*((sxx-syy)/2) + txy*txy)
	return mean + radius
}
