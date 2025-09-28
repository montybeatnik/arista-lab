package arista

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type eosClient struct {
	url        string
	httpClient *http.Client
}

func NewEosClient(url string) eosClient {
	// Configure a custom http.Transport with a modified TLS configuration
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create a custom client with a timeout
	cl := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}
	client := eosClient{url: url, httpClient: cl}
	return client
}

func (c eosClient) Run(reqBody []byte) {
	// Create a new POST request with a body and custom headers
	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewReader(reqBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	username := "admin"
	password := "admin"
	// Set Basic Authentication headers.
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println("Error performing request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// fmt.Println("Body:", string(body))

	var bgpEvpnSummaryResp BGPEvpnSummaryResponse
	if err := json.Unmarshal(body, &bgpEvpnSummaryResp); err != nil {
		panic(err)
	}

	for _, result := range bgpEvpnSummaryResp.Result {
		for _, vrf := range result.Vrfs {
			fmt.Println(vrf.ASN)
			for k, v := range vrf.Peers {
				fmt.Println(k)
				fmt.Printf("prefix rcvd: %v\n", v.PrefixReceived)
				fmt.Printf("prefix adv: %v\n", v.PrefixAdvertised)
				// fmt.Printf("neighbor: %v\n")
			}
		}
	}

}
