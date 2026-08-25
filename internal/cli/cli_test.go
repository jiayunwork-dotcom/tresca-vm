package cli

import (
	"testing"

	"tresca-vm/internal/yield"
)

func TestParseList(t *testing.T) {
	values, err := parseList("100,0,-50")
	if err != nil {
		t.Fatal(err)
	}
	if len(values) != 3 || values[0] != 100 || values[2] != -50 {
		t.Fatalf("values = %v", values)
	}
	if _, err := parseList("1,2"); err == nil {
		t.Fatal("parseList accepted two values")
	}
}

func TestValidateSigmaY(t *testing.T) {
	if err := validateSigmaY(200); err != nil {
		t.Fatal(err)
	}
	if err := validateSigmaY(0); err == nil {
		t.Fatal("validateSigmaY accepted zero")
	}
}

func TestLoadExample(t *testing.T) {
	caseValue, err := loadExample("../../example/uniaxial.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(caseValue.Principal) != 3 {
		t.Fatalf("principal = %v", caseValue.Principal)
	}
}

func TestExampleCases(t *testing.T) {
	cases := ExampleCases()
	if len(cases) != 3 {
		t.Fatalf("cases = %d", len(cases))
	}
	for _, c := range cases {
		result, err := yield.Evaluate(c.Principal, c.SigmaY)
		if err != nil {
			t.Fatal(err)
		}
		if result.TrescaStress < 0 || result.MisesStress < 0 {
			t.Fatalf("negative equivalent stress: %+v", result)
		}
	}
}

func TestTensorSymmetric(t *testing.T) {
	tensor := tensorSymmetric(1, 2, 3, 4, 5, 6)
	if tensor.XX != 1 || tensor.YY != 2 || tensor.ZZ != 3 {
		t.Fatalf("tensor = %+v", tensor)
	}
}
