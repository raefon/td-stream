package commands

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	stockCmd = &cobra.Command{
		Use:   "stock",
		Short: "Query stock",
	}
	listStockCmd = &cobra.Command{
		Use:   "list",
		Short: "List stock",
		RunE:  listStock,
	}
)

func init() {
	listStockCmd.Flags().Bool("all", false, "Include out-of-stock instances")
	stockCmd.AddCommand(listStockCmd)
	rootCmd.AddCommand(stockCmd)
}

func listStock(cmd *cobra.Command, args []string) error {
	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	res, err := client.ListStock()
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	t.AppendHeader(table.Row{"HostNode ID", "GPU", "Region", "Available Units", "GPU Price", "Location", "External Ports"})

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	// Iterate over HostNode map
	for hostID, host := range res.HostNode {
		location := host.Location
		region := location.Region
		city := location.City
		locationStr := fmt.Sprintf("%s, %s", city, region)
		networking := host.Networking
		ports := networking.Ports

		// Sort ports to easily find min and max
		var portsStr string
		random.Shuffle(len(ports), func(i, j int) { ports[i], ports[j] = ports[j], ports[i] })
		if len(ports) > 10 {
			ports = ports[:10] // Slice the first 10 ports after shuffling
		}
		if len(ports) > 0 {
			portsStr = fmt.Sprintf("%v", ports)
		} else {
			portsStr = "No ports"
		}

		specs := host.Specs
		gpus := specs.GPU
		for gpuName, gpuDetails := range gpus {
			amount := gpuDetails.Amount
			price := gpuDetails.Price
			if amount > 0 || all {
				t.AppendRow(table.Row{hostID, gpuName, region, amount, price, locationStr, portsStr})
			}
		}
	}

	// Sort
	t.SortBy([]table.SortBy{
		//{Name: "GPU", Mode: table.Asc},
		//{Name: "Region", Mode: table.Asc},
		{Name: "GPU Price", Mode: table.Asc},
	})

	t.Render()
	return nil
}
