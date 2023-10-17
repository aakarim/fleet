package command

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var Module = &starlarkstruct.Module{
	Name: "command",
	Members: starlark.StringDict{
		"run": starlark.NewBuiltin("run", Run),
	},
}

func Run(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name, commandArgs starlark.Value

	if err := starlark.UnpackPositionalArgs("run", args, kwargs, 1, &name, &commandArgs); err != nil {
		return nil, err
	}

	var commandArgsString []string

	// check to see if we have the command, if not then expand using shell
	path, err := exec.LookPath(name.String())
	if err != nil {
		if !errors.Is(err, exec.ErrNotFound) {
			return nil, fmt.Errorf("command %s not found; if this is a shell command did you try adding your shell e.g. /bin/sh?", name.String())
		}

		path = os.ExpandEnv("$SHELL")

		commandArgsString = append(commandArgsString, "-c", name.String())
	}

	if commandArgs != nil {
		if commandArgs.Type() != "list" {
			return nil, fmt.Errorf("command args must be a list")
		}

		it := commandArgs.(*starlark.List).Iterate()
		defer it.Done()

		var currentVal starlark.Value
		for it.Next(&currentVal) {
			if currentVal.Type() != "string" {
				return nil, fmt.Errorf("command args must be a list of strings")
			}
			commandArgsString = append(commandArgsString, currentVal.String())
		}
	}

	cmd := exec.Command(path, commandArgsString...)
	out, err := cmd.Output()
	if err != nil {
		return starlark.None, err
	}

	print(string(out))

	return starlark.None, nil
}
