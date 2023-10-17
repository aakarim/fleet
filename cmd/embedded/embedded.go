package main

import (
	_ "embed"

	"github.com/aakarim/fleet/pkg/starctl"
)

//go:embed embedded.star
var embeddedPlaybook string

func main() {
	starctl.PlaybookStr = embeddedPlaybook
	starctl.Main()
}
