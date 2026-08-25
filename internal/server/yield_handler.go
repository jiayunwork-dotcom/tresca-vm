package server

import (
	"net/http"

	"tresca-vm/internal/tensor"
	"tresca-vm/internal/yield"
)

type yieldRequest struct {
	Principal []float64     `json:"principal"`
	SigmaY    float64       `json:"sigma_y"`
	Tensor    *tensor.Input `json:"tensor"`
	XX        *float64      `json:"sxx"`
	YY        *float64      `json:"syy"`
	ZZ        *float64      `json:"szz"`
	XY        *float64      `json:"sxy"`
	XZ        *float64      `json:"sxz"`
	YZ        *float64      `json:"syz"`
}

func yieldHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	var req yieldRequest
	if !readJSON(w, r, &req) {
		return
	}
	var result yield.Result
	var err error
	switch {
	case len(req.Principal) > 0:
		result, err = yield.Evaluate(req.Principal, req.SigmaY)
	case req.Tensor != nil:
		t, parseErr := tensor.FromInput(*req.Tensor, true)
		if parseErr != nil {
			badRequest(w, parseErr.Error())
			return
		}
		result, err = yield.EvaluateTensor(t, req.SigmaY)
	case hasComponent(req.XX) || hasComponent(req.YY) || hasComponent(req.ZZ) ||
		hasComponent(req.XY) || hasComponent(req.XZ) || hasComponent(req.YZ):
		in := tensor.Input{XX: req.XX, YY: req.YY, ZZ: req.ZZ, XY: req.XY, XZ: req.XZ, YZ: req.YZ}
		t, parseErr := tensor.FromInput(in, true)
		if parseErr != nil {
			badRequest(w, parseErr.Error())
			return
		}
		result, err = yield.EvaluateTensor(t, req.SigmaY)
	default:
		badRequest(w, "provide principal or tensor components")
		return
	}
	if err != nil {
		writeValidationOutcome(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func hasComponent(value *float64) bool {
	return value != nil
}

func tensorFromRequest(req yieldRequest) (tensor.Tensor, error) {
	if req.Tensor != nil {
		return tensor.FromInput(*req.Tensor, true)
	}
	in := tensor.Input{XX: req.XX, YY: req.YY, ZZ: req.ZZ, XY: req.XY, XZ: req.XZ, YZ: req.YZ}
	return tensor.FromInput(in, true)
}
