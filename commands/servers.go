package commands

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/raefon/td-stream/api"
	"github.com/spf13/cobra"
)

var (
	serversCmd = &cobra.Command{
		Use:   "servers",
		Short: "Manage servers",
	}
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List servers",
		RunE:  serverList,
	}
	infoCmd = &cobra.Command{
		Use:   "info [flags] server_id",
		Short: "Get server info",
		Args:  cobra.ExactArgs(1),
		RunE:  serverInfo,
	}
	startCmd = &cobra.Command{
		Use:     "start [flags] server_id",
		Short:   "Start a server",
		Args:    cobra.ExactArgs(1),
		RunE:    startServer,
		PostRun: logAction("success"),
	}
	stopCmd = &cobra.Command{
		Use:     "stop [flags] server_id",
		Short:   "Stop a server",
		Args:    cobra.ExactArgs(1),
		RunE:    stopServer,
		PostRun: logAction("success"),
	}
	deleteCmd = &cobra.Command{
		Use:     "delete [flags] server_id",
		Short:   "Delete a server",
		Args:    cobra.ExactArgs(1),
		RunE:    deleteServer,
		PostRun: logAction("success"),
	}
	deployCmd = &cobra.Command{
		Use:     "deploy [flags] name password",
		Short:   "Deploy a server",
		Args:    cobra.ExactArgs(2),
		RunE:    deployServer,
		PostRun: logAction("success"),
	}
	// need to fix
	/*
		manageCmd = &cobra.Command{
			Use:   "manage server_id",
			Short: "Open server management panel in a browser",
			Args:  cobra.ExactArgs(1),
			RunE:  manageServer,
		}
	*/
	// need to implement
	/*
		sshCmd = &cobra.Command{
			Use:   "ssh server_id",
			Short: "Launch an SSH sesion with a server",
			Args:  cobra.ExactArgs(1),
			RunE:  sshServer,
		}
	*/
	restartCmd = &cobra.Command{
		Use:     "restart [flags] server_id",
		Short:   "Restart a server",
		Args:    cobra.ExactArgs(1),
		RunE:    restartServer,
		PostRun: logAction("success"),
	}
	modifyCmd = &cobra.Command{
		Use:     "modify [flags] server_id",
		Short:   "Modify a server",
		Args:    cobra.ExactArgs(1),
		RunE:    modifyServer,
		PostRun: logAction("success"),
	}
	statusCmd = &cobra.Command{
		Use:   "status server_id",
		Short: "Get server status",
		Args:  cobra.ExactArgs(1),
		RunE:  serverStatus,
	}

// need to implement
/*
	dockerCommandsViaSSHCmd = &cobra.Command{
		Use:   "docker-cmd server_id command",
		Short: "Execute a Docker command via SSH for a specified server",
		Args:  cobra.MinimumNArgs(2),
		RunE:  dockerCommandsViaSSH,
	}
	// Wolf stuff
	wolfLogsCmd = &cobra.Command{
		Use:   "wolf_logs server_id",
		Short: "Fetch logs from wolf on a server",
		Args:  cobra.ExactArgs(1), // Expects exactly one argument: server_id
		RunE: func(cmd *cobra.Command, args []string) error {
			// Hardcoded container ID
			containerID := "your_container_id_here"
			dockerCommand := "logs " + containerID

			// Prepare the arguments for dockerCommandsViaSSH
			dockerArgs := append([]string{args[0]}, dockerCommand)

			// Ensure flags are set for SSH
			cmd.Flags().String("bin", "ssh", "Name of SSH client executable (e.g., ssh, mosh)")
			cmd.Flags().String("user", "user", "User account to use for login")
			cmd.Flags().String("keyPath", "~/.ssh/id_rsa", "Path to private key used for authentication")

			// Call dockerCommandsViaSSH with the server ID and the Docker command
			return dockerCommandsViaSSH(cmd, dockerArgs)
		},
	}
*/
)

func init() {
	serversCmd.AddCommand(listCmd)

	serversCmd.AddCommand(infoCmd)

	serversCmd.AddCommand(stopCmd)

	serversCmd.AddCommand(startCmd)

	serversCmd.AddCommand(deleteCmd)

	serversCmd.AddCommand(deployCmd)
	deployCmd.Flags().String("gpuModel", "geforcertx4090-pcie-24gb", "The GPU model that you would like to provision")
	deployCmd.Flags().String("location", "", "Location")
	deployCmd.Flags().String("hostnode", "c136d11f-2ac8-469f-ad2c-8eac05e2a155", "UUID of the hostnode you want to deploy the server on. Can be omitted if location is set.")
	deployCmd.Flags().Int("gpuCount", 1, "The number of GPUs of the model you specified earlier")
	deployCmd.Flags().String("cpuModel", "AMD EPYC 75F3", "The CPU model that you would like to provision")
	deployCmd.Flags().Int("vcpus", 2, "Number of vCPUs that you would like")
	deployCmd.Flags().Int("storage", 20, "Number of GB of networked storage")
	deployCmd.Flags().Int("ram", 4, "Number of GB of RAM to be deployed.")
	deployCmd.Flags().String("operating_system", "Ubuntu 20.04 LTS", "Operating system")
	deployCmd.Flags().String("internal_ports", "80,443", "Internal ports to be used by the server")
	deployCmd.Flags().String("external_ports", "47600,46701", "External ports to be used by the server")

	//need to implement
	/*	serversCmd.AddCommand(manageCmd)

		serversCmd.AddCommand(sshCmd)
		sshCmd.Flags().String("bin", "ssh", "Name of SSH client executable (e.g. ssh, mosh)")
		sshCmd.Flags().String("user", "user", "User account to use for login")
		sshCmd.Flags().String("extraFlags", "", "Extra flags to pass to the SSH client")
		sshCmd.Flags().String("keyPath", "~/.ssh/id_rsa", "Path to private key used for authentication")
	*/
	serversCmd.AddCommand(restartCmd)

	serversCmd.AddCommand(modifyCmd)
	modifyCmd.Flags().String("gpuModel", "Quadro_4000", "The GPU model that you would like to provision")
	modifyCmd.Flags().Int("gpuCount", 1, "The number of GPUs of the model you specified earlier")
	modifyCmd.Flags().String("cpuModel", "Intel_Xeon_v4", "The CPU model that you would like to provision")
	modifyCmd.Flags().Int("vcpus", 2, "Number of vCPUs that you would like")
	modifyCmd.Flags().Int("storage", 20, "Number of GB of networked storage")
	modifyCmd.Flags().Int("ram", 4, "Number of GB of RAM to be deployed.")

	serversCmd.AddCommand(statusCmd)

	//serversCmd.AddCommand(dockerCommandsViaSSHCmd)

	rootCmd.AddCommand(serversCmd)

	// Wolf stuff
	//rootCmd.AddCommand(wolfLogsCmd)
}

func serverList(cmd *cobra.Command, args []string) error {
	res, err := client.ListServers()
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Server ID", "Status"})

	for serverID, details := range res.VirtualMachines {
		serverName := details.Name     // Adjust field names as per actual struct definition
		serverStatus := details.Status // Adjust field names as per actual struct definition
		t.AppendRow(table.Row{serverName, serverID, serverStatus})
	}
	t.Render()

	return nil
}

func serverInfo(cmd *cobra.Command, args []string) error {
	server := args[0]
	res, err := client.GetServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	props := []map[string]string{
		//{"name": "ID", "value": res.VirtualMachines.ServerID},
		{"name": "Name", "value": res.VirtualMachines.Name},
		{"name": "Location", "value": res.VirtualMachines.Location},
		{"name": "HostNode", "value": res.VirtualMachines.HostNode},
		{"name": "IP", "value": res.VirtualMachines.IP},
		{"name": "Charged Cost", "value": fmt.Sprintf("%v", res.VirtualMachines.Cost)},
		{"name": "Status", "value": res.VirtualMachines.Status},
		{"name": "Type", "value": res.VirtualMachines.Type},
		{"name": "vCPUs", "value": strconv.Itoa(res.VirtualMachines.Specs.VCPUs)},
		{"name": "RAM", "value": fmt.Sprintf("%vGB", res.VirtualMachines.Specs.RAM)},
		{"name": "Storage", "value": fmt.Sprintf("%vGB", res.VirtualMachines.Specs.STORAGE)},
		{"name": "Operating System", "value": res.VirtualMachines.OperatingSystem},
		{"name": "Port Forwards", "value": fmt.Sprintf("%v", res.VirtualMachines.PortForwards)},
		{"name": "GPU Amount", "value": strconv.Itoa(res.VirtualMachines.Specs.GPU.Amount)},
		{"name": "GPU Type", "value": res.VirtualMachines.Specs.GPU.Type},
		{"name": "Creation Timestamp", "value": res.VirtualMachines.TimestampCreation},
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Property", "Value"})
	for _, elem := range props {
		t.AppendRow(table.Row{elem["name"], elem["value"]})
	}
	t.Render()

	return nil
}

func startServer(cmd *cobra.Command, args []string) error {
	server := args[0]
	res, err := client.StartServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	return nil
}

func stopServer(cmd *cobra.Command, args []string) error {
	server := args[0]
	res, err := client.StopServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	return nil
}

func deleteServer(cmd *cobra.Command, args []string) error {
	server := args[0]
	res, err := client.DeleteServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	return nil
}

// deployServer deploys a server by making a request to the API with the specified parameters.
func deployServer(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	// Retrieve all parameters and check for their presence
	hostnode, err := flags.GetString("hostnode")
	if err != nil || hostnode == "" {
		return errors.New("hostnode is required")
	}

	gpuModel, err := flags.GetString("gpuModel")
	if err != nil || gpuModel == "" {
		return errors.New("gpuModel is required")
	}

	gpuCount, err := flags.GetInt("gpuCount")
	if err != nil {
		return errors.New("gpuCount is required")
	}

	vcpus, err := flags.GetInt("vcpus")
	if err != nil {
		return errors.New("vcpus is required")
	}

	ram, err := flags.GetInt("ram")
	if err != nil {
		return errors.New("ram is required")
	}

	storage, err := flags.GetInt("storage")
	if err != nil {
		return errors.New("storage is required")
	}

	operatingSystem, err := flags.GetString("operating_system")
	if err != nil || operatingSystem == "" {
		return errors.New("operating_system is required")
	}

	internalPorts, err := flags.GetString("internal_ports")
	if err != nil || internalPorts == "" {
		return errors.New("internal_ports is required")
	}

	internalPortsSlice := strings.Split(internalPorts, ",")

	externalPorts, err := flags.GetString("external_ports")
	if err != nil || externalPorts == "" {
		return errors.New("external_ports is required")
	}

	externalPortsSlice := strings.Split(externalPorts, ",")

	// Initialize the request with all mandatory fields
	req := api.DeployServerRequest{
		HostNode:        hostnode,
		Name:            args[0],
		Password:        args[1],
		GPUModel:        gpuModel,
		GPUCount:        gpuCount,
		VCPUs:           vcpus,
		RAM:             ram,
		Storage:         storage,
		OperatingSystem: operatingSystem,
		InternalPorts:   internalPortsSlice,
		ExternalPorts:   externalPortsSlice,
	}

	res, err := client.DeployServer(req)
	if err != nil {
		return err
	}
	if !res.Success {
		return errors.New(res.Error)
	}

	fmt.Println(res.Server)
	return nil
}

// need to fix
/* func manageServer(cmd *cobra.Command, args []string) error {
	server := args[0]
	res, err := client.GetServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	err = browser.OpenURL(res.Server.Links["dashboard"]["href"])
	if err != nil {
		return err
	}

	return nil
}
*/

// need to fix
/* func executeSSHCommand(serverId, bin, user, keyPath, command string) error {
	res, err := client.GetServer(serverId)
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	sshCmd := exec.Command(bin, "-i", keyPath, fmt.Sprintf("%v@%v", user, res.Server.Ip), command)
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	if err := sshCmd.Run(); err != nil {
		return err
	}

	return nil
}
*/

// need to fix
/*
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

	return executeSSHCommand(server, bin, user, keyPath, "")
}
*/

func logAction(message string) func(*cobra.Command, []string) {
	return func(c *cobra.Command, s []string) { log.Println(message) }
}

func restartServer(cmd *cobra.Command, args []string) error {
	server := args[0]
	res, err := client.RestartServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	return nil
}

func modifyServer(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	serverId := args[0]

	var gpuModel *string = nil
	if flags.Changed("gpuModel") {
		gpuModelVal, err := flags.GetString("gpuModel")
		if err != nil {
			return err
		}
		gpuModel = &gpuModelVal
	}

	var gpuCount *int = nil
	if flags.Changed("gpuCount") {
		gpuCountVal, err := flags.GetInt("gpuCount")
		if err != nil {
			return err
		}
		gpuCount = &gpuCountVal
	}

	var vcpus *int = nil
	if flags.Changed("vcpus") {
		vcpusVal, err := flags.GetInt("vcpus")
		if err != nil {
			return err
		}
		vcpus = &vcpusVal
	}

	var ram *int = nil
	if flags.Changed("ram") {
		ramVal, err := flags.GetInt("ram")
		if err != nil {
			return err
		}
		ram = &ramVal
	}

	var storage *int = nil
	if flags.Changed("storage") {
		storageVal, err := flags.GetInt("storage")
		if err != nil {
			return err
		}
		storage = &storageVal
	}

	req := api.ModifyServerRequest{
		ServerId: serverId,
		VCPUs:    vcpus,
		RAM:      ram,
		Storage:  storage,
	}

	req.GPUModel = gpuModel
	req.GPUCount = gpuCount

	res, err := client.ModifyServer(req)

	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	return nil
}

func serverStatus(cmd *cobra.Command, args []string) error {
	server := args[0]
	res, err := client.GetServerStatus(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	fmt.Println(res.Status)

	return nil
}

// need to fix
/*
func dockerCommandsViaSSH(cmd *cobra.Command, args []string) error {
	server := args[0]
	dockerCommand := strings.Join(args[1:], " ") // Join all arguments after the server ID as the Docker command

	// Set up flags for the SSH command
	cmd.Flags().Set("extraFlags", dockerCommand)

	// Call sshServer to handle the SSH connection and command execution
	return sshServer(cmd, []string{server})
}
*/
