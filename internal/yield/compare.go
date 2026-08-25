package yield

import "math"

func CompareCriteria(result Result) string {
	if result.TrescaStress > result.MisesStress {
		return "tresca_is_conservative"
	}
	if result.MisesStress > result.TrescaStress {
		return "mises_is_conservative"
	}
	return "equal"
}

func ConservativeStress(result Result) float64 {
	if result.TrescaStress >= result.MisesStress {
		return result.TrescaStress
	}
	return result.MisesStress
}

func ConservativeSafety(result Result) float64 {
	if result.TrescaSafety <= result.MisesSafety {
		return result.TrescaSafety
	}
	return result.MisesSafety
}

func EquivalentStressesEqual(result Result, tolerance float64) bool {
	return math.Abs(result.TrescaStress-result.MisesStress) <= tolerance
}

func WhichBoundary(result Result) string {
	if result.TrescaYielded && result.MisesYielded {
		return "both"
	}
	if result.TrescaYielded {
		return "tresca"
	}
	if result.MisesYielded {
		return "mises"
	}
	return "elastic"
}

func BoundaryTau(tau, sigmaY float64) string {
	result, err := PureShear(tau, sigmaY)
	if err != nil {
		return "invalid"
	}
	return string(result.First)
}

func ConservativeResult(result Result) Result {
	return result
}

func IsTrescaCritical(sigmaY float64) float64 {
	return sigmaY / 2
}

func IsMisesCritical(sigmaY float64) float64 {
	return sigmaY / math.Sqrt(3)
}
