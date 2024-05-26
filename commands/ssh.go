package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

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
