package opportunity

import "testing"

func TestInputValidate(t *testing.T) {
	valid := Input{Title: "Scholarship", Description: "For students", Types: []Type{Scholarship}, Link: "https://example.com"}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid input, got %v", err)
	}
	multiType := Input{Title: "Research event", Description: "For students", Types: []Type{Research, Extras}}
	if err := multiType.Validate(); err != nil {
		t.Fatalf("expected multi-type input to be valid, got %v", err)
	}
	for _, input := range []Input{{Title: "", Description: "x", Types: []Type{Scholarship}}, {Title: "x", Description: "", Types: []Type{Scholarship}}, {Title: "x", Description: "x", Types: []Type{"bad"}}, {Title: "x", Description: "x", Types: []Type{Scholarship, Scholarship}}, {Title: "x", Description: "x", Types: []Type{Scholarship}, Link: "not-a-url"}} {
		if err := input.Validate(); err == nil {
			t.Fatalf("expected validation failure for %#v", input)
		}
	}
}

func TestTypeValid(t *testing.T) {
	for _, value := range []Type{Scholarship, Hackathon, Internship, Research, Extras} {
		if !value.Valid() {
			t.Errorf("%q should be a valid opportunity type", value)
		}
	}
	if Type("research_extra").Valid() {
		t.Fatal("legacy research_extra should no longer be a valid opportunity type")
	}
}
