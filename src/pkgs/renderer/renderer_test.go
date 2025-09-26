package renderer

import (
	"fmt"
	"testing"

	"github.com/montybeatnik/arista-lab/laber/pkgs/arista"
)

func TestRenderTemplate(t *testing.T) {
	tmplPath := "../../templates/eapi_payload.tmpl"
	cmds := []string{"show bgp evpn summary"}
	payload := arista.PayloadData{
		Method:  "runCmds",
		Version: 1,
		Format:  "json",
		Cmds:    cmds,
	}
	body, err := RenderTemplate(tmplPath, payload)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(body))
}
