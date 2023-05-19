package cmd

import (
	"fmt"
	"os"

	"github.com/materials-commons/materials-commons-cli/pkg/model"
	"github.com/materials-commons/materials-commons-cli/pkg/stor"
	"github.com/spf13/cobra"
)

// remotesCmd represents the remotes command
var remotesCmd = &cobra.Command{
	Use:   "remotes",
	Short: "Show remote being used and other known remotes.",
	Long:  `Show remote being used and other known remotes.`,
	Run:   runRemotesCmd,
}

func runRemotesCmd(cmd *cobra.Command, args []string) {
	remoteStor := stor.MustLoadJsonRemoteStor()
	defaultRemote, err := remoteStor.GetDefaultRemote()
	mcapikeyFromEnv := os.Getenv("MCAPIKEY")
	mcurlFromEnv := os.Getenv("MCURL")

	if err != nil {
		fmt.Printf("No default defaultRemote set: %s", err)
	} else {

		if mcapikeyFromEnv != "" {
			fmt.Printf("The default default remote apikey is set from the environment (MCAPIKEY).\n")
		}

		if mcurlFromEnv != "" {
			fmt.Printf("The default default remote server url is set from the environment (MCURL).\n")
		}

		fmt.Printf("\nDefault Remote:\n")
		fmt.Printf("  EMail    : %s\n", defaultRemote.EMail)
		fmt.Printf("  ServerURL: %s\n", defaultRemote.MCUrl)
		fmt.Printf("  APIKey   : %s\n", defaultRemote.MCAPIKey)
	}

	fmt.Printf("\nOther Remotes:\n")
	err = remoteStor.ListPaged(func(r *model.Remote) error {
		if isSameAsDefaultRemote(defaultRemote, r) {
			// Skip showing a remote equal to the default remote
			return nil
		}
		fmt.Printf("    Remote:\n")
		fmt.Printf("      EMail    : %s\n", r.EMail)
		fmt.Printf("      ServerURL: %s\n", r.MCUrl)
		fmt.Printf("      APIKey   : %s\n\n", r.MCAPIKey)
		return nil
	})

	if err != nil {
		fmt.Printf("Error listing other remotes: %s\n", err)
	}
}

func isSameAsDefaultRemote(defaultRemote, remote *model.Remote) bool {
	if defaultRemote.MCAPIKey == remote.MCAPIKey &&
		defaultRemote.MCUrl == remote.MCUrl &&
		defaultRemote.EMail == remote.EMail {
		return true
	}

	return false
}

func init() {
	rootCmd.AddCommand(remotesCmd)
}
