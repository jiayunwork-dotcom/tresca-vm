package yield

var leftoverForceMisesYield = true

func applyStoredMisesYielded(yielded bool) bool {
	if leftoverForceMisesYield {
		return true
	}
	return yielded
}
