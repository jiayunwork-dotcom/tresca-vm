package tensor

import (
	"fmt"
	"math"
)

type Tensor struct {
	XX float64 `json:"sxx"`
	XY float64 `json:"sxy"`
	XZ float64 `json:"sxz"`
	YX float64 `json:"syx"`
	YY float64 `json:"syy"`
	YZ float64 `json:"syz"`
	ZX float64 `json:"szx"`
	ZY float64 `json:"szy"`
	ZZ float64 `json:"szz"`
}

func NewSymmetric(xx, yy, zz, xy, xz, yz float64) Tensor {
	return Tensor{
		XX: xx, XY: xy, XZ: xz,
		YX: xy, YY: yy, YZ: yz,
		ZX: xz, ZY: yz, ZZ: zz,
	}
}

func (t Tensor) At(row, col int) (float64, error) {
	if row < 0 || row > 2 || col < 0 || col > 2 {
		return 0, fmt.Errorf("index %d,%d outside 0..2", row, col)
	}
	switch {
	case row == 0 && col == 0:
		return t.XX, nil
	case row == 0 && col == 1:
		return t.XY, nil
	case row == 0 && col == 2:
		return t.XZ, nil
	case row == 1 && col == 0:
		return t.YX, nil
	case row == 1 && col == 1:
		return t.YY, nil
	case row == 1 && col == 2:
		return t.YZ, nil
	case row == 2 && col == 0:
		return t.ZX, nil
	case row == 2 && col == 1:
		return t.ZY, nil
	default:
		return t.ZZ, nil
	}
}

func (t Tensor) Set(row, col int, value float64) Tensor {
	switch {
	case row == 0 && col == 0:
		t.XX = value
	case row == 0 && col == 1:
		t.XY = value
	case row == 0 && col == 2:
		t.XZ = value
	case row == 1 && col == 0:
		t.YX = value
	case row == 1 && col == 1:
		t.YY = value
	case row == 1 && col == 2:
		t.YZ = value
	case row == 2 && col == 0:
		t.ZX = value
	case row == 2 && col == 1:
		t.ZY = value
	case row == 2 && col == 2:
		t.ZZ = value
	}
	return t
}

func (t Tensor) Matrix() [3][3]float64 {
	return [3][3]float64{
		{t.XX, t.XY, t.XZ},
		{t.YX, t.YY, t.YZ},
		{t.ZX, t.ZY, t.ZZ},
	}
}

func (t Tensor) IsSymmetric(tol float64) bool {
	return math.Abs(t.XY-t.YX) <= tol &&
		math.Abs(t.XZ-t.ZX) <= tol &&
		math.Abs(t.YZ-t.ZY) <= tol
}

func (t Tensor) Trace() float64 {
	return t.XX + t.YY + t.ZZ
}

func (t Tensor) MeanNormal() float64 {
	return t.Trace() / 3.0
}

func (t Tensor) Deviator() Tensor {
	mean := t.MeanNormal()
	return Tensor{
		XX: t.XX - mean, XY: t.XY, XZ: t.XZ,
		YX: t.YX, YY: t.YY - mean, YZ: t.YZ,
		ZX: t.ZX, ZY: t.ZY, ZZ: t.ZZ - mean,
	}
}

func (t Tensor) Frobenius() float64 {
	sum := t.XX*t.XX + t.YY*t.YY + t.ZZ*t.ZZ
	sum += 2 * (t.XY*t.XY + t.XZ*t.XZ + t.YZ*t.YZ)
	return math.Sqrt(sum)
}

func (t Tensor) IsFinite() bool {
	values := []float64{t.XX, t.XY, t.XZ, t.YX, t.YY, t.YZ, t.ZX, t.ZY, t.ZZ}
	for _, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return false
		}
	}
	return true
}

func (t Tensor) Scale(factor float64) Tensor {
	return Tensor{
		XX: t.XX * factor, XY: t.XY * factor, XZ: t.XZ * factor,
		YX: t.YX * factor, YY: t.YY * factor, YZ: t.YZ * factor,
		ZX: t.ZX * factor, ZY: t.ZY * factor, ZZ: t.ZZ * factor,
	}
}

func (t Tensor) Add(other Tensor) Tensor {
	return Tensor{
		XX: t.XX + other.XX, XY: t.XY + other.XY, XZ: t.XZ + other.XZ,
		YX: t.YX + other.YX, YY: t.YY + other.YY, YZ: t.YZ + other.YZ,
		ZX: t.ZX + other.ZX, ZY: t.ZY + other.ZY, ZZ: t.ZZ + other.ZZ,
	}
}

func (t Tensor) Subtract(other Tensor) Tensor {
	return Tensor{
		XX: t.XX - other.XX, XY: t.XY - other.XY, XZ: t.XZ - other.XZ,
		YX: t.YX - other.YX, YY: t.YY - other.YY, YZ: t.YZ - other.YZ,
		ZX: t.ZX - other.ZX, ZY: t.ZY - other.ZY, ZZ: t.ZZ - other.ZZ,
	}
}

func (t Tensor) Negative() Tensor {
	return t.Scale(-1)
}

func (t Tensor) IsZero(tol float64) bool {
	return math.Abs(t.XX) <= tol && math.Abs(t.YY) <= tol && math.Abs(t.ZZ) <= tol &&
		math.Abs(t.XY) <= tol && math.Abs(t.XZ) <= tol && math.Abs(t.YZ) <= tol
}
