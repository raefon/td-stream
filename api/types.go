package api

type VirtualMachine struct {
	Cost            float32           `json:"cost"`
	Location        string            `json:"location"`
	HostNode        string            `json:"hostnode"`
	Name            string            `json:"name"`
	OperatingSystem string            `json:"operating_system"`
	PortForwards    map[string]string `json:"port_forwards"`
	IP              string            `json:"ip_address"`
	Type            string            `json:"type"`

	Specs struct {
		GPU struct {
			Amount int    `json:"amount"`
			Type   string `json:"type"`
		} `json:"gpu"`
		RAM     int `json:"ram"`
		VCPUs   int `json:"vcpus"`
		STORAGE int `json:"storage"`
	} `json:"specs"`
	Status            string `json:"status"`
	TimestampCreation string `json:"timestamp_creation"`
}

type BillingDetails struct {
	Balance            float32 `json:"balance"`
	HourlySpendingRate float32 `json:"hourly_spending_rate"`
}
