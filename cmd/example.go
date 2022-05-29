package cmd

import (
	"context"
	"fmt"

	"github.com/mikelorant/go-cli-framework/internal/database"
	"github.com/mikelorant/go-cli-framework/internal/ssh"
	"github.com/mikelorant/go-cli-framework/pkg/config"
	"github.com/mikelorant/go-cli-framework/pkg/logging"

	"github.com/spf13/cobra"
)

func NewExampleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "example",
		Short: "A brief description of your application",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := RunExample(cmd); err != nil {
				return fmt.Errorf("unable to run SSH command: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().String("key", "id_rsa", "SSH key")
	cmd.Flags().String("user", "root", "SSH user")
	cmd.Flags().String("host", "", "SSH host")
	cmd.Flags().String("db-user", "admin", "Database user")
	cmd.Flags().String("db-password", "", "Database password")
	cmd.Flags().String("db-host", "", "Database host")
	cmd.Flags().String("db-port", "3306", "Database port")

	cmd.Flags().SetAnnotation("key", "bindWithKey", []string{"ssh.key"})
	cmd.Flags().SetAnnotation("user", "bindWithKey", []string{"ssh.user"})
	cmd.Flags().SetAnnotation("host", "bindWithKey", []string{"ssh.host"})
	cmd.Flags().SetAnnotation("db-user", "bindWithKey", []string{"database.user"})
	cmd.Flags().SetAnnotation("db-password", "bindWithKey", []string{"database.password"})
	cmd.Flags().SetAnnotation("db-host", "bindWithKey", []string{"database.host"})
	cmd.Flags().SetAnnotation("db-port", "bindWithKey", []string{"database.port"})

	return cmd
}

func RunExample(cmd *cobra.Command) error {
	ctx := context.Background()
	log := logging.New(ctx,
		logging.WithLevel(true),
	)
	cfg := config.New(ctx)

	err := cfg.Load(
		config.WithEnvPrefix(_configEnvPrefix),
		config.WithFilename(_configFilename),
		config.WithFlags(cmd.Flags()),
	)
	if err != nil {
		return fmt.Errorf("unable to load config: %w", err)
	}

	log.SetLevel("debug")

	log.Info("Hello from Example Command")

	ssh := &ssh.SSH{}
	cfg.UnmarshalWithKey("ssh", ssh)
	log.Debug(ssh)

	db := &database.Database{}
	cfg.UnmarshalWithKey("database", db)
	log.Debug(db)

	return nil
}
