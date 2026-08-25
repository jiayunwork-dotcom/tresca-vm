package cli

import (
	"fmt"
	"os"
)

func Run(args []string) int {
	if len(args) == 0 {
		return runServe([]string{})
	}
	switch args[0] {
	case "yield":
		return runYield(args[1:])
	case "tensor":
		return runTensor(args[1:])
	case "pure-shear":
		return runPureShear(args[1:])
	case "plane":
		return runPlane(args[1:])
	case "example":
		return runExample(args[1:])
	case "serve":
		return runServe(args[1:])
	case "help", "-h", "--help":
		printHelp()
		return 0
	case "version":
		fmt.Println("tresca-vm 1.0.0")
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", args[0])
		printHelp()
		return 2
	}
}

func printHelp() {
	fmt.Println(`tresca-vm: Tresca and von Mises yield calculator

Usage:
  tresca-vm                         start HTTP server on :8080
  tresca-vm yield -principal 100,0,-50 -sigmaY 200
  tresca-vm tensor -sxx 100 -syy -50 -sigmaY 200
  tresca-vm pure-shear -tau 100 -sigmaY 200
  tresca-vm plane -sxx 100 -syy -50 -txy 30 -sigmaY 200
  tresca-vm example -file example/uniaxial.json
  tresca-vm serve -addr :8080

HTTP:
  POST /api/yield       {"principal":[100,0,-50],"sigma_y":200}
  POST /api/pure-shear  {"tau":100,"sigma_y":200}
  POST /api/plane-stress {"sxx":100,"syy":-50,"txy":30,"sigma_y":200}`)
}

func fail(err error) int {
	fmt.Fprintln(os.Stderr, "error:", err)
	return 1
}
