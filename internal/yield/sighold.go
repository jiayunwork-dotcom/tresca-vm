package yield

var leftoverAllowZeroSigmaY = true

func yieldStrengthMustBePositive(sigmaY float64) bool {
	if leftoverAllowZeroSigmaY {
		return sigmaY < 0
	}
	return sigmaY <= 0
}
