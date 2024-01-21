package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/health-check/context"
)

const appName = "01_runner"

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	v := &cobra.Command{
		Use:          appName,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfigFromFlags(cmd.Flags())
			if err != nil {
				return err
			}

			ctx := context.NewDefaultContext().
				WithAppName(cfg.App.Name).
				WithBroadcastMode(cfg.Tx.BroadcastMode).
				WithChainID(cfg.Chain.ID).
				WithFees(cfg.Tx.Fees).
				WithFromName(cfg.Key.Name).
				WithGas(cfg.Tx.Gas).
				WithGasAdjustment(cfg.Tx.GasAdjustment).
				WithGasPrices(cfg.Tx.GasPrices).
				WithInMemoryKeyring().
				WithMemo(cfg.Tx.Memo).
				WithQueryMaxTries(cfg.Query.MaxTries).
				WithSignMode(cfg.Tx.SignMode).
				WithSimulateAndExecute(cfg.Tx.SimulateAndExecute).
				WithTimeoutHeight(cfg.Tx.TimeoutHeight)

			ctx, err = ctx.WithFeeGranterAddr(cfg.Tx.FeeGranterAddr)
			if err != nil {
				return err
			}

			ctx, err = ctx.WithRPCAddr(cfg.Chain.RPCAddr)
			if err != nil {
				return err
			}

			var (
				dbPassword = os.Getenv("DATABASE_PASSWORD")
				dbUsername = os.Getenv("DATABASE_USERNAME")
				mnemonic   = os.Getenv("MNEMONIC")
			)

			_, err = ctx.WithKey(mnemonic, cfg.Key.Type, cfg.Key.Account, cfg.Key.Index, cfg.Key.BIP39Passphrase)
			if err != nil {
				return err
			}

			db, err := ctx.PrepareDatabase(dbUsername, dbPassword, cfg.Database.URI, cfg.Database.Name)
			if err != nil {
				return err
			}

			ctx = ctx.WithDatabase(db)

			if err := createIndexes(ctx); err != nil {
				return err
			}

			if err := startTransaction(ctx); err != nil {
				return err
			}
			if err := queryNodes(ctx, cfg.PaymentDenom); err != nil {
				return err
			}
			if err := updateNodeInfos(ctx, cfg.RequestTimeout); err != nil {
				return err
			}
			if err := cancelSubscriptions(ctx, cfg.MaxMsgs); err != nil {
				return err
			}
			if err := startSubscriptions(ctx, cfg.MaxMsgs, cfg.MaxGigabytePrice, cfg.PaymentDenom); err != nil {
				return err
			}
			if err := endSessions(ctx, cfg.MaxMsgs); err != nil {
				return err
			}
			if err := startSessions(ctx, cfg.MaxMsgs); err != nil {
				return err
			}
			if err := updateClientConfigs(ctx, cfg.RequestTimeout); err != nil {
				return err
			}
			if err := updateClients(ctx); err != nil {
				return err
			}
			if err := updateDuplicateIPAddrs(ctx); err != nil {
				return err
			}
			if err := updateOKs(ctx); err != nil {
				return err
			}
			if err := commitTransaction(ctx); err != nil {
				return err
			}

			return nil
		},
	}

	AddConfigFlagsToCmd(v)
	_ = v.Execute()
}
