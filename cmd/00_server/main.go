package main

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/health-check/api/record"
	"github.com/sentinel-official/health-check/context"
)

const (
	appName = "00_server"
)

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

			router := gin.Default()
			router.Use(cors.Default())

			record.RegisterRoutes(ctx, router)

			if err := http.ListenAndServe(cfg.ListenAddr, router); err != nil {
				return err
			}

			return nil
		},
	}

	AddConfigFlagsToCmd(v)
	_ = v.Execute()
}
