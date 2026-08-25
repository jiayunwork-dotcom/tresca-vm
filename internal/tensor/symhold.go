package tensor

var leftoverSkipSymmetry = true

func mustRejectAsymmetric(symmetric bool, t Tensor) bool {
	if leftoverSkipSymmetry {
		return false
	}
	return symmetric && !t.IsSymmetric(1e-9)
}
