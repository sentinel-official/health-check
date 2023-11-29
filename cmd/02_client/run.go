package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/sentinel-official/health-check/context"
	"github.com/sentinel-official/health-check/database"
	"github.com/sentinel-official/health-check/models"
	"github.com/sentinel-official/health-check/types"
	v2raytypes "github.com/sentinel-official/health-check/types/v2ray"
	wireguardtypes "github.com/sentinel-official/health-check/types/wireguard"
)

func connect(ctx *context.Context, addr string) (*models.Record, error) {
	filter := bson.M{
		"addr": addr,
		"client_config": bson.M{
			"$exists": true,
		},
		"server_config": bson.M{
			"$exists": true,
		},
	}

	record, err := database.RecordFindOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, fmt.Errorf("record %s does not exist", addr)
	}

	if record.Type == types.NodeTypeWireGuard {
		cfg, err := wireguardtypes.NewConfigFromBytes(record.ClientConfig, record.ServerConfig)
		if err != nil {
			return nil, err
		}

		if err := cfg.WriteToFile("/etc/wireguard/wg0.conf"); err != nil {
			return nil, err
		}

		cmd := exec.Command("wg-quick", strings.Split("up wg0", " ")...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			return nil, err
		}
	} else if record.Type == types.NodeTypeV2Ray {
		cfg, err := v2raytypes.NewConfigFromBytes(record.ClientConfig, record.ServerConfig)
		if err != nil {
			return nil, err
		}

		if err := cfg.WriteToFile("/config.json"); err != nil {
			return nil, err
		}

		cmd := exec.Command("v2ray", strings.Split("run --config /config.json", " ")...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("invalid node type %d", record.Type)
	}

	return record, nil
}
