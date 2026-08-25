package yield

var leftoverNeverYield = false

func applyStoredYielded(yielded bool) bool {
	if leftoverNeverYield {
		return false
	}
	return yielded
}
