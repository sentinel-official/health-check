package main

import (
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sentinel-official/health-check/context"
	"github.com/sentinel-official/health-check/database"
	"github.com/sentinel-official/health-check/libs/geoip"
	"github.com/sentinel-official/health-check/types"
	"github.com/sentinel-official/health-check/utils"
)

const appName = "02_client"

func main() {
	v := &cobra.Command{
		Use:          appName,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfigFromFlags(cmd.Flags())
			if err != nil {
				return err
			}

			ctx := context.NewDefaultContext().
				WithAppName(cfg.App.Name)

			var (
				dbPassword = os.Getenv("DATABASE_PASSWORD")
				dbUsername = os.Getenv("DATABASE_USERNAME")
			)

			db, err := ctx.PrepareDatabase(dbUsername, dbPassword, cfg.Database.URI, cfg.Database.Name)
			if err != nil {
				return err
			}

			ctx = ctx.WithDatabase(db)

			if err := createIndexes(ctx); err != nil {
				return err
			}

			record, err := connect(ctx, cfg.Address)
			if err != nil {
				return err
			}

			time.Sleep(5 * time.Second)

			var transport *http.Transport
			if record.Type == types.NodeTypeWireGuard {
				transport, err = utils.NewTransport("")
				if err != nil {
					return err
				}

			} else if record.Type == types.NodeTypeV2Ray {
				transport, err = utils.NewTransport("socks5://127.0.0.1:1080")
				if err != nil {
					return err
				}
			}

			filter := bson.M{
				"addr": cfg.Address,
			}
			update := bson.M{}

			location, err := geoip.Location(transport)
			if err != nil {
				update = bson.M{
					"$set": bson.M{
						"location_fetch_error":     err.Error(),
						"location_fetch_timestamp": time.Now().UTC(),
					},
				}
			} else {
				update = bson.M{
					"$set": bson.M{
						"city":                     location.City,
						"country":                  location.Country,
						"ip_addr":                  location.IP,
						"latitude":                 location.Latitude,
						"location_fetch_error":     "",
						"location_fetch_timestamp": time.Now().UTC(),
						"longitude":                location.Longitude,
					},
				}
			}

			if _, err := database.TempRecordFindOneAndUpdate(ctx, filter, update); err != nil {
				return err
			}

			return nil
		},
	}

	AddConfigFlagsToCmd(v)
	_ = v.Execute()
}
