package systemd

import (
	"github.com/aakarim/fleet/lib/systemd/util"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var Module = &starlarkstruct.Module{
	Name: "systemd",
	Members: starlark.StringDict{
		"util": util.Module,
	},
}
