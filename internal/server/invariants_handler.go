package server

import (
	"net/http"

	"tresca-vm/internal/tensor"
)

type invariantsRequest struct {
	Tensor *tensor.Input `json:"tensor"`
	XX     *float64      `json:"sxx"`
	YY     *float64      `json:"syy"`
	ZZ     *float64      `json:"szz"`
	XY     *float64      `json:"sxy"`
	XZ     *float64      `json:"sxz"`
	YZ     *float64      `json:"syz"`
}

func invariantsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	var req invariantsRequest
	if !readJSON(w, r, &req) {
		return
	}
	t, err := invariantsTensor(req)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	inv := tensor.ComputeInvariants(t)
	principal, err := tensor.PrincipalStresses(t)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"invariants": inv,
		"principal":  principal,
		"deviator":   t.Deviator(),
	})
}

func invariantsTensor(req invariantsRequest) (tensor.Tensor, error) {
	if req.Tensor != nil {
		return tensor.FromInput(*req.Tensor, true)
	}
	in := tensor.Input{XX: req.XX, YY: req.YY, ZZ: req.ZZ, XY: req.XY, XZ: req.XZ, YZ: req.YZ}
	return tensor.FromInput(in, true)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"name":    "tresca-vm",
		"version": "1.0.0",
	})
}
