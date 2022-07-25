package consul

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/go-leo/config"
	"github.com/hashicorp/consul/api"
)

func consulFactory() (*api.Client, error) {
	return api.NewClient(api.DefaultConfig())
}

func TestResource_Load_Consul(t *testing.T) {
	client, err := consulFactory()
	if err != nil {
		t.Errorf("factory() error = %v", err)
		return
	}

	_, err = client.KV().Put(&api.KVPair{
		Key:   "consul",
		Value: []byte("TEST_KEY=test_value"),
	}, nil)
	if err != nil {
		t.Errorf("Put() error = %v", err)
		return
	}

	defer func() {
		_, err = client.KV().Delete("consul", nil)
		if err != nil {
			t.Errorf("Delete() error = %v", err)
			return
		}
	}()

	time.Sleep(time.Second)

	r := Resource{
		Formatter: config.Env{},
		Client:    client,
		Key:       "consul",
	}
	ctx := context.Background()
	content, err := r.Load(ctx)
	if err != nil {
		t.Errorf("Load() error = %v", err)
		return
	}

	if !strings.Contains(string(content), "TEST_KEY=test_value") {
		t.Errorf("Load() data = %v, want data to contain 'TEST_KEY=test_value'", string(content))
	}
	time.Sleep(time.Second)
}

func TestResource_Watch_Consul(t *testing.T) {
	client, err := consulFactory()
	if err != nil {
		t.Errorf("factory() error = %v", err)
		return
	}

	_, err = client.KV().Put(&api.KVPair{
		Key:   "consul",
		Value: []byte("TEST_KEY=test_value"),
	}, nil)
	if err != nil {
		t.Errorf("PublishConfig() error = %v", err)
		return
	}

	defer func() {
		_, err = client.KV().Delete("consul", nil)
		if err != nil {
			t.Errorf("PublishConfig() error = %v", err)
			return
		}
	}()

	time.Sleep(time.Second)

	r := Resource{
		Formatter: config.Env{},
		Client:    client,
		Key:       "consul",
	}

	notifyC := make(chan *config.Event, 1)
	// Start watching
	ctx := context.Background()
	stopFunc, err := r.Watch(ctx, notifyC)
	if err != nil {
		t.Errorf("Watch() error = %v", err)
		return
	}

	// Give some time for the watcher to detect the change
	go func() {
		for {
			time.Sleep(time.Second)
			_, err = client.KV().Put(&api.KVPair{
				Key:   "nacos",
				Value: []byte("TEST_KEY_NEW=test_value_new" + time.Now().Format(time.RFC3339)),
			}, nil)
			if err != nil {
				t.Errorf("PublishConfig() error = %v", err)
				return
			}
		}
	}()

	// Wait for the event
	select {
	case event := <-notifyC:
		if data, ok := event.AsDataEvent(); !ok || data.Data == nil {
			t.Error("Expected DataEvent with non-nil data")
		}
	case <-time.After(100 * time.Second):
		t.Error("No event received within the timeout")
	}

	stopFunc()

	_, err = client.KV().Put(&api.KVPair{
		Key:   "nacos",
		Value: []byte("TEST_KEY=another_test_value"),
	}, nil)
	if err != nil {
		t.Errorf("PublishConfig() error = %v", err)
		return
	}

	select {
	case event := <-notifyC:
		err, ok := event.AsErrorEvent()
		if !ok || err.Err == nil || !errors.Is(err.Err, config.ErrStopWatch) {
			t.Error("Did not expect to receive an event after stopping the watcher")
		}
	case <-time.After(100 * time.Millisecond):
		// Expected behavior
	}
}
