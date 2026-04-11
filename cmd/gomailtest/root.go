package main

import (
	"github.com/spf13/cobra"
	"msgraphtool/internal/common/bootstrap"
	"msgraphtool/internal/protocols/imap"
	"msgraphtool/internal/protocols/jmap"
	"msgraphtool/internal/protocols/msgraph"
	"msgraphtool/internal/protocols/pop3"
	"msgraphtool/internal/protocols/smtp"
)

var rootCmd = &cobra.Command{
	Use:   "gomailtest",
	Short: "Email and calendar protocol testing tool",
	Long: `gomailtest is a unified CLI for testing email and calendar protocols.

It supports Microsoft Graph API (Exchange Online) and SMTP with more protocols coming in 3.0.

Run 'gomailtest <protocol> --help' for protocol-specific usage.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := bootstrap.SetupSignalContext()
		cmd.SetContext(ctx)
		_ = cancel // context lifetime tied to cmd context
		return nil
	},
}

func init() {
	rootCmd.AddCommand(msgraph.NewCmd())
	rootCmd.AddCommand(smtp.NewCmd())
	rootCmd.AddCommand(pop3.NewCmd())
	rootCmd.AddCommand(imap.NewCmd())
	rootCmd.AddCommand(jmap.NewCmd())
}

// Execute runs the root command and returns any error.
func Execute() error {
	return rootCmd.Execute()
}
