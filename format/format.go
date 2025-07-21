package format

import (
	"strings"
	"sync"

	"google.golang.org/protobuf/types/known/structpb"
)

// Global formatters registry mapping file extensions to their corresponding parsers
var (
	// formatters stores registered format parsers
	formatters = make(map[string]Formatter)
	// mutex to protect concurrent access to formatters
	mutex sync.RWMutex
)

// Formatter interface defines the standard method for parsing configuration data
type Formatter interface {
	// Parse converts byte data into a protobuf Struct object
	//
	// Args:
	//   data ([]byte): Raw configuration data
	//
	// Returns:
	//   *structpb.Struct: Parsed structured data
	//   error: Error if parsing fails
	Parse(data []byte) (*structpb.Struct, error)
}

// RegisterFormatter associates a file extension with a configuration parser
//
// Args:
//
//	ext (string): File extension (e.g., "yaml", "toml")
//	formatter (Formatter): Implementation of the Formatter interface
func RegisterFormatter(ext string, formatter Formatter) {
	mutex.Lock()
	formatters[strings.ToLower(ext)] = formatter
	mutex.Unlock()
}

// GetFormatter retrieves the parser associated with a specific file extension
//
// Args:
//
//	ext (string): File extension to look up
//
// Returns:
//
//	Formatter: Registered parser or nil if not found
func GetFormatter(ext string) (Formatter, bool) {
	mutex.RLock()
	formatter, ok := formatters[strings.ToLower(ext)]
	mutex.RUnlock()
	return formatter, ok
}
