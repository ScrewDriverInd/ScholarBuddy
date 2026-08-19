package opportunity

import "testing"

func TestInputValidate(t *testing.T) {
	valid := Input{Title: "Scholarship", Description: "For students", Type: Scholarship, Link: "https://example.com"}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid input, got %v", err)
	}
	for _, input := range []Input{{Title: "", Description: "x", Type: Scholarship}, {Title: "x", Description: "", Type: Scholarship}, {Title: "x", Description: "x", Type: "bad"}, {Title: "x", Description: "x", Type: Scholarship, Link: "not-a-url"}} {
		if err := input.Validate(); err == nil {
			t.Fatalf("expected validation failure for %#v", input)
		}
	}
}
