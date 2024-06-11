package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	setupCmd = &cobra.Command{
		Use:   "setup server_id",
		Short: "Setup server for use with wolf",
		Args:  cobra.ExactArgs(1), // Expects exactly one argument: server_id
		RunE: func(cmd *cobra.Command, args []string) error {
			serverID := args[0]

			cmd.Flags().String("bin", "ssh", "Name of SSH client executable (e.g., ssh, mosh)")
			cmd.Flags().String("user", "user", "User account to use for login")
			cmd.Flags().String("command", "", "Command to execute over SSH")

			return setupServerCmd(cmd, serverID)
		},
	}
)

func init() {
	rootCmd.AddCommand(setupCmd)
}

func setupServerCmd(cmd *cobra.Command, server string) error {
	err := getSetupFiles(cmd, server)
	if err != nil {
		return fmt.Errorf("error getting setup files: %w", err)
	}

	startScriptCommand := "bash /home/user/setup.sh"

	// Set the command to be executed over SSH
	cmd.Flags().Set("command", startScriptCommand)

	// Execute the SSH command to run the script
	return sshServer(cmd, []string{server})
}

func getSetupFiles(cmd *cobra.Command, server string) error {
	// Define the files and their target location
	files := []string{"setup.sh", "88-wolf-virtual-inputs.rules", "89-blacklist-vfio.rules"}
	targetDir := "/home/user/"

	// Build the curl command to fetch files from GitHub and place them in the target directory
	curlCommands := make([]string, len(files))
	for i, file := range files {
		url := fmt.Sprintf("https://raw.githubusercontent.com/raefon/td-stream/main/setup/%s", file)
		curlCommands[i] = fmt.Sprintf("curl -o %s%s %s", targetDir, file, url)
	}
	fullCurlCommand := strings.Join(curlCommands, " && ")

	// Set the command to be executed over SSH
	cmd.Flags().Set("command", fullCurlCommand)

	// Execute the SSH command
	return sshServer(cmd, []string{server})
}
