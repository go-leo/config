package toml

import (
	"reflect"
	"testing"
)

// TestParse_Success tests successful parsing of valid TOML data.
func TestParse_Success(t *testing.T) {
	data := []byte("name = 'Alice'\nage = 30")
	parser := Toml{}
	result, err := parser.Parse(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatalf("Expected non-nil result, got nil")
	}

	expectedMap := map[string]interface{}{
		"name": "Alice",
		"age":  float64(30),
	}

	if !reflect.DeepEqual(expectedMap, result.AsMap()) {
		t.Errorf("Expected map %v, got %v", expectedMap, result.AsMap())
	}
}

// TestParse_InvalidToml tests error handling for invalid TOML format.
func TestParse_InvalidToml(t *testing.T) {
	data := []byte("key = invalid_value") // Invalid TOML (no quotes around string)
	parser := Toml{}
	result, err := parser.Parse(data)

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}
