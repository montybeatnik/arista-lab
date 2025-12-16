package arista

func (c eosClient) BGPSummary() (BGPEvpnSummaryResponse, error) {
	cmds := []string{"show bgp summary"}
	tmplPath := "templates/eapi_payload.tmpl"
	fmt.Println("rendering template...")
	body, err := renderer.RenderTemplate(tmplPath, arista.PayloadData{
		Method:  "runCmds",
		Version: 1,
		Format:  "json",
		Cmds:    cmds,
	})
	if err != nil {
		fmt.Printf("failed to render template: %v\n", err)
		return bgpEvpnSummaryResp{}, fmt.Errorf("failed to render template: %v", err)
	}
	fmt.Println("running cmd...")
	var bgpEvpnSummaryResp arista.BGPEvpnSummaryResponse
	if err := client.Run(body, &bgpEvpnSummaryResp); err != nil {
		fmt.Printf("Run failed: %v\n", err)
		return bgpEvpnSummaryResp{}, fmt.Errorf("run failed: %v", err)
	}
	return bgpEvpnSummaryResp, nil
}