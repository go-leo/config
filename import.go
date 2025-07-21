package config

import (
	// Environment variable format support
	// Automatically registers env format decoder when imported
	_ "github.com/go-leo/config/format/env"

	// JSON format support
	// Automatically registers json format decoder when imported
	_ "github.com/go-leo/config/format/json"

	// TOML format support
	// Automatically registers toml format decoder when imported
	_ "github.com/go-leo/config/format/toml"

	// YAML format support
	// Automatically registers yaml format decoder when imported
	_ "github.com/go-leo/config/format/yaml"

	// Sample merger implementation
	// Automatically registers sample merger when imported
	_ "github.com/go-leo/config/merge/sample"
)
