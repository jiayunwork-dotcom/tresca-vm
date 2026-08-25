package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func postJSON(t *testing.T, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	Routes().ServeHTTP(rec, req)
	return rec
}

func TestYieldEndpointPrincipal(t *testing.T) {
	rec := postJSON(t, "/api/yield", map[string]interface{}{
		"principal": []float64{100, 0, 0},
		"sigma_y":   100,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body=%s", rec.Code, rec.Body.String())
	}
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["first"] != "both" && body["first"] != "tresca" {
		t.Fatalf("body = %s", rec.Body.String())
	}
}

func TestYieldEndpointTensor(t *testing.T) {
	rec := postJSON(t, "/api/yield", map[string]interface{}{
		"sxx":     100,
		"sigma_y": 200,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body=%s", rec.Code, rec.Body.String())
	}
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["mises_stress"] == nil {
		t.Fatalf("body = %s", rec.Body.String())
	}
}

func TestPureShearEndpoint(t *testing.T) {
	rec := postJSON(t, "/api/pure-shear", map[string]interface{}{
		"tau":     100,
		"sigma_y": 200,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body=%s", rec.Code, rec.Body.String())
	}
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["tresca_yielded"] != true {
		t.Fatalf("body = %s", rec.Body.String())
	}
}

func TestPlaneStressEndpoint(t *testing.T) {
	rec := postJSON(t, "/api/plane-stress", map[string]interface{}{
		"sxx":     120,
		"syy":     -40,
		"txy":     30,
		"sigma_y": 300,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body=%s", rec.Code, rec.Body.String())
	}
}

func TestInvalidSigmaYReturns400(t *testing.T) {
	rec := postJSON(t, "/api/yield", map[string]interface{}{
		"principal": []float64{100, 0, 0},
		"sigma_y":   -1,
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestInvalidTensorReturns400(t *testing.T) {
	rec := postJSON(t, "/api/yield", map[string]interface{}{
		"tensor":  map[string]interface{}{"sxx": 10, "sxy": 3, "syx": 4},
		"sigma_y": 200,
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/yield", nil)
	rec := httptest.NewRecorder()
	Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
}
