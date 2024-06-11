package commands

import (
	"github.com/spf13/cobra"
)

var (
	nvidiaCmd = &cobra.Command{
		Use:   "nvidia",
		Short: "Manage nvidia drivers",
	}
	nvidiaInstallCmd = &cobra.Command{
		Use:   "install server_id",
		Short: "Install nvidia drivers on a specified server",
		Args:  cobra.ExactArgs(1), // Expects exactly one argument: server_id
		RunE: func(cmd *cobra.Command, args []string) error {
			serverID := args[0]

			cmd.Flags().String("bin", "ssh", "Name of SSH client executable (e.g., ssh, mosh)")
			cmd.Flags().String("user", "user", "User account to use for login")
			cmd.Flags().String("command", "", "Command to execute over SSH")

			return nvidiaInstall(cmd, serverID)
		},
	}
)

func init() {
	nvidiaCmd.AddCommand(nvidiaInstallCmd)
	rootCmd.AddCommand(nvidiaCmd)
}

func nvidiaInstall(cmd *cobra.Command, server string) error {

	startScriptCommand := "sudo add-apt-repository ppa:graphics-drivers/ppa -y && sudo apt update && sudo apt install nvidia-driver-550 -y && echo 'Complete.... Rebooting.' && sudo reboot"

	// Set the command to be executed over SSH
	cmd.Flags().Set("command", startScriptCommand)

	// Execute the SSH command to run the script
	return sshServer(cmd, []string{server})
}
