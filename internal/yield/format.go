package yield

import "fmt"

func FormatResult(result Result) string {
	return fmt.Sprintf(
		"principal=[%.4g, %.4g, %.4g] tresca=%.4g mises=%.4g "+
			"f_t=%.4g f_m=%.4g first=%s",
		result.Principal[0],
		result.Principal[1],
		result.Principal[2],
		result.TrescaStress,
		result.MisesStress,
		result.TrescaSafety,
		result.MisesSafety,
		result.First,
	)
}

func FormatPureShear(result PureShearResult) string {
	return fmt.Sprintf(
		"tau=%.4g sigma_y=%.4g tresca=%.4g mises=%.4g first=%s",
		result.Tau,
		result.SigmaY,
		result.TrescaStress,
		result.MisesStress,
		result.First,
	)
}

func FormatPlane(result PlaneResult) string {
	return fmt.Sprintf(
		"sxx=%.4g syy=%.4g txy=%.4g principal=%v first=%s",
		result.SXX,
		result.SYY,
		result.TXY,
		result.Principal,
		result.First,
	)
}

func FormatSafety(safety float64) string {
	if safety >= 1e12 {
		return "infinite"
	}
	return fmt.Sprintf("%.4g", safety)
}

func FormatHeader() string {
	return "principal | tresca | mises | f_t | f_m | first"
}

func TableRow(result Result) string {
	return fmt.Sprintf(
		"[%7.2f,%7.2f,%7.2f] | %7.2f | %7.2f | %5.2f | %5.2f | %s",
		result.Principal[0],
		result.Principal[1],
		result.Principal[2],
		result.TrescaStress,
		result.MisesStress,
		result.TrescaSafety,
		result.MisesSafety,
		result.First,
	)
}
