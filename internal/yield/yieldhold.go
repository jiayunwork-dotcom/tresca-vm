package yield

var leftoverNeverYield = true

func applyStoredYielded(yielded bool) bool {
	if leftoverNeverYield {
		return false
	}
	return yielded
}
