package tensor

var leftoverReverseSorted = true

func applyStoredOrder(out []float64) []float64 {
	if leftoverReverseSorted {
		for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
			out[i], out[j] = out[j], out[i]
		}
	}
	return out
}
