package env

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/exp/slices"

	"github.com/go-leo/config/format"
	"google.golang.org/protobuf/types/known/structpb"
)

type Resource struct {
	prefix string

	formatter format.Formatter
	data      atomic.Value
}

func (r *Resource) Load(ctx context.Context) (*structpb.Struct, error) {
	data, err := r.load(ctx)
	if err != nil {
		return nil, err
	}
	r.data.Store(data)
	return r.formatter.Parse(data)
}

func (r *Resource) load(ctx context.Context) ([]byte, error) {
	var environs [][]byte
	for _, environ := range os.Environ() {
		if strings.HasPrefix(environ, r.prefix) {
			environs = append(environs, []byte(environ))
		}
	}
	if len(environs) <= 0 {
		return nil, fmt.Errorf("config: no environment variables found with prefix %s", r.prefix)
	}
	slices.SortFunc(environs, bytes.Compare)
	return bytes.Join(environs, []byte("\n")), nil
}

func (r *Resource) Watch(ctx context.Context, notifyC chan<- *structpb.Struct, errC chan<- error) (func(ctx context.Context) error, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	stopC := make(chan struct{})
	stop := func(ctx context.Context) error {
		close(stopC)
		return nil
	}
	go func() {
		for {
			fmt.Println("watch")
			select {
			case <-ctx.Done():
				return
			case <-stopC:
				return
			case <-time.After(time.Second):
				data, err := r.load(ctx)
				if err != nil {
					errC <- err
					continue
				}
				preData := r.data.Load()
				if preData != nil && bytes.Equal(preData.([]byte), data) {
					continue
				}
				newValue, err := r.formatter.Parse(data)
				if err != nil {
					errC <- err
					continue
				}
				notifyC <- newValue
				r.data.Store(data)
			}
		}
	}()
	return stop, nil
}

func New(prefix string) (*Resource, error) {
	ext := "env"
	formatter, ok := format.GetFormatter(ext)
	if !ok {
		return nil, fmt.Errorf("config: not found formatter for %s", ext)
	}
	return &Resource{
		prefix:    prefix,
		formatter: formatter,
	}, nil
}
