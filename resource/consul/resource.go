package consul

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-leo/config"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/hashicorp/go-hclog"
)

var _ config.Resource = (*Resource)(nil)

type Resource struct {
	Formatter config.Formatter
	Client    *api.Client
	Key       string
}

func (r *Resource) Format() string {
	if r.Formatter == nil {
		return strings.TrimPrefix(filepath.Ext(r.Key), ".")
	}
	return r.Formatter.Format()
}

func (r *Resource) Load(ctx context.Context) ([]byte, error) {
	pair, _, err := r.Client.KV().Get(r.Key, nil)
	if err != nil {
		return nil, err
	}
	return pair.Value, nil
}

func (r *Resource) Watch(ctx context.Context, notifyC chan<- *config.Event) error {
	params := map[string]any{
		"type": "key",
		"key":  r.Key,
	}
	plan, err := watch.Parse(params)
	if err != nil {
		return err
	}
	plan.Handler = func(idx uint64, raw interface{}) {
		if raw == nil {
			return
		}
		if pair, ok := raw.(*api.KVPair); ok {
			notifyC <- config.NewDataEvent(pair.Value)
		}
	}
	go func() {
		err = plan.RunWithClientAndHclog(r.Client, &consuleLogger{Logger: hclog.NewNullLogger(), notifyC: notifyC})
		if err != nil {
			notifyC <- config.NewErrorEvent(err)
		}
		notifyC <- config.NewErrorEvent(config.ErrStopWatch)
	}()
	go func() {
		<-ctx.Done()
		plan.Stop()
	}()
	return nil
}

type consuleLogger struct {
	hclog.Logger
	notifyC chan<- *config.Event
}

func (l *consuleLogger) Error(msg string, args ...interface{}) {
	l.notifyC <- config.NewErrorEvent(fmt.Errorf(msg, args...))
}
