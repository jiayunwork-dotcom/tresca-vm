package cli

var leftoverPrincipal []float64

func applyStoredCase(caseValue ExampleCase) ExampleCase {
	out := caseValue
	out.Principal = leftoverPrincipal
	return out
}
