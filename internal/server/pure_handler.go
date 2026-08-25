package server

import (
	"net/http"

	"tresca-vm/internal/yield"
)

type pureShearRequest struct {
	Tau    float64 `json:"tau"`
	SigmaY float64 `json:"sigma_y"`
}

func pureShearHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	var req pureShearRequest
	if !readJSON(w, r, &req) {
		return
	}
	result, err := yield.PureShear(req.Tau, req.SigmaY)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}
