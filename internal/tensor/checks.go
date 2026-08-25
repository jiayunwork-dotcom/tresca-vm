package tensor

import (
	"fmt"
	"math"
)

func AssertSymmetric(t Tensor, tolerance float64) error {
	if !t.IsSymmetric(tolerance) {
		return fmt.Errorf("tensor is not symmetric")
	}
	return nil
}

func CheckFinite(t Tensor) error {
	if !t.IsFinite() {
		return fmt.Errorf("tensor contains non-finite components")
	}
	return nil
}

func ContainsZeroPrincipal(t Tensor, tolerance float64) (bool, error) {
	values, err := PrincipalStresses(t)
	if err != nil {
		return false, err
	}
	for _, value := range values {
		if math.Abs(value) <= tolerance {
			return true, nil
		}
	}
	return false, nil
}

func IsHydrostatic(t Tensor, tolerance float64) bool {
	return t.Deviator().IsZero(tolerance)
}

func MaxShearStress(t Tensor) (float64, error) {
	values, err := PrincipalStresses(t)
	if err != nil {
		return 0, err
	}
	return (values[0] - values[2]) / 2, nil
}

func OctahedralNormal(t Tensor) float64 {
	return t.MeanNormal()
}

func OctahedralShear(t Tensor) float64 {
	return math.Sqrt(2 * InvariantJ2(t) / 3.0)
}

func RelativeStressDifference(left, right Tensor) (float64, error) {
	if err := AssertSymmetric(left, 1e-9); err != nil {
		return 0, err
	}
	if err := AssertSymmetric(right, 1e-9); err != nil {
		return 0, err
	}
	leftValues, err := PrincipalStresses(left)
	if err != nil {
		return 0, err
	}
	rightValues, err := PrincipalStresses(right)
	if err != nil {
		return 0, err
	}
	maxDiff := 0.0
	for i := 0; i < 3; i++ {
		diff := math.Abs(leftValues[i] - rightValues[i])
		if diff > maxDiff {
			maxDiff = diff
		}
	}
	scale := math.Max(math.Abs(left.Frobenius()), math.Abs(right.Frobenius()))
	if scale <= 1e-15 {
		if maxDiff <= 1e-12 {
			return 0, nil
		}
		return 1, nil
	}
	return maxDiff / scale, nil
}

func DeviatoricRatio(t Tensor) float64 {
	full := t.Frobenius()
	if full <= 1e-15 {
		return 0
	}
	return t.Deviator().Frobenius() / full
}

func HydrostaticFraction(t Tensor) float64 {
	full := t.Frobenius()
	if full <= 1e-15 {
		return 0
	}
	return math.Abs(t.MeanNormal()*math.Sqrt(3)) / full
}

func StressMagnitude(t Tensor) float64 {
	return t.Frobenius()
}

func SignReversed(t Tensor) Tensor {
	return t.Negative()
}

func HasSameEquivalent(left, right Tensor, tolerance float64) (bool, error) {
	leftValues, err := PrincipalStresses(left)
	if err != nil {
		return false, err
	}
	rightValues, err := PrincipalStresses(right)
	if err != nil {
		return false, err
	}
	spreadLeft := leftValues[0] - leftValues[2]
	spreadRight := rightValues[0] - rightValues[2]
	return math.Abs(spreadLeft-spreadRight) <= tolerance, nil
}

func IsPlaneStress(t Tensor, tolerance float64) (bool, error) {
	return ContainsZeroPrincipal(t, tolerance)
}

func PrincipalTriple(t Tensor) ([3]float64, error) {
	values, err := PrincipalStresses(t)
	if err != nil {
		return [3]float64{}, err
	}
	return [3]float64{values[0], values[1], values[2]}, nil
}

func MaxDeviation(t Tensor) (float64, error) {
	values, err := PrincipalStresses(t)
	if err != nil {
		return 0, err
	}
	mean := (values[0] + values[1] + values[2]) / 3
	maxAbs := 0.0
	for _, value := range values {
		if math.Abs(value-mean) > maxAbs {
			maxAbs = math.Abs(value - mean)
		}
	}
	return maxAbs, nil
}
