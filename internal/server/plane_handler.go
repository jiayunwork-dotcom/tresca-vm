package server

import (
	"net/http"

	"tresca-vm/internal/yield"
)

type planeStressRequest struct {
	SXX    float64 `json:"sxx"`
	SYY    float64 `json:"syy"`
	TXY    float64 `json:"txy"`
	SigmaY float64 `json:"sigma_y"`
}

func planeStressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	var req planeStressRequest
	if !readJSON(w, r, &req) {
		return
	}
	result, err := yield.PlaneStress(req.SXX, req.SYY, req.TXY, req.SigmaY)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}
