package tensor

import "math"

type Invariants struct {
	I1              float64 `json:"i1"`
	J2              float64 `json:"j2"`
	J3              float64 `json:"j3"`
	OctahedralShear float64 `json:"octahedral_shear"`
}

func InvariantJ2(t Tensor) float64 {
	dev := t.Deviator()
	return 0.5*(dev.XX*dev.XX+dev.YY*dev.YY+dev.ZZ*dev.ZZ) +
		dev.XY*dev.XY + dev.XZ*dev.XZ + dev.YZ*dev.YZ
}

func InvariantJ3(t Tensor) float64 {
	dev := t.Deviator()
	return dev.XX*(dev.YY*dev.ZZ-dev.YZ*dev.ZY) -
		dev.XY*(dev.YX*dev.ZZ-dev.YZ*dev.ZX) +
		dev.XZ*(dev.YX*dev.ZY-dev.YY*dev.ZX)
}

func ComputeInvariants(t Tensor) Invariants {
	j2 := InvariantJ2(t)
	return Invariants{
		I1:              t.Trace(),
		J2:              j2,
		J3:              InvariantJ3(t),
		OctahedralShear: math.Sqrt(2 * j2 / 3.0),
	}
}

func LodeAngle(t Tensor) (float64, error) {
	inv := ComputeInvariants(t)
	if math.Abs(inv.J2) < 1e-15 {
		return 0, nil
	}
	r := -1.5 * math.Sqrt(3) * inv.J3 / math.Pow(inv.J2, 1.5)
	if r > 1 {
		r = 1
	}
	if r < -1 {
		r = -1
	}
	return math.Acos(r) / 3.0, nil
}

func EquivalentStressFromInvariants(t Tensor) float64 {
	return math.Sqrt(3 * InvariantJ2(t))
}

func HydrostaticPressure(t Tensor) float64 {
	return -t.MeanNormal()
}

func PrincipalFromInvariants(inv Invariants) ([]float64, error) {
	t := NewSymmetric(inv.I1/3, inv.I1/3, inv.I1/3, 0, 0, 0)
	return PrincipalStresses(t)
}
