package parser

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	b, err := os.ReadFile("testdata/success/http.idl")
	if err != nil {
		t.Fatal(err)
	}
	doc, err := Parse(string(b))
	if err != nil {
		t.Fatal(err)
	}
	b, err = os.ReadFile("testdata/success/http.formated.idl")
	if err != nil {
		t.Fatal(err)
	}
	s := Dump(doc)
	if s != string(b) {
		t.Fatalf("expected:\n%s\nbut got:\n%s", string(b), s)
	}
	b, err = os.ReadFile("testdata/success/http.idl.json")
	if err != nil {
		t.Fatal(err)
	}
	v, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	b = bytes.TrimSpace(b)
	if bytes.Compare(v, b) != 0 {
		t.Fatalf("expected:\n%s\nbut got:\n%s", string(b), string(v))
	}
}
