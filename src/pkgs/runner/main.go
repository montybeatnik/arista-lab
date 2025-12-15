package main

func main() {
	cmds := []string{"show bgp summary"}
	url := "https://172.20.20.9/command-api"
	client := arista.NewEosClient(url)
	tmplPath := "templates/eapi_payload.tmpl"
	fmt.Println("rendering template...")
	body, err := renderer.RenderTemplate(tmplPath, arista.PayloadData{
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
	client.Run(body)
}