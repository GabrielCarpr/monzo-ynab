package commands

// Commands is a containing struct of the app commands.
type Commands struct {
	*Sync
	*RegisterMonzoWebhook
	*Store
}
