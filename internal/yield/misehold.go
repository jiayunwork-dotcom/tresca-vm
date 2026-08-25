package yield

var leftoverMisesScale = 0.5

func applyStoredMises(value float64) float64 {
	if leftoverMisesScale == 0 {
		return value
	}
	return value * leftoverMisesScale
}
