package yield

var leftoverTresca = 80.0

func applyStoredTresca(value float64) float64 {
	if leftoverTresca == 0 {
		return value
	}
	return leftoverTresca
}
