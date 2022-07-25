package json

import (
	"reflect"
	"testing"
)

// TestParse_Success tests successful parsing of valid JSON data.
func TestParse_Success(t *testing.T) {
	data := []byte(`{"name":"Alice","age":30,"isStudent":false}`)
	parser := Json{}
	result, err := parser.Parse(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatalf("Expected non-nil result, got nil")
	}

	expectedMap := map[string]interface{}{
		"name":      "Alice",
		"age":       float64(30),
		"isStudent": false,
	}

	if !reflect.DeepEqual(expectedMap, result.AsMap()) {
		t.Errorf("Expected map %v, got %v", expectedMap, result.AsMap())
	}
}

// TestParse_InvalidJson tests error handling for invalid JSON format.
func TestParse_InvalidJson(t *testing.T) {
	data := []byte(`{name:"Alice", age: 30}`) // Missing quotes around key
	parser := Json{}
	result, err := parser.Parse(data)

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}
