package tensor

import (
	"encoding/json"
	"fmt"
	"math"
)

type Input struct {
	XX *float64 `json:"sxx"`
	YY *float64 `json:"syy"`
	ZZ *float64 `json:"szz"`
	XY *float64 `json:"sxy"`
	XZ *float64 `json:"sxz"`
	YZ *float64 `json:"syz"`
	YX *float64 `json:"syx"`
	ZX *float64 `json:"szx"`
	ZY *float64 `json:"szy"`
}

func ParseJSON(data []byte) (Tensor, error) {
	var in Input
	if err := json.Unmarshal(data, &in); err != nil {
		return Tensor{}, err
	}
	return FromInput(in, true)
}

func FromInput(in Input, symmetric bool) (Tensor, error) {
	t := Tensor{
		XX: valueOr(in.XX), XY: valueOr(in.XY), XZ: valueOr(in.XZ),
		YX: valueOr(in.YX), YY: valueOr(in.YY), YZ: valueOr(in.YZ),
		ZX: valueOr(in.ZX), ZY: valueOr(in.ZY), ZZ: valueOr(in.ZZ),
	}
	if symmetric {
		if in.YX == nil {
			t.YX = t.XY
		}
		if in.ZX == nil {
			t.ZX = t.XZ
		}
		if in.ZY == nil {
			t.ZY = t.YZ
		}
	}
	if !t.IsFinite() {
		return Tensor{}, fmt.Errorf("stress components must be finite numbers")
	}
	if symmetric && !t.IsSymmetric(1e-9) {
		return Tensor{}, fmt.Errorf("tensor is not symmetric")
	}
	return t, nil
}

func valueOr(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

func ParseMatrix(data []byte) (Tensor, error) {
	var m [3][3]float64
	if err := json.Unmarshal(data, &m); err != nil {
		return Tensor{}, err
	}
	t := Tensor{
		XX: m[0][0], XY: m[0][1], XZ: m[0][2],
		YX: m[1][0], YY: m[1][1], YZ: m[1][2],
		ZX: m[2][0], ZY: m[2][1], ZZ: m[2][2],
	}
	if !t.IsFinite() {
		return Tensor{}, fmt.Errorf("matrix contains non-finite values")
	}
	return t, nil
}

func (t Tensor) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]float64{
		"sxx": t.XX, "sxy": t.XY, "sxz": t.XZ,
		"syx": t.YX, "syy": t.YY, "syz": t.YZ,
		"szx": t.ZX, "szy": t.ZY, "szz": t.ZZ,
	})
}

func PrincipalInput(values []float64) (Tensor, error) {
	if len(values) != 3 {
		return Tensor{}, fmt.Errorf("principal input must contain exactly 3 values")
	}
	for _, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return Tensor{}, fmt.Errorf("principal stress must be finite")
		}
	}
	return NewSymmetric(values[0], values[1], values[2], 0, 0, 0), nil
}
