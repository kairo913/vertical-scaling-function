package main

import (
	"context"
	"encoding/json"
	"io"
	"log"

	fdk "github.com/fnproject/fdk-go"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/common/auth"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/oracle/oci-go-sdk/v65/example/helpers"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	ActiveOCPU     int    `env:"ACTIVE_OCPU, required"`
	ActiveMemory   int    `env:"ACTIVE_MEMORY, required"`
	InActiveOCPU   int    `env:"INACTIVE_OCPU, required"`
	InActiveMemory int    `env:"INACTIVE_MEMORY, required"`
	InstanceId     string `env:"INSTANCE_ID, required"`
}

func main() {
	fdk.Handle(fdk.HandlerFunc(resizeHandler))
}

func resizeHandler(ctx context.Context, in io.Reader, out io.Writer) {
	var cfg Config
	err := envconfig.Process(ctx, &cfg)
	helpers.FatalIfError(err)

	log.Printf("Instance ID: %s", cfg.InstanceId)
	log.Printf("Active OCPU: %d", cfg.ActiveOCPU)
	log.Printf("Active Memory: %d", cfg.ActiveMemory)
	log.Printf("Inactive OCPU: %d", cfg.InActiveOCPU)
	log.Printf("Inactive Memory: %d", cfg.InActiveMemory)

	provider, err := auth.ResourcePrincipalConfigurationProvider()
	helpers.FatalIfError(err)

	client, err := core.NewComputeClientWithConfigurationProvider(provider)
	helpers.FatalIfError(err)

	// Get the instance current details
	getReq := core.GetInstanceRequest{
		InstanceId: common.String(cfg.InstanceId),
	}

	resp, err := client.GetInstance(context.Background(), getReq)
	helpers.FatalIfError(err)

	// Check current instance shape and set target shape
	var targetOCPU, targetMemory int

	if *resp.ShapeConfig.Ocpus == float32(cfg.ActiveOCPU) && *resp.ShapeConfig.MemoryInGBs == float32(cfg.ActiveMemory) {
		targetOCPU = cfg.InActiveOCPU
		targetMemory = cfg.InActiveMemory
	} else if *resp.ShapeConfig.Ocpus == float32(cfg.InActiveOCPU) && *resp.ShapeConfig.MemoryInGBs == float32(cfg.InActiveMemory) {
		targetOCPU = cfg.ActiveOCPU
		targetMemory = cfg.ActiveMemory
	} else {
		json.NewEncoder(out).Encode("Instance shape is none of the specified shapes")
		return
	}

	// Resize the instance
	updateReq := core.UpdateInstanceRequest{
		InstanceId: common.String(cfg.InstanceId),
		UpdateInstanceDetails: core.UpdateInstanceDetails{
			ShapeConfig: &core.UpdateInstanceShapeConfigDetails{
				Ocpus:       common.Float32(float32(targetOCPU)),
				MemoryInGBs: common.Float32(float32(targetMemory)),
			},
		},
	}

	_, err = client.UpdateInstance(context.Background(), updateReq)
	helpers.FatalIfError(err)

	json.NewEncoder(out).Encode("Instance resized successfully")
}
