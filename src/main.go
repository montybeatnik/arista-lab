package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	arista "github.com/montybeatnik/arista-lab/laber/pkgs"
)

type stringSlice []string

func (s *stringSlice) String() string { return strings.Join(*s, ",") }
func (s *stringSlice) Set(v string) error {
	*s = append(*s, v)
	return nil
}

type PayloadData struct {
	Method  string
	Version int
	Format  string
	Cmds    []string
	ID      int
}

func renderTemplate(tplPath string, data PayloadData) ([]byte, error) {
	funcs := template.FuncMap{
		"toJSON": func(v any) (string, error) {
			b, err := json.Marshal(v)
			return string(b), err
		},
	}
	base := filepath.Base(tplPath)
	tpl, err := template.New(base).Funcs(funcs).ParseFiles(tplPath)
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, base, data); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}
	return buf.Bytes(), nil
}

type eosClient struct {
	url        string
	httpClient *http.Client
}

func newEosClient(url string) eosClient {
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

func (c eosClient) run(reqBody []byte) {
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

	var bgpEvpnSummaryResp arista.BGPEvpnSummaryResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		panic(err)
	}

	vrf := bgpEvpnSummaryResp.Result[0].Vrfs["default"]
	for nbr, p := range vrf.Peers {
		fmt.Printf("Peer %s state=%s rx=%d tx=%d\n", nbr, p.PeerState, p.MsgReceived, p.MsgSent)
	}
}

func main() {
	cmds := []string{"show bgp evpn summary"}
	url := "https://172.20.20.9/command-api"
	client := newEosClient(url)
	tmplPath := "src/templates/eapi_payload.tmpl"
	fmt.Println("rendering template...")
	body, err := renderTemplate(tmplPath, PayloadData{
		Method:  "runCmds",
		Version: 1,
		Format:  "json",
		Cmds:    cmds,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("running cmd...")
	client.run(body)
}
