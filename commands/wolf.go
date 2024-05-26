package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	wolfCmd = &cobra.Command{
		Use:   "wolf",
		Short: "Manage wolf instance",
	}
	wolfLogsCmd = &cobra.Command{
		Use:   "logs server_id",
		Short: "Get wolf logs",
		Args:  cobra.ExactArgs(1), // Expects exactly one argument: server_id
		RunE: func(cmd *cobra.Command, args []string) error {
			// Hardcoded container ID
			containerID := "wolf-wolf-1"
			dockerCommand := "docker logs " + containerID

			// Prepare the arguments for dockerCommandsViaSSH
			dockerArgs := append([]string{args[0]}, dockerCommand)

			// Ensure flags are set for SSH
			cmd.Flags().String("bin", "ssh", "Name of SSH client executable (e.g., ssh, mosh)")
			cmd.Flags().String("user", "user", "User account to use for login")
			cmd.Flags().String("command", "", "Command to execute over SSH")

			// Call dockerCommandsViaSSH with the server ID and the Docker command
			return dockerCommandsViaSSH(cmd, dockerArgs)
		},
	}
	wolfInstallCmd = &cobra.Command{
		Use:   "install server_id",
		Short: "Install wolf on a specified server",
		Args:  cobra.ExactArgs(1), // Expects exactly one argument: server_id
		RunE: func(cmd *cobra.Command, args []string) error {
			serverID := args[0]

			cmd.Flags().String("bin", "ssh", "Name of SSH client executable (e.g., ssh, mosh)")
			cmd.Flags().String("user", "user", "User account to use for login")
			cmd.Flags().String("command", "", "Command to execute over SSH")

			return wolfInstall(cmd, serverID)
		},
	}
)

func init() {
	wolfCmd.AddCommand(wolfLogsCmd)
	wolfCmd.AddCommand(wolfInstallCmd)
	rootCmd.AddCommand(wolfCmd)
}

func sshServer(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	server := args[0]
	bin, err := flags.GetString("bin")
	if err != nil {
		return err
	}

	user, err := flags.GetString("user")
	if err != nil {
		return err
	}

	keyPath, err := flags.GetString("keyPath")
	if err != nil {
		return err
	}
	command, err := flags.GetString("command")
	if err != nil {
		return err
	}

	return executeSSHCommand(server, bin, user, keyPath, command)
}

func executeSSHCommand(serverId, bin, user, keyPath, command string) error {
	res, err := client.GetServer(serverId)
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	// Default SSH port
	sshPort := "22"

	// Check for port forwarding and adjust the SSH port accordingly
	for externalPort, internalPort := range res.VirtualMachines.PortForwards {
		if internalPort == "22" {
			sshPort = externalPort
			break
		}
	}

	sshCmd := exec.Command(bin, "-i", keyPath, "-p", sshPort, fmt.Sprintf("%v@%v", user, res.VirtualMachines.IP), command)
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	if err := sshCmd.Run(); err != nil {
		return err
	}

	return nil
}

func dockerCommandsViaSSH(cmd *cobra.Command, args []string) error {
	server := args[0]
	dockerCommand := strings.Join(args[1:], " ") // Join all arguments after the server ID as the Docker command

	// Set up flags for the SSH command
	cmd.Flags().Set("command", dockerCommand)

	// Call sshServer to handle the SSH connection and command execution
	return sshServer(cmd, []string{server})
}

func getWolfFiles(cmd *cobra.Command, server string) error {
	// Define the files and their target location
	files := []string{"docker-compose.nvidia.yml", "docker-nvidia-start.sh"}
	targetDir := "/home/user/"

	// Build the curl command to fetch files from GitHub and place them in the target directory
	curlCommands := make([]string, len(files))
	for i, file := range files {
		url := fmt.Sprintf("https://raw.githubusercontent.com/raefon/td-stream/main/wolf/%s", file)
		curlCommands[i] = fmt.Sprintf("curl -o %s%s %s", targetDir, file, url)
	}
	fullCurlCommand := strings.Join(curlCommands, " && ")

	// Set the command to be executed over SSH
	cmd.Flags().Set("command", fullCurlCommand)

	// Execute the SSH command
	return sshServer(cmd, []string{server})
}

func wolfInstall(cmd *cobra.Command, server string) error {
	// First, get the wolf files
	err := getWolfFiles(cmd, server)
	if err != nil {
		return fmt.Errorf("error getting wolf files: %w", err)
	}

	// Define the command to run docker-nvidia-start.sh
	startScriptCommand := "bash /home/user/docker-nvidia-start.sh /home/user/docker-compose.nvidia.yml"

	// Set the command to be executed over SSH
	cmd.Flags().Set("command", startScriptCommand)

	// Execute the SSH command to run the script
	return sshServer(cmd, []string{server})
}
