package consul

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/go-leo/config/format"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/hashicorp/go-hclog"
	
	"google.golang.org/protobuf/types/known/structpb"
)

type Resource struct {
	client *api.Client
	key    string
	ext    string

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
	pair, _, err := r.client.KV().Get(r.key, nil)
	if err != nil {
		return nil, err
	}
	return pair.Value, nil
}

func (r *Resource) Watch(ctx context.Context, notifyC chan<- *structpb.Struct, errC chan<- error) (func(ctx context.Context) error, error) {
	params := map[string]any{
		"type": "key",
		"key":  r.key,
	}
	plan, err := watch.Parse(params)
	if err != nil {
		return nil, err
	}
	plan.Handler = func(idx uint64, raw interface{}) {
		if raw == nil {
			return
		}
		pair, ok := raw.(*api.KVPair)
		if !ok {
			return
		}
		data := pair.Value
		preData := r.data.Load()
		if preData != nil && bytes.Equal(preData.([]byte), data) {
			return
		}
		newValue, err := r.formatter.Parse(data)
		if err != nil {
			errC <- err
			return
		}
		notifyC <- newValue
		r.data.Store(data)
	}
	go func() {
		_ = plan.RunWithClientAndHclog(
			r.client,
			&consuleLogger{
				Logger: hclog.NewNullLogger(),
				errC:   errC,
			})
	}()
	stop := func(ctx context.Context) error {
		plan.Stop()
		return nil
	}
	return stop, nil
}

type consuleLogger struct {
	hclog.Logger
	errC chan<- error
}

func (l *consuleLogger) Error(msg string, args ...interface{}) {
	l.errC <- fmt.Errorf(msg, args...)
}

func New(client *api.Client, key string) (*Resource, error) {
	ext := strings.TrimPrefix(filepath.Ext(key), ".")
	if ext == "" {
		return nil, fmt.Errorf("config: key extension is empty")
	}
	formatter, ok := format.GetFormatter(ext)
	if !ok {
		return nil, fmt.Errorf("config: not found formatter for %s", ext)
	}
	return &Resource{
		client:    client,
		key:       key,
		ext:       ext,
		formatter: formatter,
	}, nil
}
