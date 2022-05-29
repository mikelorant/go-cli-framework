package cmd

import (
	"context"
	"fmt"

	"test5/internal/config"
	"test5/internal/database"
	"test5/internal/logging"
	"test5/internal/ssh"

	"github.com/spf13/cobra"
)

func NewSSHCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssh",
		Short: "A brief description of your application",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := RunSSH(cmd); err != nil {
				return fmt.Errorf("unable to run SSH command: %w", err)
			}

			return nil
		},
	}

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

func RunSSH(cmd *cobra.Command) error {
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

	log.Info("Hello from SSH Command")

	ssh := &ssh.SSH{}
	cfg.UnmarshalWithKey("ssh", ssh)
	log.SetOption(
		logging.WithField("ssh"),
	)
	log.Debug(ssh)

	db := &database.Database{}
	cfg.UnmarshalWithKey("database", db)
	log.SetOption(
		logging.WithField("database"),
	)
	log.Debug(db)

	return nil
}
