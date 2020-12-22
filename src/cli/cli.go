package cli

import (
	"fmt"
	"log"
	"monzo-ynab/commands"
	"monzo-ynab/internal/config"
	"monzo-ynab/rest"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

// NewCLI returns a configured CLI
func NewCLI(c *commands.Commands, r *rest.Handler, i *Installer, cfg config.Config) *CLI {
	return &CLI{appCommands: c, restHandler: r, installer: i, config: cfg}
}

// CLI is the apps command line interface
type CLI struct {
	config      config.Config
	appCommands *commands.Commands
	restHandler *rest.Handler
	installer   *Installer

	rootCmd *cobra.Command

	restPort string
}

// Init builds the CLI
func (c *CLI) Init() {
	rootCmd := &cobra.Command{Use: "monzo-ynab"}

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the server",
		Long:  "Run the server and listen for Monzo webhooks",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("Listening for transactions on port %s", c.restPort)
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", c.restPort), c.restHandler))
		},
	}
	runCmd.Flags().StringVarP(&c.restPort, "port", "p", "80", "Server port")

	syncCmd := &cobra.Command{
		Use:   "sync [days]",
		Short: "Syncs Monzo to YNAB",
		Long:  "Syncs Monzo transactions to YNAB for the past specified days",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			days, _ := strconv.Atoi(args[0])
			err := c.appCommands.Sync.Execute(days)
			if err != nil {
				log.Printf("Something went wrong: %s", err)
			}
		},
	}

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Installs Monzo-to-YNAB",
		Long:  "Runs the interactive installer the configure and setup Monzo-to-YNAB",
		Run: func(cmd *cobra.Command, args []string) {
			c.installer.Install()
		},
	}

	configCmd := &cobra.Command{
		Use:   "configure [option] [value]",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			val := args[1]

			err := c.config.Set(key, val)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	webhookCmd := &cobra.Command{
		Use:   "register-webhooks",
		Short: "Register Monzo webhooks",
		Run: func(cmd *cobra.Command, args []string) {
			c.appCommands.RegisterMonzoWebhook.Execute("/events/monzo/")
		},
	}

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(webhookCmd)
	c.rootCmd = rootCmd
}

// Run runs the CLI
func (c CLI) Run() {
	c.rootCmd.Execute()
}
