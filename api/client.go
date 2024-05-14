package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
)

var CLIENT_VERSION = "0.7.1"

type Response struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type GetServerResponse struct {
	Response
	VirtualMachines VirtualMachine `json:"virtualmachines"`
}

type GetServerStatusResponse struct {
	Response
	Status string `json:"status"`
}

type ListServersResponse struct {
	Response
	VirtualMachines map[string]VirtualMachine `json:"virtualmachines"`
}

type GetBillingDetailsResponse struct {
	Response
	BillingDetails
}

type DeployServerRequest struct {
	Password        string   `mapstructure:"password"`
	GPUModel        string   `mapstructure:"gpu_model"`
	GPUCount        int      `mapstructure:"gpu_count"`
	VCPUs           int      `mapstructure:"vcpus"`
	RAM             int      `mapstructure:"ram"`
	Storage         int      `mapstructure:"storage"`
	OperatingSystem string   `mapstructure:"operating_system"`
	Location        string   `mapstructure:"location,omitempty"`
	HostNode        string   `mapstructure:"hostnode"`
	Name            string   `mapstructure:"name"`
	InternalPorts   []string `mapstructure:"internal_ports"`
	ExternalPorts   []string `mapstructure:"external_ports"`
}

type ModifyServerRequest struct {
	ServerId string  `mapstructure:"server_id"`
	GPUModel *string `mapstructure:"gpu_model,omitempty"`
	GPUCount *int    `mapstructure:"gpu_count,omitempty"`
	CPUModel *string `mapstructure:"cpu_model,omitempty"`
	VCPUs    *int    `mapstructure:"vcpus"`
	RAM      *int    `mapstructure:"ram"`
	Storage  *int    `mapstructure:"storage"`
}

type DeployServerResponse struct {
	Response
	Cost struct {
		ComputePrice float64 `json:"compute_price"`
		StoragePrice float64 `json:"storage_price"`
		TotalPrice   float64 `json:"total_price"`
	} `json:"cost"`
	IP           string            `json:"ip"`
	PortForwards map[string]string `json:"port_forwards"`
	Server       string            `json:"server"`
}

type ListStockResponse struct {
	Response
	HostNode map[string]struct {
		Location struct {
			City    string `json:"city"`
			Country string `json:"country"`
			Region  string `json:"region"`
		} `json:"location"`
		Networking struct {
			Ports []int `json:"ports"`
		} `json:"networking"`
		Specs struct {
			CPU struct {
				Amount int `json:"amount"`
				Name   string
				Price  float64 `json:"price"`
			} `json:"cpu"`
			GPU map[string]struct {
				Amount int `json:"amount"`
				Name   string
				Price  float64 `json:"price"`
			} `json:"gpu"`
			RAM struct {
				Amount int     `json:"amount"`
				Price  float64 `json:"price"`
			} `json:"ram"`
			Storage struct {
				Amount int     `json:"amount"`
				Price  float64 `json:"price"`
			} `json:"storage"`
		} `json:"specs"`
	} `json:"hostnodes"`
}

type Client struct {
	BaseUrl  string
	ApiKey   string
	ApiToken string
	Debug    bool
	KeyPath  string
}

func (client *Client) do(method string, path string, params map[string]string, headers map[string]string, body []byte) (*json.RawMessage, error) {
	query := url.Values{}
	for key, elem := range params {
		query.Add(key, elem)
	}

	url := fmt.Sprintf("%v/%v?%v", client.BaseUrl, path, query.Encode())
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for key, elem := range headers {
		req.Header.Add(key, elem)
	}

	if client.Debug {
		reqDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		fmt.Println(string(reqDump))
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if client.Debug {
		resDump, err := httputil.DumpResponse(res, true)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		fmt.Println(string(resDump))
	}

	bytes, _ := io.ReadAll(res.Body)

	// HACK: Workaround for API issue which causes endpoint to
	// return an HTML Page with a 200 Status code
	if strings.HasPrefix(res.Header.Get("Content-Type"), "text/html") {
		bytes, err := json.Marshal(Response{Success: false, Error: "api call failed"})
		if err != nil {
			return nil, err
		}
		msg := json.RawMessage(bytes)
		return &msg, nil
	}

	var raw map[string]interface{}
	err = json.Unmarshal(bytes, &raw)
	if err != nil {
		return nil, err
	}

	// HACK: Some endpoints return a string boolean on the `success`` field
	if val, ok := raw["success"]; ok {
		if reflect.ValueOf(val).Kind() == reflect.String {
			success, err := strconv.ParseBool(val.(string))
			if err != nil {
				success = false
			}
			raw["success"] = success
		}
	}

	// HACK: on some endpoints, the success key is only
	// present on OK response codes, if the response code
	// is good then just assume the value is true
	if _, ok := raw["success"]; !ok {
		if res.StatusCode >= 200 && res.StatusCode <= 300 {
			raw["success"] = true
		}
	}

	bytes, err = json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	msg := json.RawMessage(bytes)
	return &msg, nil
}

func (client *Client) get(path string, params map[string]string, auth bool) (*json.RawMessage, error) {
	newParams := map[string]string{}

	if auth {
		newParams["api_key"] = client.ApiKey
		newParams["api_token"] = client.ApiToken
	}

	for key, elem := range params {
		newParams[key] = elem
	}

	headers := map[string]string{}
	headers["User-Agent"] = fmt.Sprintf("td-stream/%v", CLIENT_VERSION)

	return client.do(http.MethodGet, path, newParams, headers, nil)
}

func (client *Client) post(path string, body map[string]string, auth bool) (*json.RawMessage, error) {
	newBody := url.Values{}

	if auth {
		newBody.Add("api_key", client.ApiKey)
		newBody.Add("api_token", client.ApiToken)
	}

	for key, elem := range body {
		newBody.Add(key, elem)
	}

	headers := map[string]string{}
	headers["User-Agent"] = fmt.Sprintf("td-stream/%v", CLIENT_VERSION)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	// Encode the form data and replace '+' with '%20'
	encodedData := newBody.Encode()
	correctedData := strings.ReplaceAll(encodedData, "+", "%20")

	return client.do(
		http.MethodPost,
		path,
		nil,
		headers,
		[]byte(correctedData),
	)
}

func (client *Client) ListServers() (*ListServersResponse, error) {
	raw, err := client.post("list", nil, true)
	if err != nil {
		return nil, err
	}

	var res ListServersResponse
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *Client) StopServer(server string) (*Response, error) {
	raw, err := client.get("stop/single", map[string]string{"server": server}, true)
	if err != nil {
		return nil, err
	}

	var res Response
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *Client) StartServer(server string) (*Response, error) {
	raw, err := client.get("start/single", map[string]string{"server": server}, true)
	if err != nil {
		return nil, err
	}

	var res Response
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *Client) DeleteServer(server string) (*Response, error) {
	raw, err := client.get("delete/single", map[string]string{"server": server}, true)
	if err != nil {
		return nil, err
	}

	var res Response
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *Client) GetServer(server string) (*GetServerResponse, error) {
	raw, err := client.post("get/single", map[string]string{"server": server}, true)
	if err != nil {
		return nil, err
	}

	var res GetServerResponse
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *Client) DeployServer(req DeployServerRequest) (*DeployServerResponse, error) {
	body := url.Values{} // Define 'body' here

	var rawBody map[string]interface{}
	err := mapstructure.Decode(req, &rawBody)
	if err != nil {
		return nil, err
	}

	for key, elem := range rawBody {
		val := reflect.ValueOf(elem)
		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				continue // Skip nil pointer fields
			}
			val = val.Elem() // Dereference the pointer to get the value
		}
		if (key == "internal_ports" || key == "external_ports") && val.Kind() == reflect.Slice {
			// Special handling for ports if it is a slice
			slice := make([]string, val.Len())
			for i := 0; i < val.Len(); i++ {
				slice[i] = fmt.Sprintf("%v", val.Index(i).Interface())
			}
			body.Add(key, "{"+strings.Join(slice, ", ")+"}")
		} else if val.Kind() == reflect.Slice {
			// General handling for other slices
			slice := make([]string, val.Len())
			for i := 0; i < val.Len(); i++ {
				slice[i] = fmt.Sprintf("%v", val.Index(i).Interface())
			}
			body.Add(key, strings.Join(slice, ", "))
		} else {
			body.Add(key, fmt.Sprintf("%v", val))
		}
	}

	// Add authentication parameters
	body.Add("api_key", client.ApiKey)
	body.Add("api_token", client.ApiToken)

	headers := map[string]string{
		"User-Agent":   fmt.Sprintf("td-stream/%v", CLIENT_VERSION),
		"Content-Type": "application/x-www-form-urlencoded",
	}

	// Encode the form data and replace '+' with '%20'
	encodedData := body.Encode()
	correctedData := strings.ReplaceAll(encodedData, "+", "%20")

	raw, err := client.do(
		http.MethodPost,
		"deploy/single",
		nil,
		headers,
		[]byte(correctedData),
	)
	if err != nil {
		return nil, err
	}

	var res DeployServerResponse
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *Client) GetBillingDetails() (*GetBillingDetailsResponse, error) {
	raw, err := client.get("billing", nil, true)
	if err != nil {
		return nil, err
	}

	var res GetBillingDetailsResponse
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func NewClient(baseUrl string, apiKey string, apiToken string, debug bool, keyPath string) *Client {
	return &Client{baseUrl, apiKey, apiToken, debug, keyPath}
}

func (client *Client) RestartServer(server string) (*Response, error) {
	raw, err := client.get("restart/single", map[string]string{"server": server}, true)
	if err != nil {
		return nil, err
	}

	var res Response
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *Client) ListStock() (*ListStockResponse, error) {
	raw, err := client.get("deploy/hostnodes", nil, false)
	if err != nil {
		return nil, err
	}

	var res ListStockResponse
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	for _, host := range res.HostNode {
		for gpuName, gpuDetails := range host.Specs.GPU {
			gpuDetails.Name = gpuName
			host.Specs.GPU[gpuName] = gpuDetails
		}
	}

	return &res, nil
}

func (client *Client) ModifyServer(req ModifyServerRequest) (*Response, error) {
	var rawBody map[string]interface{}
	err := mapstructure.Decode(req, &rawBody)
	if err != nil {
		return nil, err
	}

	// convert to map[string]string skipping nil pointers
	body := map[string]string{}
	for key, elem := range rawBody {
		val := reflect.ValueOf(elem)
		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				continue
			}
			val = val.Elem()
		}
		body[key] = fmt.Sprintf("%v", val)
	}

	raw, err := client.post("modify/single", body, true)
	if err != nil {
		return nil, err
	}

	var res Response
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *Client) GetServerStatus(server string) (*GetServerStatusResponse, error) {
	raw, err := client.post("deploy/status", map[string]string{"server": server}, true)
	if err != nil {
		return nil, err
	}

	var res GetServerStatusResponse
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
