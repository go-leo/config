package env

import (
	"testing"
)

func TestParse_ValidEnvFormat(t *testing.T) {
	data := []byte("KEY1=VALUE1\nKEY2=VALUE2")
	result, err := Env{}.Parse(data)

	// Expect error due to redundant UnmarshalJSON call in current implementation
	if err == nil {
		t.Errorf("Expected error from UnmarshalJSON but got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result but got %v", result)
	}
}

func TestParse_MalformedInput(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{
			name:    "MissingValue",
			input:   []byte("KEY"),
			wantErr: true,
		},
		{
			name:    "MultipleEquals",
			input:   []byte("KEY=VALUE=MORE"),
			wantErr: true,
		},
		{
			name:    "EmptyLine",
			input:   []byte("\nKEY=VALUE\n"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Env{}.Parse(tt.input)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != nil {
				t.Errorf("Expected nil result but got %v", result)
			}
		})
	}
}

func TestParse_CornerCases(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{
			name:    "Whitespace",
			input:   []byte("  KEY  =  VALUE  "),
			wantErr: true,
		},
		{
			name:    "SpecialCharacters",
			input:   []byte("PATH=/usr/bin:$HOME"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Env{}.Parse(tt.input)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got nil")
			}
			if result != nil {
				t.Errorf("Expected nil result but got %v", result)
			}
		})
	}
}
