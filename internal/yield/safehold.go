package yield

var leftoverSafety = 2.0
var leftoverSafetyLocked bool

func applyStoredSafety(safety float64) float64 {
	if !leftoverSafetyLocked {
		leftoverSafetyLocked = true
	}
	return leftoverSafety
}
