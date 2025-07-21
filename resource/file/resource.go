package file

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/go-leo/config/format"
	"google.golang.org/protobuf/types/known/structpb"
)

type Resource struct {
	filename string
	ext      string

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
	return os.ReadFile(r.filename)
}

func (r *Resource) Watch(ctx context.Context, notifyC chan<- *structpb.Struct, errC chan<- error) (func(ctx context.Context) error, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if err := fsWatcher.Add(filepath.Dir(r.filename)); err != nil {
		return nil, err
	}
	stopC := make(chan struct{})
	stop := func(ctx context.Context) error {
		close(stopC)
		return nil
	}
	go func() {
		defer func() {
			if err := fsWatcher.Close(); err != nil {
				errC <- err
			}
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case <-stopC:
				return
			case event, ok := <-fsWatcher.Events:
				if !ok {
					return
				}
				if filepath.Clean(event.Name) != r.filename {
					continue
				}
				if !event.Has(fsnotify.Write) && !event.Has(fsnotify.Create) {
					continue
				}
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
			case err, ok := <-fsWatcher.Errors:
				if !ok {
					return
				}
				errC <- err
			}
		}
	}()
	return stop, nil
}

func New(filename string) (*Resource, error) {
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")
	if ext == "" {
		return nil, fmt.Errorf("config: file extension is empty")
	}
	formatter, ok := format.GetFormatter(ext)
	if !ok {
		return nil, fmt.Errorf("config: not found formatter for %s", ext)
	}
	return &Resource{
		filename:  filename,
		ext:       ext,
		formatter: formatter,
	}, nil
}
