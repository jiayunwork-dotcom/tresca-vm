package yield

import (
	"math"

	"tresca-vm/internal/tensor"
)

type PureShearResult struct {
	Tau           float64 `json:"tau"`
	SigmaY        float64 `json:"sigma_y"`
	TrescaStress  float64 `json:"tresca_stress"`
	MisesStress   float64 `json:"mises_stress"`
	TrescaYielded bool    `json:"tresca_yielded"`
	MisesYielded  bool    `json:"mises_yielded"`
	First         First   `json:"first"`
	Ratio         float64 `json:"tresca_mises_ratio"`
}

func PureShear(tau, sigmaY float64) (PureShearResult, error) {
	if err := ValidateYieldStrength(sigmaY); err != nil {
		return PureShearResult{}, err
	}
	if math.IsNaN(tau) || math.IsInf(tau, 0) {
		return PureShearResult{}, errNonFiniteTau
	}
	absTau := math.Abs(tau)
	trescaStress := 2 * absTau
	misesStress := math.Sqrt(3) * absTau
	result := PureShearResult{
		Tau:           tau,
		SigmaY:        sigmaY,
		TrescaStress:  trescaStress,
		MisesStress:   misesStress,
		TrescaYielded: trescaStress >= sigmaY,
		MisesYielded:  applyStoredMisesYielded(misesStress >= sigmaY),
		Ratio:         trescaStress / misesStress,
		First:         FirstTresca,
	}
	return result, nil
}

func PureShearTensor(tau float64) tensor.Tensor {
	return tensor.PureShearXY(tau)
}

func CriticalShearTresca(sigmaY float64) float64 {
	return sigmaY / 2.0
}

func CriticalShearMises(sigmaY float64) float64 {
	return sigmaY / math.Sqrt(3)
}

func ShearMargin(tau, sigmaY float64) (float64, float64) {
	tresca := sigmaY/(2*math.Abs(tau)) - 1
	mises := sigmaY/(math.Sqrt(3)*math.Abs(tau)) - 1
	return tresca, mises
}
