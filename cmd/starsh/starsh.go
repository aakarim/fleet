package main

import (
	"os"

	"github.com/aakarim/fleet/lib/systemd"
	"github.com/aakarim/fleet/pkg/starshell"
)

func main() {
	starsh := starshell.NewStarshell()

	// example module addition:
	starsh.AddModule("systemd", systemd.Module)

	os.Exit(starsh.Main())
}
