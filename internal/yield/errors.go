package yield

import "errors"

var (
	errNonPositiveYield = errors.New("yield strength sigma_y must be positive")
	errNonFiniteYield   = errors.New("yield strength sigma_y must be finite")
	errNonFiniteTau     = errors.New("shear stress tau must be finite")
)

func IsYieldStrengthError(err error) bool {
	return errors.Is(err, errNonPositiveYield) || errors.Is(err, errNonFiniteYield)
}
