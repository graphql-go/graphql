package graphql

import (
	"testing"
)

func TestDefaultResolveFn(t *testing.T) {
	p := ResolveParams{
		Source: &struct {
			A string `json:"a"`
			B string `json:"b"`
			C string `json:"c"`
			D string `json:"d"`
			E string `json:"e"`
			F string `json:"f"`
			G string `json:"g"`
			H string `json:"h"`
		}{
			F: "testing",
		},
		Info: ResolveInfo{
			FieldName: "F",
		},
	}
	v, err := defaultResolveFn(p)
	if err != nil {
		t.Fatal(err)
	}
	if s, ok := v.(string); !ok {
		t.Fatalf("Expected string, got %T", v)
	} else if s != "testing" {
		t.Fatalf("Expected 'testing'")
	}

	p = ResolveParams{
		Source: map[string]interface{}{
			"A": "a",
			"B": "b",
			"C": "c",
			"D": "d",
			"E": "e",
			"F": "testing",
			"G": func() interface{} { return "g" },
			"H": "h",
		},
		Info: ResolveInfo{
			FieldName: "F",
		},
	}
	v, err = defaultResolveFn(p)
	if err != nil {
		t.Fatal(err)
	}
	if s, ok := v.(string); !ok {
		t.Fatalf("Expected string, got %T", v)
	} else if s != "testing" {
		t.Fatalf("Expected 'testing'")
	}

	p.Info.FieldName = "G"
	v, err = defaultResolveFn(p)
	if err != nil {
		t.Fatal(err)
	}
	if s, ok := v.(string); !ok {
		t.Fatalf("Expected string, got %T", v)
	} else if s != "g" {
		t.Fatalf("Expected 'testing'")
	}
}

func BenchmarkDefaultResolveFnStruct(b *testing.B) {
	p := ResolveParams{
		Source: &struct {
			A string `json:"a"`
			B string `json:"b"`
			C string `json:"c"`
			D string `json:"d"`
			E string `json:"e"`
			F string `json:"f"`
			G string `json:"g"`
			H string `json:"h"`
		}{
			F: "testing",
		},
		Info: ResolveInfo{
			FieldName: "F",
		},
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		defaultResolveFn(p)
	}
}

func BenchmarkDefaultResolveFnMap(b *testing.B) {
	p := ResolveParams{
		Source: map[string]interface{}{
			"A": "a",
			"B": "b",
			"C": "c",
			"D": "d",
			"E": "e",
			"F": "testing",
			"G": "g",
			"H": "h",
		},
		Info: ResolveInfo{
			FieldName: "F",
		},
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		defaultResolveFn(p)
	}
}
