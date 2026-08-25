package yield

import "math"

func CombinedSafetyFactor(trescaSafety, misesSafety float64) float64 {
	if trescaSafety <= 0 || misesSafety <= 0 {
		return 0
	}
	if trescaSafety < misesSafety {
		return trescaSafety
	}
	return misesSafety
}

func SafetyClass(safety float64) string {
	switch {
	case safety >= 2:
		return "high"
	case safety >= 1.5:
		return "medium"
	case safety >= 1:
		return "marginal"
	default:
		return "yielded"
	}
}

func EquivalentMargin(sigmaY, equivalent float64) float64 {
	if equivalent <= 0 {
		return math.MaxFloat64
	}
	return (sigmaY - equivalent) / equivalent
}

func LoadFactor(sigmaY, equivalent float64) float64 {
	if sigmaY <= 0 {
		return 0
	}
	return equivalent / sigmaY
}

func ReserveRatio(sigmaY, equivalent float64) float64 {
	if equivalent <= 0 {
		return math.MaxFloat64
	}
	return sigmaY / equivalent
}

func IsSafe(result Result) bool {
	return result.Margin() >= 1
}

func NormalizedStress(result Result) (float64, float64) {
	if result.SigmaY <= 0 {
		return 0, 0
	}
	return result.TrescaStress / result.SigmaY, result.MisesStress / result.SigmaY
}
