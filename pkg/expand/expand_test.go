package expand

import (
	"os"
	"testing"
)

func TestExpandEnvStrings_ExpandsNestedStructAndSlice(t *testing.T) {
	t.Setenv("EXPAND_TEST_VALUE", "world")

	type child struct {
		Message string
	}
	type sample struct {
		Title  string
		Child  child
		Values []string
	}

	v := sample{
		Title: "hello-${EXPAND_TEST_VALUE}",
		Child: child{Message: "child-${EXPAND_TEST_VALUE}"},
		Values: []string{
			"a-${EXPAND_TEST_VALUE}",
			"b-${EXPAND_TEST_VALUE}",
		},
	}

	ExpandEnvStrings(&v)

	if v.Title != "hello-world" {
		t.Fatalf("unexpected title: %q", v.Title)
	}
	if v.Child.Message != "child-world" {
		t.Fatalf("unexpected child message: %q", v.Child.Message)
	}
	if v.Values[0] != "a-world" || v.Values[1] != "b-world" {
		t.Fatalf("unexpected values: %#v", v.Values)
	}
}

func TestExpandEnvStrings_ExpandsMapStringValues(t *testing.T) {
	t.Setenv("EXPAND_TEST_MAP", "ok")

	v := map[string]string{
		"first":  "${EXPAND_TEST_MAP}",
		"second": "value-${EXPAND_TEST_MAP}",
	}

	ExpandEnvStrings(&v)

	if v["first"] != "ok" {
		t.Fatalf("unexpected first value: %q", v["first"])
	}
	if v["second"] != "value-ok" {
		t.Fatalf("unexpected second value: %q", v["second"])
	}
}

func TestExpandEnvStrings_UsesCurrentEnvironment(t *testing.T) {
	_ = os.Setenv("EXPAND_DIRECT", "42")
	t.Cleanup(func() { _ = os.Unsetenv("EXPAND_DIRECT") })

	v := "value-${EXPAND_DIRECT}"
	ExpandEnvStrings(&v)

	if v != "value-42" {
		t.Fatalf("unexpected value: %q", v)
	}
}

func TestExpandEnvStrings_ExpandsMapWithAnyValues(t *testing.T) {
	t.Setenv("EXPAND_MAP_ANY", "resolved")

	v := map[string]any{
		"str": "${EXPAND_MAP_ANY}",
		"num": 42,
	}

	ExpandEnvStrings(&v)

	if v["num"] != 42 {
		t.Fatalf("unexpected num: %v", v["num"])
	}
}

func TestExpandEnvStrings_NoExpansionOnNonStringSlice(t *testing.T) {
	t.Setenv("EXPAND_NUM", "99")

	v := []int{1, 2, 3}
	ExpandEnvStrings(&v)

	if v[0] != 1 || v[1] != 2 || v[2] != 3 {
		t.Fatalf("unexpected slice: %v", v)
	}
}
