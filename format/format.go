package format

import (
	"strings"
	"sync"

	"google.golang.org/protobuf/types/known/structpb"
)

// Global formatters registry mapping file extensions to their corresponding parsers
var (
	formatters   = make(map[string]Formatter) // Stores registered format parsers
	parsersMutex sync.RWMutex                 // Mutex to protect concurrent access to formatters
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
	parsersMutex.Lock()
	formatters[strings.ToLower(ext)] = formatter
	parsersMutex.Unlock()
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
	var formatter Formatter
	var ok bool
	parsersMutex.RLock()
	formatter, ok = formatters[strings.ToLower(ext)]
	parsersMutex.RUnlock()
	return formatter, ok
}
