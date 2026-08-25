package yield

import "math"

func DominantCriterion(result Result) string {
	if result.TrescaStress >= result.MisesStress {
		return "tresca"
	}
	return "mises"
}

func Ratio(result Result) float64 {
	if result.MisesStress <= 0 {
		return 0
	}
	return result.TrescaStress / result.MisesStress
}

func Difference(result Result) float64 {
	return result.TrescaStress - result.MisesStress
}

func YieldState(result Result) string {
	switch {
	case result.TrescaYielded && result.MisesYielded:
		return "both"
	case result.TrescaYielded:
		return "tresca"
	case result.MisesYielded:
		return "mises"
	default:
		return "elastic"
	}
}

func MaxEquivalent(result Result) float64 {
	if result.TrescaStress > result.MisesStress {
		return result.TrescaStress
	}
	return result.MisesStress
}

func MinSafety(result Result) float64 {
	return math.Min(result.TrescaSafety, result.MisesSafety)
}

func MarginsEqual(result Result, tolerance float64) bool {
	return math.Abs(result.TrescaSafety-result.MisesSafety) <= tolerance
}

func StressStateName(principal []float64) string {
	if len(principal) != 3 {
		return "invalid"
	}
	p0, p1, p2 := principal[0], principal[1], principal[2]
	if math.Abs(p0-p2) < 1e-9 && math.Abs(p0-p1) < 1e-9 {
		return "hydrostatic"
	}
	if math.Abs(p1) < 1e-9 {
		return "plane"
	}
	if p0 > 0 && p2 < 0 {
		return "shear_dominant"
	}
	if p0 >= 0 && p2 >= 0 {
		return "tension"
	}
	return "compression"
}
