package gflatten_test

import (
	"encoding/json"
	"testing"

	"github.com/cloudmatelabs/gflatten"
)

var src = map[string]any{
	"foo": []any{
		"bar", "baz",
	},
	"foobar": map[string]any{
		"foo": []any{
			"baz",
			map[string]any{
				"bar": map[string]any{
					"baz": "foobar",
				},
			},
		},
	},
}

func TestFlatten(t *testing.T) {
	option := gflatten.Option{
		ParameterDelimiter: ".",
		ArrayWrap:          gflatten.WRAP.SQUARE_BRACKET,
	}
	flattenTest(t, src, option, map[string]any{
		"foo[0]":                "bar",
		"foo[1]":                "baz",
		"foobar.foo[0]":         "baz",
		"foobar.foo[1].bar.baz": "foobar",
	})
}

func TestFlattenPostgresStyle(t *testing.T) {
	option := gflatten.Option{
		ParameterDelimiter: "->",
		ArrayDelimiter:     "->",
		ParameterWrap:      gflatten.WRAP.SINGLE_QUOTE,
	}
	flattenTest(t, src, option, map[string]any{
		"'foo'->0":                         "bar",
		"'foo'->1":                         "baz",
		"'foobar'->'foo'->0":               "baz",
		"'foobar'->'foo'->1->'bar'->'baz'": "foobar",
	})
}

func TestFlattenMySQLStyle(t *testing.T) {
	option := gflatten.Option{
		Prefix:             "$",
		ParameterDelimiter: ".",
		ArrayWrap:          gflatten.WRAP.SQUARE_BRACKET,
	}
	flattenTest(t, src, option, map[string]any{
		"$.foo[0]":                "bar",
		"$.foo[1]":                "baz",
		"$.foobar.foo[0]":         "baz",
		"$.foobar.foo[1].bar.baz": "foobar",
	})
}

func TestFlattenStructInput(t *testing.T) {
	type bar struct {
		Baz string `json:"baz"`
	}
	type foobar struct {
		Foo []bar `json:"foo"`
	}
	type input struct {
		Foo    []string `json:"foo"`
		Foobar foobar   `json:"foobar"`
	}

	src := input{
		Foo: []string{"bar", "baz"},
		Foobar: foobar{
			Foo: []bar{
				{Baz: "baz"},
				{Baz: "foobar"},
			},
		},
	}
	option := gflatten.Option{
		ParameterDelimiter: ".",
		ArrayWrap:          gflatten.WRAP.SQUARE_BRACKET,
	}
	flattenTest(t, src, option, map[string]any{
		"foo[0]":            "bar",
		"foo[1]":            "baz",
		"foobar.foo[0].baz": "baz",
		"foobar.foo[1].baz": "foobar",
	})
}

func flattenTest(t *testing.T, src any, option gflatten.Option, expect map[string]any) {
	dest, err := gflatten.Flatten(src, option)
	if err != nil {
		t.Fatalf("Error: %v", err)
		return
	}

	fatal := false

	for k, v := range expect {
		if dest[k] != v {
			fatal = true
			t.Errorf("dest[%s] != %v", k, v)
		}
	}

	if fatal {
		d, _ := json.MarshalIndent(dest, "", "  ")
		t.Fatal(string(d))
	}
}
