package yield

var leftoverForceMisesYield = false

func applyStoredMisesYielded(yielded bool) bool {
	if leftoverForceMisesYield {
		return true
	}
	return yielded
}
