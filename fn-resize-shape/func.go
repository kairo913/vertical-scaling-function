package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"

	fdk "github.com/fnproject/fdk-go"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/oracle/oci-go-sdk/v65/example/helpers"

	"github.com/caarlos0/env"
)

type Config struct {
	TargetOCPU   int    `env:"TARGET_OCPU"`
	TargetMemory int    `env:"TARGET_MEMORY"`
	TENANT_OCID  string `env:"TENANT_OCID"`
	USER_OCID    string `env:"USER_OCID"`
	REGION       string `env:"REGION"`
	PRIVATE_KEY  string `env:"PRIVATE_KEY"`
	FINGERPRINT  string `env:"FINGERPRINT"`
	PASSPHRASE   string `env:"PASSPHRASE"`
	InstanceId   string `env:"INSTANCE_ID"`
}

func main() {
	fdk.Handle(fdk.HandlerFunc(resizeHandler))
}

func resizeHandler(ctx context.Context, in io.Reader, out io.Writer) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		helpers.FatalIfError(err)
	}

	if cfg.InstanceId == "" {
		helpers.FatalIfError(errors.New("instance ID is required"))
		return
	}

	if cfg.TargetOCPU == 0 || cfg.TargetMemory == 0 {
		helpers.FatalIfError(errors.New("target OCPU and memory are required"))
		return
	}

	// Create a new client
	privateKeyLocation := "/function/" + cfg.PRIVATE_KEY
	privateKey, err := os.ReadFile(privateKeyLocation)
	if err != nil {
		helpers.FatalIfError(err)
	}

	rawConfigProvider := common.NewRawConfigurationProvider(cfg.TENANT_OCID, cfg.USER_OCID, cfg.REGION, cfg.FINGERPRINT, string(privateKey), common.String(cfg.PASSPHRASE))
	client, err := core.NewComputeClientWithConfigurationProvider(rawConfigProvider)
	if err != nil {
		helpers.FatalIfError(err)
	}

	// Resize the instance
	req := core.UpdateInstanceRequest{
		InstanceId: common.String(cfg.InstanceId),
		UpdateInstanceDetails: core.UpdateInstanceDetails{
			ShapeConfig: &core.UpdateInstanceShapeConfigDetails{
				Ocpus:       common.Float32(float32(cfg.TargetOCPU)),
				MemoryInGBs: common.Float32(float32(cfg.TargetMemory)),
			},
		},
	}

	_, err = client.UpdateInstance(context.Background(), req)
	if err != nil {
		helpers.FatalIfError(err)
	}

	json.NewEncoder(out).Encode("Instance resized successfully")
}
