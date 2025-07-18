package config

import (
	"testing"

	"google.golang.org/protobuf/types/known/structpb"
)

func TestMergeStructs(t *testing.T) {
	m := merger{}

	// Create test structs
	s1, err := structpb.NewStruct(map[string]interface{}{
		"name": "John",
		"age":  30,
		"address": map[string]interface{}{
			"city": "New York",
		},
	})
	if err != nil {
		t.Fatalf("Failed to create test struct: %v", err)
	}

	s2, err := structpb.NewStruct(map[string]interface{}{
		"age": 31,
		"address": map[string]interface{}{
			"zip": "10001",
		},
		"email": "john@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create test struct: %v", err)
	}

	// Merge the structs
	result := m.Merge(s1, s2)

	// Verify results
	if result.GetFields()["name"].GetStringValue() != "John" {
		t.Errorf("Expected name 'John', got %v", result.GetFields()["name"])
	}

	if result.GetFields()["age"].GetNumberValue() != 31 {
		t.Errorf("Expected age 31, got %v", result.GetFields()["age"])
	}

	address := result.GetFields()["address"].GetStructValue()
	if address.GetFields()["city"].GetStringValue() != "" {
		t.Errorf("Expected city '', got %v", address.GetFields()["city"])
	}
	if address.GetFields()["zip"].GetStringValue() != "10001" {
		t.Errorf("Expected zip '10001', got %v", address.GetFields()["zip"])
	}

	if result.GetFields()["email"].GetStringValue() != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got %v", result.GetFields()["email"])
	}
}

func TestMergeWithEmptyStruct(t *testing.T) {
	m := merger{}

	s1, err := structpb.NewStruct(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	s2, err := structpb.NewStruct(map[string]interface{}{
		"key": "value",
	})
	if err != nil {
		t.Fatal(err)
	}

	result := m.Merge(s1, s2)

	if len(result.GetFields()) != 1 {
		t.Errorf("Expected 1 field, got %d", len(result.GetFields()))
	}
	if result.GetFields()["key"].GetStringValue() != "value" {
		t.Errorf("Expected value 'value', got %v", result.GetFields()["key"])
	}
}

func TestMergeLists(t *testing.T) {
	m := merger{}

	list1, err := structpb.NewList([]interface{}{1, 2})
	if err != nil {
		t.Fatal(err)
	}

	list2, err := structpb.NewList([]interface{}{3, 4})
	if err != nil {
		t.Fatal(err)
	}

	s1 := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"s1_numbers": structpb.NewListValue(list1),
		},
	}

	s2 := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"s2_numbers": structpb.NewListValue(list2),
		},
	}

	result := m.Merge(s1, s2)

	mergedList := result.GetFields()["s2_numbers"].GetListValue()
	if len(mergedList.GetValues()) != 2 {
		t.Errorf("Expected list length 2, got %d", len(mergedList.GetValues()))
	}
}

func TestCopyValue(t *testing.T) {
	m := merger{}

	tests := []struct {
		name  string
		value *structpb.Value
		check func(*structpb.Value) bool
	}{
		{
			name:  "number",
			value: structpb.NewNumberValue(42),
			check: func(v *structpb.Value) bool {
				return v.GetNumberValue() == 42
			},
		},
		{
			name:  "string",
			value: structpb.NewStringValue("test"),
			check: func(v *structpb.Value) bool {
				return v.GetStringValue() == "test"
			},
		},
		{
			name:  "bool",
			value: structpb.NewBoolValue(true),
			check: func(v *structpb.Value) bool {
				return v.GetBoolValue() == true
			},
		},
		{
			name: "nested struct",
			value: structpb.NewStructValue(&structpb.Struct{
				Fields: map[string]*structpb.Value{
					"key": structpb.NewStringValue("value"),
				},
			}),
			check: func(v *structpb.Value) bool {
				return v.GetStructValue().GetFields()["key"].GetStringValue() == "value"
			},
		},
		{
			name: "null",
			value: structpb.NewNullValue(),
			check: func(v *structpb.Value) bool {
				return v.GetKind() == nil || v.GetNullValue() == 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copied := m.copyValue(tt.value)
			if !tt.check(copied) {
				t.Errorf("Copy failed for %s", tt.name)
			}
		})
	}
}