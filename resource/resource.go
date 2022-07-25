package resource

import (
	"context"

	"google.golang.org/protobuf/types/known/structpb"
)

type Resource interface {
	Load(ctx context.Context) (*structpb.Struct, error)
	Watch(ctx context.Context, notifyC chan<- *structpb.Struct) (func(), error)
}
