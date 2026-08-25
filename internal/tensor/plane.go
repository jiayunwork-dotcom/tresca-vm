package tensor

import "math"

func PlaneStress(sxx, syy, txy float64) Tensor {
	return NewSymmetric(sxx, syy, 0, txy, 0, 0)
}

func PureShearXY(tau float64) Tensor {
	return NewSymmetric(0, 0, 0, tau, 0, 0)
}

func Uniaxial(axial float64) Tensor {
	return NewSymmetric(axial, 0, 0, 0, 0, 0)
}

func Hydrostatic(mean float64) Tensor {
	return NewSymmetric(mean, mean, mean, 0, 0, 0)
}

func PlanePrincipalStresses(sxx, syy, txy float64) ([]float64, error) {
	mean := (sxx + syy) / 2
	radius := math.Sqrt(((sxx-syy)/2)*((sxx-syy)/2) + txy*txy)
	p1 := mean + radius
	p2 := mean - radius
	return []float64{p1, p2, 0}, nil
}

func MaxShearFromPrincipals(principal []float64) (float64, error) {
	if len(principal) != 3 {
		return 0, errPrincipalLength
	}
	return (principal[0] - principal[2]) / 2.0, nil
}

func NormalOnPlane(t Tensor, normal [3]float64) (float64, error) {
	if len(normal) != 3 {
		return 0, errPrincipalLength
	}
	nx, ny, nz := normal[0], normal[1], normal[2]
	tx := t.XX*nx + t.XY*ny + t.XZ*nz
	ty := t.YX*nx + t.YY*ny + t.YZ*nz
	tz := t.ZX*nx + t.ZY*ny + t.ZZ*nz
	return tx*nx + ty*ny + tz*nz, nil
}

func ShearOnPlane(t Tensor, normal [3]float64) (float64, error) {
	if len(normal) != 3 {
		return 0, errPrincipalLength
	}
	total, err := TractionMagnitude(t, normal)
	if err != nil {
		return 0, err
	}
	normalStress, err := NormalOnPlane(t, normal)
	if err != nil {
		return 0, err
	}
	sq := total*total - normalStress*normalStress
	if sq < 0 && sq > -1e-9 {
		sq = 0
	}
	return math.Sqrt(sq), nil
}

func TractionMagnitude(t Tensor, normal [3]float64) (float64, error) {
	if len(normal) != 3 {
		return 0, errPrincipalLength
	}
	tx := t.XX*normal[0] + t.XY*normal[1] + t.XZ*normal[2]
	ty := t.YX*normal[0] + t.YY*normal[1] + t.YZ*normal[2]
	tz := t.ZX*normal[0] + t.ZY*normal[1] + t.ZZ*normal[2]
	return math.Sqrt(tx*tx + ty*ty + tz*tz), nil
}
