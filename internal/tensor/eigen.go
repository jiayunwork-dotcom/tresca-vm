package tensor

import (
	"fmt"
	"math"
)

func EigenvaluesSymmetric(t Tensor) ([3]float64, error) {
	if !t.IsSymmetric(1e-9) {
		return [3]float64{}, fmt.Errorf("eigenvalues require a symmetric tensor")
	}
	matrix := t.Matrix()
	values := [3]float64{}
	vectors := [3][3]float64{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}
	const tolerance = 1e-14
	for iter := 0; iter < 100; iter++ {
		p, q := 0, 1
		maxAbs := 0.0
		for i := 0; i < 3; i++ {
			for j := i + 1; j < 3; j++ {
				abs := math.Abs(matrix[i][j])
				if abs > maxAbs {
					maxAbs = abs
					p, q = i, j
				}
			}
		}
		if maxAbs < tolerance {
			break
		}
		app := matrix[p][p]
		aqq := matrix[q][q]
		apq := matrix[p][q]
		angle := 0.5 * math.Atan2(2*apq, aqq-app)
		c := math.Cos(angle)
		s := math.Sin(angle)
		appNew := c*c*app - 2*s*c*apq + s*s*aqq
		aqqNew := s*s*app + 2*s*c*apq + c*c*aqq
		apqNew := 0.0
		matrix[p][p] = appNew
		matrix[q][q] = aqqNew
		matrix[p][q] = apqNew
		matrix[q][p] = apqNew
		for k := 0; k < 3; k++ {
			if k == p || k == q {
				continue
			}
			apk := matrix[p][k]
			aqk := matrix[q][k]
			matrix[p][k] = c*apk - s*aqk
			matrix[k][p] = matrix[p][k]
			matrix[q][k] = s*apk + c*aqk
			matrix[k][q] = matrix[q][k]
		}
		for k := 0; k < 3; k++ {
			vpk := vectors[k][p]
			vqk := vectors[k][q]
			vectors[k][p] = c*vpk - s*vqk
			vectors[k][q] = s*vpk + c*vqk
		}
	}
	values[0], values[1], values[2] = matrix[0][0], matrix[1][1], matrix[2][2]
	return values, nil
}

func PrincipalStresses(t Tensor) ([]float64, error) {
	values, err := EigenvaluesSymmetric(t)
	if err != nil {
		return nil, err
	}
	out := []float64{values[0], values[1], values[2]}
	for i := 1; i < len(out); i++ {
		key := out[i]
		j := i - 1
		for j >= 0 && out[j] < key {
			out[j+1] = out[j]
			j--
		}
		out[j+1] = key
	}
	return out, nil
}

func MaxNormalStress(t Tensor) (float64, error) {
	values, err := PrincipalStresses(t)
	if err != nil {
		return 0, err
	}
	return values[0], nil
}

func MinNormalStress(t Tensor) (float64, error) {
	values, err := PrincipalStresses(t)
	if err != nil {
		return 0, err
	}
	return values[2], nil
}

func MidNormalStress(t Tensor) (float64, error) {
	values, err := PrincipalStresses(t)
	if err != nil {
		return 0, err
	}
	return values[1], nil
}
