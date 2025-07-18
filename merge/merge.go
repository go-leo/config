package config

import (
	"google.golang.org/protobuf/types/known/structpb"
)

// Merger defines an interface for merging multiple protobuf Structs into one
type Merger interface {
	// Merge combines multiple structpb.Struct values into a single struct
	Merge(values ...*structpb.Struct) *structpb.Struct
}

// merger implements the Merger interface
type merger struct{}

// Merge combines multiple structpb.Struct values into a single struct
// It creates a new target struct and merges all source structs into it
func (m merger) Merge(values ...*structpb.Struct) *structpb.Struct {
	target := ignoreError(structpb.NewStruct(map[string]any{}))
	for _, value := range values {
		m.mergeStruct(target, value)
	}
	return target
}

// mergeStruct merges fields from source struct into target struct
// For each field in source, it makes a deep copy and adds to target
func (m merger) mergeStruct(target *structpb.Struct, source *structpb.Struct) {
	for key, field := range source.GetFields() {
		target.Fields[key] = m.copyValue(field)
	}
}

// mergeList merges values from source ListValue into target ListValue
// Appends each item from source to target after making a copy
func (m merger) mergeList(target *structpb.ListValue, source *structpb.ListValue) {
	for _, item := range source.GetValues() {
		target.Values = append(target.Values, m.copyValue(item))
	}
}

// copyValue creates a deep copy of a protobuf Value
// Handles all types of Value including structs, lists and primitive types
// Returns NullValue for nil or unknown types
func (m merger) copyValue(value *structpb.Value) *structpb.Value {
	if value == nil {
		return structpb.NewNullValue()
	}
	switch v := value.GetKind().(type) {
	case *structpb.Value_NumberValue:
		if v == nil {
			return structpb.NewNullValue()
		}
		return structpb.NewNumberValue(v.NumberValue)
	case *structpb.Value_StringValue:
		if v == nil {
			return structpb.NewNullValue()
		}
		return structpb.NewStringValue(v.StringValue)
	case *structpb.Value_BoolValue:
		if v == nil {
			return structpb.NewNullValue()
		}
		return structpb.NewBoolValue(v.BoolValue)
	case *structpb.Value_StructValue:
		if v == nil {
			return structpb.NewNullValue()
		}
		subValue := ignoreError(structpb.NewStruct(map[string]any{}))
		m.mergeStruct(subValue, v.StructValue)
		return structpb.NewStructValue(subValue)
	case *structpb.Value_ListValue:
		if v == nil {
			return structpb.NewNullValue()
		}
		subList := ignoreError(structpb.NewList([]any{}))
		m.mergeList(subList, v.ListValue)
		return structpb.NewListValue(subList)
	case *structpb.Value_NullValue:
		return structpb.NewNullValue()
	default:
		return structpb.NewNullValue()
	}
}

// ignoreError is a helper function to ignore the error return value
// from functions that return (T, error) when we know error won't occur
func ignoreError[T any](v T, _ error) T {
	return v
}
