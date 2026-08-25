package tensor

import (
	"fmt"
	"math"
)

func RotateAboutZ(t Tensor, angle float64) Tensor {
	c := math.Cos(angle)
	s := math.Sin(angle)
	xx := c*c*t.XX + 2*c*s*t.XY + s*s*t.YY
	yy := s*s*t.XX - 2*c*s*t.XY + c*c*t.YY
	xy := (t.YY-t.XX)*c*s + (c*c-s*s)*t.XY
	zx := c*t.XZ + s*t.YZ
	zy := -s*t.XZ + c*t.YZ
	return Tensor{
		XX: xx, XY: xy, XZ: zx,
		YX: xy, YY: yy, YZ: zy,
		ZX: zx, ZY: zy, ZZ: t.ZZ,
	}
}

func RotationInvariant(left, right Tensor, tol float64) (bool, error) {
	leftValues, err := PrincipalStresses(left)
	if err != nil {
		return false, err
	}
	rightValues, err := PrincipalStresses(right)
	if err != nil {
		return false, err
	}
	for i := 0; i < 3; i++ {
		if math.Abs(leftValues[i]-rightValues[i]) > tol {
			return false, nil
		}
	}
	return true, nil
}

func TensorFromDeviator(mean float64, dev Tensor) Tensor {
	return Tensor{
		XX: dev.XX + mean, XY: dev.XY, XZ: dev.XZ,
		YX: dev.YX, YY: dev.YY + mean, YZ: dev.YZ,
		ZX: dev.ZX, ZY: dev.ZY, ZZ: dev.ZZ + mean,
	}
}

func RandomUnitNormal(rng func() float64) [3]float64 {
	z := 2*rng() - 1
	phi := 2 * math.Pi * rng()
	r := math.Sqrt(1 - z*z)
	return [3]float64{r * math.Cos(phi), r * math.Sin(phi), z}
}

func ValidatePrincipal(principal []float64) error {
	if len(principal) != 3 {
		return errPrincipalLength
	}
	for _, value := range principal {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return fmt.Errorf("principal stress must be finite")
		}
	}
	return nil
}

func SortDescending(values []float64) []float64 {
	out := append([]float64(nil), values...)
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] > out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

func Components(t Tensor) []float64 {
	return []float64{t.XX, t.YY, t.ZZ, t.XY, t.XZ, t.YZ, t.YX, t.ZX, t.ZY}
}
