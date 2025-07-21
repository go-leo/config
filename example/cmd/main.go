package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-leo/config/example/configs"
	"github.com/go-leo/config/resource/env"
	"github.com/go-leo/config/resource/file"
)

func main() {
	// prepare config

	os.Setenv("LEO_RUN_ENV", "dev")
	defer os.Unsetenv("LEO_RUN_ENV")

	tmpDir := os.TempDir()

	jsonFilename := tmpDir + "/config.json"
	if err := os.WriteFile(jsonFilename, genConfigJSON(), 0o644); err != nil {
		panic(err)
	}
	defer os.Remove(jsonFilename)

	yamlFilename := tmpDir + "/config.yaml"
	if err := os.WriteFile(yamlFilename, genConfigYaml(), 0o644); err != nil {
		panic(err)
	}
	defer os.Remove(yamlFilename)

	// load config
	envRsc, err := env.New("LEO_")
	if err != nil {
		panic(err)
	}
	jsonRsc, err := file.New(jsonFilename)
	if err != nil {
		panic(err)
	}
	yamlRsc, err := file.New(yamlFilename)
	if err != nil {
		panic(err)
	}
	if err := configs.LoadApplicationConfig(context.TODO(), envRsc, jsonRsc, yamlRsc); err != nil {
		panic(err)
	}
	fmt.Println(configs.GetApplicationConfig())

	sigC, stop, err := configs.WatchApplicationConfig(context.TODO(), envRsc, jsonRsc, yamlRsc)
	if err != nil {
		panic(err)
	}
	go func() {
		time.Sleep(10 * time.Second)
		stop(context.TODO())
	}()

	go func() {
		for range sigC {
			fmt.Println(configs.GetApplicationConfig())
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second)
			if err := os.WriteFile(jsonFilename, genConfigJSON(), 0o644); err != nil {
				panic(err)
			}
			if err := os.WriteFile(yamlFilename, genConfigYaml(), 0o644); err != nil {
				panic(err)
			}
		}
	}()

	time.Sleep(11*time.Second)
}

func genConfigJSON() []byte {
	return []byte(fmt.Sprintf(`{"grpc":{"addr":"127.0.0.1","port":%d}}`, time.Now().Unix()))
}

func genConfigYaml() []byte {
	return []byte(fmt.Sprintf(`
redis:
  addr: 127.0.0.1:6379
  network: tcp
  password: oqnevaqm
  db: %d`, time.Now().Unix()))
}
