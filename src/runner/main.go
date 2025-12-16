package main

import (
	"fmt"
	"os"
	"github.com/montybeatnik/arista-lab/laber/pkgs/arista"
)

func main() {
	cmds := []string{"show bgp summary"}
	url := "https://172.20.20.9/command-api"
	client := arista.NewEosClient(url)
	bgpEvpnSummaryResp, err := client.BGPSummary()
	if err != nil {
		fmt.Printf("Run failed: %v\n", err)
	}
	fmt.Println(bgpEvpnSummaryResp)
}