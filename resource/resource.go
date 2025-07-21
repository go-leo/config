package resource

import (
	"context"

	"google.golang.org/protobuf/types/known/structpb"
)

// Resource defines the core interface for configuration resource providers.
// Implementations should handle both synchronous loading and change monitoring
// of configuration data.
type Resource interface {
	// Load retrieves the current configuration state.
	// Args:
	//   - ctx: Context for cancellation and timeouts
	// Returns:
	//   - *structpb.Struct: Current configuration data in protobuf Struct format
	//   - error: Loading error if any

	Load(ctx context.Context) (*structpb.Struct, error)

	// Watch establishes a continuous monitoring of configuration changes.
	// Args:
	//   - ctx: Context for cancellation
	//   - notifyC: Channel for receiving configuration updates (send-only)
	//   - errC: Channel for receiving monitoring errors (send-only)
	// Returns:
	//   - func(context.Context) error: Cleanup function that stops watching,
	//     takes context for graceful shutdown, returns any cleanup error
	//   - error: Immediate error if watch setup fails
	Watch(ctx context.Context, notifyC chan<- *structpb.Struct, errC chan<- error) (func(context.Context) error, error)
}
