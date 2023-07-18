package main

import (
	"fmt"
	"os"

	"github.com/redhatinsights/platform-changelog-go/internal/config"
	"github.com/redhatinsights/platform-changelog-go/internal/logging"
	"github.com/spf13/cobra"
)

func createCommands(cfg *config.Config) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use: "changelog",
	}

	var serverCmd = &cobra.Command{
		Use:   "api",
		Short: "Runs the API server",
		Run: func(cmd *cobra.Command, args []string) {
			startAPI(cfg)
		},
	}

	rootCmd.AddCommand(serverCmd)

	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Runs the DB migrations",
	}

	// upward migration
	var upCmd = &cobra.Command{
		Use:   "up",
		Short: "Runs the DB migrations with latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Upward migration")
			return migrateDB(cfg, "up")
		},
	}

	// downward migration
	var downCmd = &cobra.Command{
		Use:   "down",
		Short: "Migrates the DB down one version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Downward migration")
			return migrateDB(cfg, "down")
		},
	}

	// seed the DB
	var seedCmd = &cobra.Command{
		Use:   "seed",
		Short: "Seeds the DB with the services from the config",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Seeding DB")
			return seedDB(cfg)
		},
	}

	rootCmd.AddCommand(seedCmd)
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(upCmd)
	migrateCmd.AddCommand(downCmd)

	return rootCmd
}

func main() {
	logging.InitLogger()
	cfg := config.Get()

	rootCmd := createCommands(cfg)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
