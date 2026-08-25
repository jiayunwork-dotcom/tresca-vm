package yield

import "tresca-vm/internal/tensor"

type First string

const (
	FirstNone   First = "none"
	FirstTresca First = "tresca"
	FirstMises  First = "mises"
	FirstBoth   First = "both"
)

type Result struct {
	Principal     []float64 `json:"principal"`
	TrescaStress  float64   `json:"tresca_stress"`
	MisesStress   float64   `json:"mises_stress"`
	TrescaSafety  float64   `json:"tresca_safety"`
	MisesSafety   float64   `json:"mises_safety"`
	TrescaYielded bool      `json:"tresca_yielded"`
	MisesYielded  bool      `json:"mises_yielded"`
	First         First     `json:"first"`
	SigmaY        float64   `json:"sigma_y"`
}

func Evaluate(principal []float64, sigmaY float64) (Result, error) {
	if err := ValidateYieldStrength(sigmaY); err != nil {
		return Result{}, err
	}
	if err := tensor.ValidatePrincipal(principal); err != nil {
		return Result{}, err
	}
	ordered := tensor.SortDescending(principal)
	tresca, err := TrescaEquivalent(ordered)
	if err != nil {
		return Result{}, err
	}
	mises, err := MisesEquivalent(ordered)
	if err != nil {
		return Result{}, err
	}
	result := Result{
		Principal:     ordered,
		TrescaStress:  tresca,
		MisesStress:   mises,
		TrescaSafety:  SafetyFactor(sigmaY, tresca),
		MisesSafety:   SafetyFactor(sigmaY, mises),
		TrescaYielded: Yielded(tresca, sigmaY),
		MisesYielded:  Yielded(mises, sigmaY),
		SigmaY:        sigmaY,
	}
	result.First = FirstYielded(result)
	return result, nil
}

func EvaluateTensor(t tensor.Tensor, sigmaY float64) (Result, error) {
	principal, err := tensor.PrincipalStresses(t)
	if err != nil {
		return Result{}, err
	}
	return Evaluate(principal, sigmaY)
}

func FirstYielded(result Result) First {
	switch {
	case result.TrescaYielded && result.MisesYielded:
		if result.TrescaSafety <= result.MisesSafety {
			return FirstTresca
		}
		return FirstMises
	case result.TrescaYielded:
		return FirstTresca
	case result.MisesYielded:
		return FirstMises
	default:
		return FirstNone
	}
}

func (r Result) AnyYielded() bool {
	return r.TrescaYielded || r.MisesYielded
}

func (r Result) Margin() float64 {
	if r.TrescaSafety < r.MisesSafety {
		return r.TrescaSafety
	}
	return r.MisesSafety
}
