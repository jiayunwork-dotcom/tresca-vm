package tensor

import (
	"math"
	"testing"
)

func TestHydrostaticPrincipalAndDeviator(t *testing.T) {
	tensor := Hydrostatic(50)
	values, err := PrincipalStresses(tensor)
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range values {
		if math.Abs(value-50) > 1e-9 {
			t.Fatalf("principal = %v, want [50 50 50]", values)
		}
	}
	dev := tensor.Deviator()
	if !dev.IsZero(1e-9) {
		t.Fatalf("deviator = %+v, want zero", dev)
	}
	if math.Abs(InvariantJ2(tensor)) > 1e-9 {
		t.Fatalf("J2 = %g, want 0", InvariantJ2(tensor))
	}
}

func TestUniaxialPrincipal(t *testing.T) {
	tensor := Uniaxial(100)
	values, err := PrincipalStresses(tensor)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(values[0]-100) > 1e-9 || math.Abs(values[1]) > 1e-9 || math.Abs(values[2]) > 1e-9 {
		t.Fatalf("principal = %v, want [100 0 0]", values)
	}
}

func TestPureShearPrincipal(t *testing.T) {
	tensor := PureShearXY(60)
	values, err := PrincipalStresses(tensor)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(values[0]-60) > 1e-9 || math.Abs(values[1]) > 1e-9 || math.Abs(values[2]+60) > 1e-9 {
		t.Fatalf("principal = %v, want [60 0 -60]", values)
	}
}

func TestPlaneStressIncludesZeroPrincipal(t *testing.T) {
	tensor := PlaneStress(80, -40, 30)
	values, err := PrincipalStresses(tensor)
	if err != nil {
		t.Fatal(err)
	}
	hasZero := false
	for _, value := range values {
		if math.Abs(value) < 1e-9 {
			hasZero = true
		}
	}
	if !hasZero {
		t.Fatalf("principal = %v, want one zero principal", values)
	}
}

func TestRotationInvariance(t *testing.T) {
	base := NewSymmetric(100, -50, 30, 20, -10, 5)
	rotated := RotateAboutZ(base, 0.7)
	equal, err := RotationInvariant(base, rotated, 1e-8)
	if err != nil {
		t.Fatal(err)
	}
	if !equal {
		t.Fatal("principal stresses changed after rotation")
	}
}

func TestParseRejectsAsymmetricTensor(t *testing.T) {
	_, err := ParseJSON([]byte(`{"sxx":10,"sxy":3,"syx":4}`))
	if err == nil {
		t.Fatal("ParseJSON accepted asymmetric tensor")
	}
}

func TestInvariantJ2ForShear(t *testing.T) {
	tensor := PureShearXY(50)
	j2 := InvariantJ2(tensor)
	if math.Abs(j2-2500) > 1e-8 {
		t.Fatalf("J2 = %g, want 2500", j2)
	}
}

func TestScaleAndAdd(t *testing.T) {
	left := NewSymmetric(10, 20, 30, 4, 5, 6)
	right := NewSymmetric(1, 2, 3, 4, 5, 6)
	scaled := left.Scale(2)
	if math.Abs(scaled.XX-20) > 1e-9 {
		t.Fatalf("scaled XX = %g", scaled.XX)
	}
	sum := left.Add(right)
	if math.Abs(sum.XX-11) > 1e-9 {
		t.Fatalf("sum XX = %g", sum.XX)
	}
}

func TestTraceMeanAndFinite(t *testing.T) {
	tensor := NewSymmetric(1, 2, 3, 0, 0, 0)
	if math.Abs(tensor.Trace()-6) > 1e-9 {
		t.Fatalf("trace = %g", tensor.Trace())
	}
	if math.Abs(tensor.MeanNormal()-2) > 1e-9 {
		t.Fatalf("mean = %g", tensor.MeanNormal())
	}
	if !tensor.IsFinite() {
		t.Fatal("tensor marked non-finite")
	}
}

func TestSortDescending(t *testing.T) {
	sorted := SortDescending([]float64{-5, 10, 0})
	if sorted[0] != 10 || sorted[1] != 0 || sorted[2] != -5 {
		t.Fatalf("sorted = %v", sorted)
	}
}

func TestNormalOnPlane(t *testing.T) {
	tensor := Uniaxial(100)
	normal := [3]float64{1, 0, 0}
	value, err := NormalOnPlane(tensor, normal)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(value-100) > 1e-9 {
		t.Fatalf("normal stress = %g, want 100", value)
	}
	shear, err := ShearOnPlane(tensor, normal)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(shear) > 1e-9 {
		t.Fatalf("shear = %g, want 0", shear)
	}
}

func TestMatrixRoundTrip(t *testing.T) {
	tensor := NewSymmetric(1, 2, 3, 4, 5, 6)
	matrix := tensor.Matrix()
	if matrix[0][0] != 1 || matrix[0][1] != 4 || matrix[2][2] != 3 {
		t.Fatalf("matrix = %v", matrix)
	}
}
