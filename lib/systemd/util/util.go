package util

import (
	"github.com/coreos/go-systemd/util"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var Module = &starlarkstruct.Module{
	Name: "util",
	Members: starlark.StringDict{
		"get_machine_id": starlark.NewBuiltin("get_machine_id", GetMachineID),
	},
}

func GetMachineID(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	machineID, err := util.GetMachineID()
	if err != nil {
		return nil, err
	}

	return starlark.String(machineID), nil
}
