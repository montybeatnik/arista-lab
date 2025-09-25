package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func newClient() *http.Client {
        // Configure a custom http.Transport with a modified TLS configuration
        tr := &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }

        // Create a custom client with a timeout
        client := &http.Client{
                Timeout: 10 * time.Second,
                Transport: tr,
        }
	return client
}

func runCmd(cmd string, client *http.Client) {
	fmt.Println("running %v", cmd)
        // Create a new POST request with a body and custom headers
        jsonBody := `{
        "jsonrpc": "2.0",
        "method": "runCmds",
        "params": {
          "version": 1,
          "format": "json",
          "cmds": ["show version"]
        },
        "id": 1
        }`
        req, err := http.NewRequest(http.MethodPost, "https://172.20.20.9/command-api", strings.NewReader(jsonBody))
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
        resp, err := client.Do(req)
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
        fmt.Println("Body:", string(body))
}

func main() {

	client := newClient()
        runCmd("show version", client)
}

