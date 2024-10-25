package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"

	fdk "github.com/fnproject/fdk-go"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/common/auth"
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
	PRIVATE_KEY  string `env:"OCI_RESOURCE_PRINCIPAL_PRIVATE_PEM"`
	FINGERPRINT  string `env:"FINGERPRINT"`
	PASSPHRASE   string `env:"PASSPHRASE"`
	InstanceId   string `env:"INSTANCE_ID"`
}

func main() {
	fdk.Handle(fdk.HandlerFunc(resizeHandler))
}

func resizeHandler(ctx context.Context, in io.Reader, out io.Writer) {
	var cfg Config
	err := env.Parse(&cfg)
	helpers.FatalIfError(err)

	log.Printf("Target OCPU: %d", cfg.TargetOCPU)
	log.Printf("Target Memory: %d", cfg.TargetMemory)
	log.Printf("Instance ID: %s", cfg.InstanceId)

	if cfg.InstanceId == "" {
		helpers.FatalIfError(errors.New("instance ID is required"))
	}

	if cfg.TargetOCPU == 0 || cfg.TargetMemory == 0 {
		helpers.FatalIfError(errors.New("target OCPU and memory are required"))
	}

	provider, err := auth.ResourcePrincipalConfigurationProvider()
	helpers.FatalIfError(err)

	client, err := core.NewComputeClientWithConfigurationProvider(provider)
	helpers.FatalIfError(err)

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
