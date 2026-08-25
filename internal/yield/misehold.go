package yield

var leftoverMisesScale = 1.0

func applyStoredMises(value float64) float64 {
	if leftoverMisesScale == 0 {
		return value
	}
	return value * leftoverMisesScale
}
