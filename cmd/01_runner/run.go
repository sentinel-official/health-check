package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/hashicorp/go-uuid"
	hubtypes "github.com/sentinel-official/hub/types"
	nodetypes "github.com/sentinel-official/hub/x/node/types"
	sessiontypes "github.com/sentinel-official/hub/x/session/types"
	subscriptiontypes "github.com/sentinel-official/hub/x/subscription/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"

	"github.com/sentinel-official/health-check/context"
	"github.com/sentinel-official/health-check/database"
	"github.com/sentinel-official/health-check/types"
	wireguardtypes "github.com/sentinel-official/health-check/types/wireguard"
)

func startTransaction(ctx *context.Context) error {
	filter := bson.M{}
	projection := bson.M{
		"addr":            1,
		"client_config":   1,
		"server_config":   1,
		"session_id":      1,
		"subscription_id": 1,
	}
	opts := options.Find().
		SetProjection(projection)

	items, err := database.RecordFindAll(ctx, filter, opts)
	if err != nil {
		return err
	}

	if _, err := database.TempRecordDeleteMany(ctx, filter); err != nil {
		return err
	}

	if _, err := database.TempRecordInsertMany(ctx, items); err != nil {
		return err
	}

	return nil
}

func queryNodes(ctx *context.Context, paymentDenom string) error {
	pagination := &query.PageRequest{
		Limit: 1e9,
	}

	nodes, err := ctx.QueryNodes(hubtypes.StatusActive, pagination)
	if err != nil {
		return err
	}

	var addrs []string
	for i := 0; i < len(nodes); i++ {
		price, ok := nodes[i].GigabytePrice(paymentDenom)
		if !ok {
			continue
		}

		filter := bson.M{
			"addr": nodes[i].Address,
		}
		update := bson.M{
			"$set": bson.M{
				"gigabyte_price": price.Amount.Int64(),
				"remote_url":     nodes[i].RemoteURL,
				"status":         hubtypes.StatusActive,
			},
		}
		opts := options.FindOneAndUpdate().
			SetUpsert(true)

		addrs = append(addrs, nodes[i].Address)
		if _, err := database.TempRecordFindOneAndUpdate(ctx, filter, update, opts); err != nil {
			return err
		}
	}

	filter := bson.M{
		"addr": bson.M{
			"$nin": addrs,
		},
	}
	update := bson.M{
		"$set": bson.M{
			"status": hubtypes.StatusInactive,
		},
	}

	if _, err := database.TempRecordUpdateMany(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func updateNodeInfos(ctx *context.Context, timeout time.Duration) error {
	filter := bson.M{
		"status": hubtypes.StatusActive,
	}

	nodes, err := database.TempRecordFindAll(ctx, filter)
	if err != nil {
		return err
	}

	group := &errgroup.Group{}
	group.SetLimit(64)

	for i := 0; i < len(nodes); i++ {
		var (
			nodeAddr  = nodes[i].Addr
			remoteURL = nodes[i].RemoteURL
		)

		group.Go(func() error {
			filter := bson.M{
				"addr": nodeAddr,
			}
			update := bson.M{}

			info, err := types.FetchNewNodeInfo(remoteURL, timeout)
			if err != nil {
				update = bson.M{
					"$set": bson.M{
						"info_fetch_error":     err.Error(),
						"info_fetch_timestamp": time.Now().UTC(),
					},
				}
			} else {
				update = bson.M{
					"$set": bson.M{
						"info_fetch_error":     "",
						"info_fetch_timestamp": time.Now().UTC(),
						"type":                 types.NewNodeTypeFromUInt64(info.Type),
					},
				}
			}

			if _, err := database.TempRecordFindOneAndUpdate(ctx, filter, update); err != nil {
				return err
			}

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}

func cancelSubscriptions(ctx *context.Context, maxMsgs int) error {
	log.Println("cancelSubscriptions")

	fromAddr, err := ctx.FromAddr()
	if err != nil {
		return err
	}

	pagination := &query.PageRequest{
		Limit: 1e9,
	}

	subscriptions, err := ctx.QuerySubscriptionsForAccount(fromAddr, pagination)
	if err != nil {
		return err
	}

	var (
		activeIDs      []uint64
		msgs           []sdk.Msg
		bech32FromAddr = fromAddr.String()
	)

	for i := 0; i < len(subscriptions); i++ {
		if !subscriptions[i].GetStatus().Equal(hubtypes.StatusActive) {
			continue
		}

		filter := bson.M{
			"subscription_id": subscriptions[i].GetID(),
		}

		record, err := database.TempRecordFindOne(ctx, filter)
		if err != nil {
			return err
		}
		if record == nil {
			msgs = append(
				msgs,
				&subscriptiontypes.MsgCancelRequest{
					From: bech32FromAddr,
					ID:   subscriptions[i].GetID(),
				},
			)

			continue
		}

		activeIDs = append(activeIDs, subscriptions[i].GetID())
	}

	filter := bson.M{}
	if len(activeIDs) > 0 {
		filter = bson.M{
			"subscription_id": bson.M{
				"$nin": activeIDs,
			},
		}
	}

	update := bson.M{
		"$unset": bson.M{
			"client_config":   1,
			"server_config":   1,
			"session_id":      1,
			"subscription_id": 1,
		},
	}

	if _, err := database.TempRecordUpdateMany(ctx, filter, update); err != nil {
		return err
	}

	log.Println("msgsLen", len(msgs))
	if len(msgs) == 0 {
		return nil
	}

	for start, end := 0, 0; start < len(msgs); start = end {
		end = start + maxMsgs
		if end > len(msgs) {
			end = len(msgs)
		}

		resp, err := ctx.Tx(msgs[start:end]...)
		if err != nil {
			return err
		}

		result, err := ctx.QueryTxWithRetry(resp.TxHash)
		if err != nil {
			return err
		}
		if result == nil {
			return fmt.Errorf("nil query result for the transaction %s", resp.TxHash)
		}
		if !result.TxResult.IsOK() {
			return fmt.Errorf("transaction %s failed with the code %d", resp.TxHash, result.TxResult.Code)
		}
	}

	return nil
}

func startSubscriptions(ctx *context.Context, maxMsgs int, maxGigabytePrice int64, paymentDenom string) error {
	log.Println("startSubscriptions")

	filter := bson.M{
		"gigabyte_price": bson.M{
			"$lt": maxGigabytePrice,
		},
		"status": hubtypes.StatusActive,
		"subscription_id": bson.M{
			"$exists": false,
		},
	}

	records, err := database.TempRecordFindAll(ctx, filter)
	if err != nil {
		return err
	}

	fromAddr, err := ctx.FromAddr()
	if err != nil {
		return err
	}

	var (
		msgs           []sdk.Msg
		bech32FromAddr = fromAddr.String()
	)

	for i := 0; i < len(records); i++ {
		msgs = append(
			msgs,
			&nodetypes.MsgSubscribeRequest{
				From:        bech32FromAddr,
				NodeAddress: records[i].Addr,
				Gigabytes:   1,
				Hours:       0,
				Denom:       paymentDenom,
			},
		)
	}

	log.Println("msgsLen", len(msgs))
	if len(msgs) == 0 {
		return nil
	}

	for start, end := 0, 0; start < len(msgs); start = end {
		end = start + maxMsgs
		if end > len(msgs) {
			end = len(msgs)
		}

		resp, err := ctx.Tx(msgs[start:end]...)
		if err != nil {
			return err
		}

		result, err := ctx.QueryTxWithRetry(resp.TxHash)
		if err != nil {
			return err
		}
		if result == nil {
			return fmt.Errorf("nil query result for the transaction %s", resp.TxHash)
		}
		if !result.TxResult.IsOK() {
			return fmt.Errorf("transaction %s failed with the code %d", resp.TxHash, result.TxResult.Code)
		}

		for _, event := range result.TxResult.Events {
			if event.Type == "sentinel.node.v2.EventCreateSubscription" {
				var (
					id       uint64
					nodeAddr string
				)

				for _, attribute := range event.Attributes {
					var (
						key   = string(attribute.Key)
						value = string(attribute.Value)
					)

					switch key {
					case "id":
						id, err = strconv.ParseUint(value[1:len(value)-1], 10, 64)
						if err != nil {
							return err
						}
					case "node_address":
						nodeAddr = value[1 : len(value)-1]
					}
				}

				filter := bson.M{
					"addr": nodeAddr,
				}
				update := bson.M{
					"$set": bson.M{
						"subscription_id": id,
					},
				}

				if _, err := database.TempRecordFindOneAndUpdate(ctx, filter, update); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func endSessions(ctx *context.Context, maxMsgs int) error {
	log.Println("endSessions")

	fromAddr, err := ctx.FromAddr()
	if err != nil {
		return err
	}

	pagination := &query.PageRequest{
		Limit: 1e9,
	}

	sessions, err := ctx.QuerySessionsForAccount(fromAddr, pagination)
	if err != nil {
		return err
	}

	var (
		activeIDs      []uint64
		msgs           []sdk.Msg
		bech32FromAddr = fromAddr.String()
	)

	for i := 0; i < len(sessions); i++ {
		if !sessions[i].Status.Equal(hubtypes.StatusActive) {
			continue
		}

		filter := bson.M{
			"session_id": sessions[i].ID,
		}

		record, err := database.TempRecordFindOne(ctx, filter)
		if err != nil {
			return err
		}
		if record == nil {
			msgs = append(
				msgs,
				&sessiontypes.MsgEndRequest{
					From:   bech32FromAddr,
					ID:     sessions[i].ID,
					Rating: 0,
				},
			)

			continue
		}

		activeIDs = append(activeIDs, sessions[i].ID)
	}

	filter := bson.M{}
	if len(activeIDs) > 0 {
		filter = bson.M{
			"session_id": bson.M{
				"$nin": activeIDs,
			},
		}
	}

	update := bson.M{
		"$unset": bson.M{
			"client_config": 1,
			"server_config": 1,
			"session_id":    1,
		},
	}

	if _, err := database.TempRecordUpdateMany(ctx, filter, update); err != nil {
		return err
	}

	log.Println("msgsLen", len(msgs))
	if len(msgs) == 0 {
		return nil
	}

	for start, end := 0, 0; start < len(msgs); start = end {
		end = start + maxMsgs
		if end > len(msgs) {
			end = len(msgs)
		}

		resp, err := ctx.Tx(msgs[start:end]...)
		if err != nil {
			return err
		}

		result, err := ctx.QueryTxWithRetry(resp.TxHash)
		if err != nil {
			return err
		}
		if result == nil {
			return fmt.Errorf("nil query result for the transaction %s", resp.TxHash)
		}
		if !result.TxResult.IsOK() {
			return fmt.Errorf("transaction %s failed with the code %d", resp.TxHash, result.TxResult.Code)
		}
	}

	return nil
}

func startSessions(ctx *context.Context, maxMsgs int) error {
	log.Println("startSessions")

	filter := bson.M{
		"session_id": bson.M{
			"$exists": false,
		},
		"status": hubtypes.StatusActive,
		"subscription_id": bson.M{
			"$exists": true,
		},
	}

	records, err := database.TempRecordFindAll(ctx, filter)
	if err != nil {
		return err
	}

	fromAddr, err := ctx.FromAddr()
	if err != nil {
		return err
	}

	var (
		msgs           []sdk.Msg
		bech32FromAddr = fromAddr.String()
	)

	for i := 0; i < len(records); i++ {
		msgs = append(
			msgs,
			&sessiontypes.MsgStartRequest{
				From:    bech32FromAddr,
				ID:      records[i].SubscriptionID,
				Address: records[i].Addr,
			},
		)
	}

	log.Println("msgsLen", len(msgs))
	if len(msgs) == 0 {
		return nil
	}

	for start, end := 0, 0; start < len(msgs); start = end {
		end = start + maxMsgs
		if end > len(msgs) {
			end = len(msgs)
		}

		resp, err := ctx.Tx(msgs[start:end]...)
		if err != nil {
			return err
		}

		result, err := ctx.QueryTxWithRetry(resp.TxHash)
		if err != nil {
			return err
		}
		if result == nil {
			return fmt.Errorf("nil query result for the transaction %s", resp.TxHash)
		}
		if !result.TxResult.IsOK() {
			return fmt.Errorf("transaction %s failed with the code %d", resp.TxHash, result.TxResult.Code)
		}

		for _, event := range result.TxResult.Events {
			if event.Type == "sentinel.session.v2.EventStart" {
				var (
					id       uint64
					nodeAddr string
				)

				for _, attribute := range event.Attributes {
					var (
						key   = string(attribute.Key)
						value = string(attribute.Value)
					)

					switch key {
					case "id":
						id, err = strconv.ParseUint(value[1:len(value)-1], 10, 64)
						if err != nil {
							return err
						}
					case "node_address":
						nodeAddr = value[1 : len(value)-1]
					}
				}

				filter := bson.M{
					"addr": nodeAddr,
				}
				update := bson.M{
					"$set": bson.M{
						"session_id": id,
					},
				}

				if _, err := database.TempRecordFindOneAndUpdate(ctx, filter, update); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func updateClientConfigs(ctx *context.Context, timeout time.Duration) error {
	log.Println("updateClientConfigs")

	fromAddr, err := ctx.FromAddr()
	if err != nil {
		return err
	}

	filter := bson.M{
		"session_id": bson.M{
			"$exists": true,
		},
		"status": hubtypes.StatusActive,
		"subscription_id": bson.M{
			"$exists": true,
		},
	}

	records, err := database.TempRecordFindAll(ctx, filter)
	if err != nil {
		return err
	}

	group := &errgroup.Group{}
	group.SetLimit(64)

	for i := 0; i < len(records); i++ {
		var (
			clientConfig = records[i].ClientConfig
			nodeAddr     = records[i].Addr
			nodeType     = records[i].Type
			remoteURL    = records[i].RemoteURL
			serverConfig = records[i].ServerConfig
			sessionID    = records[i].SessionID
		)

		group.Go(func() error {
			f := func() ([]byte, []byte, error) {
				var (
					config []byte
					key    string
				)

				switch nodeType {
				case types.NodeTypeWireGuard:
					privKey, err := wireguardtypes.NewPrivateKey()
					if err != nil {
						return nil, nil, err
					}

					config = append([]byte{}, privKey[:]...)
					key = privKey.Public().String()
				case types.NodeTypeV2Ray:
					uid, err := uuid.GenerateRandomBytes(16)
					if err != nil {
						return nil, nil, err
					}

					config = append([]byte{}, uid[:]...)
					key = base64.StdEncoding.EncodeToString(append([]byte{0x01}, uid...))
				default:
					return nil, nil, fmt.Errorf("invalid node type %s", nodeType)
				}

				signature, _, err := ctx.Sign(sdk.Uint64ToBigEndian(sessionID))
				if err != nil {
					return nil, nil, err
				}

				req, err := json.Marshal(
					map[string]interface{}{
						"key":       key,
						"signature": signature,
					},
				)
				if err != nil {
					return nil, nil, err
				}

				client := &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
						},
					},
					Timeout: timeout,
				}

				urlPath, err := url.JoinPath(remoteURL, fmt.Sprintf("/accounts/%s/sessions/1", fromAddr))
				if err != nil {
					return nil, nil, err
				}

				_, err = client.Post(urlPath, jsonrpc.ContentType, bytes.NewBuffer(req))
				if err != nil {
					return nil, nil, err
				}

				if len(clientConfig) != 0 && len(serverConfig) != 0 {
					return clientConfig, serverConfig, nil
				}

				urlPath, err = url.JoinPath(remoteURL, fmt.Sprintf("/accounts/%s/sessions/%d", fromAddr, sessionID))
				if err != nil {
					return nil, nil, err
				}

				resp, err := client.Post(urlPath, jsonrpc.ContentType, bytes.NewBuffer(req))
				if err != nil {
					return nil, nil, err
				}

				defer func() {
					if err := resp.Body.Close(); err != nil {
						panic(err)
					}
				}()

				var m map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
					return nil, nil, err
				}
				if m["error"] != nil {
					return nil, nil, fmt.Errorf("%s", m["error"])
				}

				result, err := base64.StdEncoding.DecodeString(m["result"].(string))
				if err != nil {
					return nil, nil, err
				}

				return config, result, nil
			}

			filter := bson.M{
				"addr": nodeAddr,
			}
			update := bson.M{}

			clientConfig, serverConfig, err := f()
			if err != nil {
				update = bson.M{
					"$set": bson.M{
						"config_exchange_error":     err.Error(),
						"config_exchange_timestamp": time.Now().UTC(),
					},
					"$unset": bson.M{
						"client_config": 1,
						"server_config": 1,
					},
				}
			} else {
				update = bson.M{
					"$set": bson.M{
						"client_config":             clientConfig,
						"config_exchange_error":     "",
						"config_exchange_timestamp": time.Now().UTC(),
						"server_config":             serverConfig,
					},
				}
			}

			if _, err := database.TempRecordFindOneAndUpdate(ctx, filter, update); err != nil {
				return err
			}

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}

func updateClients(ctx *context.Context) error {
	log.Println("updateClients")

	filter := bson.M{
		"client_config": bson.M{
			"$exists": true,
		},
		"server_config": bson.M{
			"$exists": true,
		},
		"session_id": bson.M{
			"$exists": true,
		},
		"status": hubtypes.StatusActive,
		"subscription_id": bson.M{
			"$exists": true,
		},
	}

	records, err := database.TempRecordFindAll(ctx, filter)
	if err != nil {
		return err
	}

	group := &errgroup.Group{}
	group.SetLimit(64)

	for i := 0; i < len(records); i++ {
		if i%100 == 0 {
			log.Println(i, records[i].Addr, records[i].SubscriptionID, records[i].SessionID)
		}

		nodeAddr := records[i].Addr
		group.Go(func() error {
			args := strings.Split(
				fmt.Sprintf("run --privileged --rm --tty health-check-client main --address=%s --database.uri=mongodb://172.17.0.1:27017", nodeAddr),
				" ")
			cmd := exec.Command("docker", args...)
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			cmd.Stdout = io.Discard

			if err := cmd.Run(); err != nil {
				return err
			}

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}

func updateDuplicateIPAddrs(ctx *context.Context) error {
	log.Println("updateDuplicateIPAddrs")

	filter := bson.M{
		"status": hubtypes.StatusActive,
		"info_fetch_timestamp": bson.M{
			"$gt": time.Time{},
		},
		"info_fetch_error": "",
		"ip_addr": bson.M{
			"$exists": true,
		},
		"config_exchange_timestamp": bson.M{
			"$gt": time.Time{},
		},
		"config_exchange_error": "",
		"location_fetch_timestamp": bson.M{
			"$gt": time.Time{},
		},
		"location_fetch_error": "",
	}

	records, err := database.TempRecordFindAll(ctx, filter)
	if err != nil {
		return err
	}

	ipAddrs := make(map[string]int)
	for i := 0; i < len(records); i++ {
		if _, ok := ipAddrs[records[i].IPAddr]; !ok {
			ipAddrs[records[i].IPAddr] = 0
		}

		ipAddrs[records[i].IPAddr] += 1
	}

	for ipAddr, i := range ipAddrs {
		filter := bson.M{
			"status": hubtypes.StatusActive,
			"info_fetch_timestamp": bson.M{
				"$gt": time.Time{},
			},
			"info_fetch_error": "",
			"ip_addr":          ipAddr,
			"config_exchange_timestamp": bson.M{
				"$gt": time.Time{},
			},
			"config_exchange_error": "",
			"location_fetch_timestamp": bson.M{
				"$gt": time.Time{},
			},
			"location_fetch_error": "",
		}
		update := bson.M{}

		if i > 1 {
			update = bson.M{
				"$set": bson.M{
					"duplicate_ip_addr": true,
				},
			}
		} else {
			update = bson.M{
				"$set": bson.M{
					"duplicate_ip_addr": false,
				},
			}
		}

		if _, err := database.TempRecordUpdateMany(ctx, filter, update); err != nil {
			return err
		}
	}

	return nil
}

func updateOKs(ctx *context.Context) error {
	log.Println("updateOKs")

	filter := bson.M{
		"status": hubtypes.StatusActive,
		"info_fetch_timestamp": bson.M{
			"$gt": time.Time{},
		},
		"info_fetch_error": "",
		"ip_addr": bson.M{
			"$exists": true,
		},
		"config_exchange_timestamp": bson.M{
			"$gt": time.Time{},
		},
		"config_exchange_error": "",
		"duplicate_ip_addr":     false,
		"location_fetch_timestamp": bson.M{
			"$gt": time.Time{},
		},
		"location_fetch_error": "",
	}
	update := bson.M{
		"$set": bson.M{
			"ok": true,
		},
	}

	if _, err := database.TempRecordUpdateMany(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func commitTransaction(ctx *context.Context) error {
	filter := bson.M{}
	items, err := database.TempRecordFindAll(ctx, filter)
	if err != nil {
		return err
	}

	if _, err := database.RecordDeleteMany(ctx, filter); err != nil {
		return err
	}

	if _, err := database.RecordInsertMany(ctx, items); err != nil {
		return err
	}

	return nil
}
