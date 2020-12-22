package cli

import (
	"fmt"
	"log"
	"monzo-ynab/commands"
	"monzo-ynab/rest"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

func NewCLI(c *commands.Commands, r *rest.Handler) *CLI {
	return &CLI{appCommands: c, restHandler: r}
}

type CLI struct {
	appCommands *commands.Commands
	restHandler *rest.Handler

	rootCmd *cobra.Command

	restPort string
}

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

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(syncCmd)
	c.rootCmd = rootCmd
}

func (c CLI) Run() {
	c.rootCmd.Execute()
}
