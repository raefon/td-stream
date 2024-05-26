package commands

import (
	"github.com/spf13/cobra"
)

var (
	vpnCmd = &cobra.Command{
		Use:   "vpn",
		Short: "Manage vpn instance",
	}
	vpnInstallCmd = &cobra.Command{
		Use:   "install server_id",
		Short: "Install vpn on a specified server",
		Args:  cobra.ExactArgs(1), // Expects exactly one argument: server_id
		RunE: func(cmd *cobra.Command, args []string) error {
			serverID := args[0]

			cmd.Flags().String("bin", "ssh", "Name of SSH client executable (e.g., ssh, mosh)")
			cmd.Flags().String("user", "user", "User account to use for login")
			cmd.Flags().String("command", "", "Command to execute over SSH")

			return vpnInstall(cmd, serverID)
		},
	}
)

func init() {
	vpnCmd.AddCommand(vpnInstallCmd)
	rootCmd.AddCommand(vpnCmd)
}

func vpnInstall(cmd *cobra.Command, server string) error {

	startScriptCommand := "wget https://git.io/wireguard -O wireguard-install.sh && sudo bash wireguard-install.sh"

	// Set the command to be executed over SSH
	cmd.Flags().Set("command", startScriptCommand)

	// Execute the SSH command to run the script
	return sshServer(cmd, []string{server})
}
