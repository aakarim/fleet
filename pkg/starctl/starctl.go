package starctl

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aakarim/fleet/internal/host"
	"github.com/aakarim/fleet/internal/player"
)

var starshPath = flag.String("starsh_path", "./dist/starsh", "the path to the starshell binary")
var privateKeyPath = flag.String("i", "$HOME/.ssh/id_rsa", "the path to your private key")
var PlaybookStr = ""

func Main() {
	flag.Parse()
	ctx := context.Background()

	// check binary exists
	sp, err := os.Open(*starshPath)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			panic("starsh_path does not exist, please build the Starshell and provide a path")
		}

		panic(err)
	}
	sp.Close()

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	hostConnectionStr := os.Args[1]

	if hostConnectionStr == "" {
		panic("you must set a host")
	}

	var pbR io.Reader

	if PlaybookStr == "" {
		pbPath := os.Args[2]
		if pbPath == "" {
			panic("you must set a playbook path")
		}

		f, err := os.Open(pbPath)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		pbR = f
	} else {
		pbR = strings.NewReader(PlaybookStr)
	}

	h, err := host.NewHost(hostConnectionStr, os.ExpandEnv(*privateKeyPath))
	if err != nil {
		panic(err)
	}

	pl := player.NewPlayer(*starshPath)

	if err := pl.Play(ctx, h, pbR); err != nil {
		panic(err)
	}
}
func printUsage() {
	fmt.Print("Usage: starctl [options] <host> <playbook>\n")
	flag.PrintDefaults()
}
