package cmd

import (
	"fmt"
	"os"

	"github.com/materials-commons/materials-commons-cli/pkg/config"
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
	remote, err := config.GetRemote()
	fmt.Printf("got remote = %#v\n", remote)
	mcapikeyFromEnv := os.Getenv("MCAPIKEY")
	mcurlFromEnv := os.Getenv("MCURL")

	if err != nil {
		fmt.Println("Error trying to read config:", err)

		if config.GetProjectMCConfig() == "" {
			fmt.Println("  Warning - No config.json found.")
		} else {
			fmt.Println("  Warning - Unable to remote config.json at:", config.GetProjectMCConfig())
		}

		if mcapikeyFromEnv != "" {
			fmt.Println("Your API Key is set from the environment (MCAPIKEY):", mcapikeyFromEnv)
		}

		if mcurlFromEnv != "" {
			fmt.Println("Your Server URL is set from the environment (MCURL):", mcurlFromEnv)
		}

		return
	}

	// If we are here then a config.json existed and was read
	fmt.Printf("Your config.json is located at %q\n", config.GetProjectMCConfig())

	if mcapikeyFromEnv != "" {
		fmt.Println("Your API Key is set from the environment (MCAPIKEY):", mcapikeyFromEnv)
	}

	if mcurlFromEnv != "" {
		fmt.Println("Your Server URL is set from the environment (MCURL):", mcurlFromEnv)
	}

	fmt.Println("Your config.json has the following settings:")
	if remote.DefaultRemote.MCAPIKey == "" && remote.DefaultRemote.MCUrl == "" {
		fmt.Println("   You have have no default remote set.")
		return
	}

	if remote.DefaultRemote.MCAPIKey == "" {
		fmt.Println("  You have a default remote, but the mcapikey isn't set.")
	} else {
		fmt.Printf("  Your default mcapikey is %q, unless overridden by the environment (MCAPIKEY).", remote.DefaultRemote.MCAPIKey)
	}

	if remote.DefaultRemote.MCUrl == "" {
		fmt.Println("  You have a default remote, but the mcurl isn't set.")
	} else {
		fmt.Printf("  Your default mcurl is %q, unless overridden by the environment (MCURL).", remote.DefaultRemote.MCUrl)
	}
}

func init() {
	rootCmd.AddCommand(remotesCmd)
}
