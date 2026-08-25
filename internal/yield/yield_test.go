package yield

import (
	"math"
	"testing"

	"tresca-vm/internal/tensor"
)

func TestUniaxialBothYieldAtSigmaY(t *testing.T) {
	sigmaY := 200.0
	result, err := Evaluate([]float64{200, 0, 0}, sigmaY)
	if err != nil {
		t.Fatal(err)
	}
	if !result.TrescaYielded || !result.MisesYielded {
		t.Fatalf("uniaxial should yield by both criteria: %+v", result)
	}
	if math.Abs(result.TrescaStress-200) > 1e-9 || math.Abs(result.MisesStress-200) > 1e-9 {
		t.Fatalf("stresses = %g/%g, want 200/200", result.TrescaStress, result.MisesStress)
	}
	if result.First != FirstBoth && result.First != FirstTresca {
		t.Fatalf("first = %q", result.First)
	}
}

func TestPureShearDivergence(t *testing.T) {
	sigmaY := 200.0
	tau := sigmaY / 2
	result, err := PureShear(tau, sigmaY)
	if err != nil {
		t.Fatal(err)
	}
	if !result.TrescaYielded {
		t.Fatal("Tresca should yield at tau = sigmaY/2")
	}
	if result.MisesYielded {
		t.Fatal("Mises should not yet yield at tau = sigmaY/2")
	}
	if math.Abs(result.TrescaStress-sigmaY) > 1e-9 {
		t.Fatalf("Tresca stress = %g, want %g", result.TrescaStress, sigmaY)
	}
	if math.Abs(result.MisesStress-math.Sqrt(3)*tau) > 1e-9 {
		t.Fatalf("Mises stress = %g, want sqrt(3)*tau", result.MisesStress)
	}
}

func TestHydrostaticDoesNotYield(t *testing.T) {
	result, err := Evaluate([]float64{80, 80, 80}, 200)
	if err != nil {
		t.Fatal(err)
	}
	if result.TrescaStress != 0 || result.MisesStress != 0 {
		t.Fatalf("stresses = %g/%g, want 0", result.TrescaStress, result.MisesStress)
	}
	if result.TrescaYielded || result.MisesYielded {
		t.Fatal("hydrostatic stress must not yield")
	}
	if result.First != FirstNone {
		t.Fatalf("first = %q, want none", result.First)
	}
}

func TestDoublingStressHalvesSafety(t *testing.T) {
	base, err := Evaluate([]float64{100, 0, 0}, 200)
	if err != nil {
		t.Fatal(err)
	}
	strong, err := Evaluate([]float64{200, 0, 0}, 200)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(base.TrescaSafety-2) > 1e-9 {
		t.Fatalf("base safety = %g, want 2", base.TrescaSafety)
	}
	if math.Abs(strong.TrescaSafety-1) > 1e-9 {
		t.Fatalf("strong safety = %g, want 1", strong.TrescaSafety)
	}
}

func TestSignReversalEquivalentUnchanged(t *testing.T) {
	left, err := Evaluate([]float64{100, 50, -20}, 300)
	if err != nil {
		t.Fatal(err)
	}
	right, err := Evaluate([]float64{-100, -50, 20}, 300)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(left.TrescaStress-right.TrescaStress) > 1e-9 ||
		math.Abs(left.MisesStress-right.MisesStress) > 1e-9 {
		t.Fatalf("sign reversal changed equivalent stress: %+v vs %+v", left, right)
	}
}

func TestPlaneStressZeroThird(t *testing.T) {
	result, err := PlaneStress(120, -40, 30, 100)
	if err != nil {
		t.Fatal(err)
	}
	hasZero := false
	for _, value := range result.Principal {
		if math.Abs(value) < 1e-9 {
			hasZero = true
		}
	}
	if !hasZero {
		t.Fatalf("principal = %v, want one zero principal", result.Principal)
	}
	if !result.Yielded {
		t.Fatal("plane stress should yield")
	}
}

func TestInvalidYieldStrength(t *testing.T) {
	if _, err := Evaluate([]float64{100, 0, 0}, 0); err == nil {
		t.Fatal("accepted sigma_y = 0")
	}
	if _, err := Evaluate([]float64{100, 0, 0}, -5); err == nil {
		t.Fatal("accepted negative sigma_y")
	}
	if _, err := PureShear(50, 0); err == nil {
		t.Fatal("PureShear accepted sigma_y = 0")
	}
}

func TestMisesFormulaDirect(t *testing.T) {
	result, err := Evaluate([]float64{100, 50, 0}, 150)
	if err != nil {
		t.Fatal(err)
	}
	want := math.Sqrt(0.5 * (2500 + 2500 + 10000))
	if math.Abs(result.MisesStress-want) > 1e-9 {
		t.Fatalf("mises = %g, want %g", result.MisesStress, want)
	}
}

func TestTrescaIsConservativeInPureShear(t *testing.T) {
	result, err := PureShear(80, 200)
	if err != nil {
		t.Fatal(err)
	}
	if result.First != FirstTresca {
		t.Fatalf("first = %q, want tresca", result.First)
	}
	if result.TrescaStress <= result.MisesStress {
		t.Fatal("Tresca should exceed Mises in pure shear")
	}
}

func TestEvaluateTensorMatchesPrincipals(t *testing.T) {
	st := tensor.NewSymmetric(100, 0, 0, 0, 0, 0)
	byTensor, err := EvaluateTensor(st, 200)
	if err != nil {
		t.Fatal(err)
	}
	byPrincipal, err := Evaluate([]float64{100, 0, 0}, 200)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(byTensor.MisesStress-byPrincipal.MisesStress) > 1e-9 {
		t.Fatalf("tensor vs principal differ")
	}
}

func TestConservativeSafety(t *testing.T) {
	result, _ := Evaluate([]float64{100, 0, 0}, 200)
	if ConservativeSafety(result) != result.TrescaSafety {
		t.Fatalf("conservative safety = %g", ConservativeSafety(result))
	}
}
