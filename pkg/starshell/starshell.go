package starshell

import (
	"encoding/json"
	"os"

	"github.com/aakarim/fleet/pkg/lib/command"
	starjson "go.starlark.net/lib/json"
	"go.starlark.net/lib/math"
	"go.starlark.net/lib/proto"
	"go.starlark.net/lib/time"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

type Starshell struct {
}

func NewStarshell() *Starshell {
	starsh := &Starshell{}
	return starsh
}

func (s *Starshell) AddModule(name string, module *starlarkstruct.Module) {
	starlark.Universe[name] = module
}

func (s *Starshell) AddStdLib() {
	starlark.Universe["json"] = starjson.Module
	starlark.Universe["time"] = time.Module
	starlark.Universe["math"] = math.Module
	starlark.Universe["proto"] = proto.Module

	starlark.Universe["command"] = command.Module

	starlark.Universe["starshell"] = starlarkstruct.FromStringDict(starlark.String("starshell"), starlark.StringDict{
		"list": starlark.NewBuiltin("list", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			list := s.List()
			starlarkList := starlark.NewList([]starlark.Value{})
			for _, item := range list {
				starlarkList.Append(starlark.String(item))
			}
			return starlarkList, nil
		}),
	})
}

func (s *Starshell) AddPlaybookModules() {
	// starlark.Universe["task"] =
}

func convertToStarlark(args map[string]any) starlark.StringDict {
	// convert to stringdict
	argsDict := starlark.StringDict{}
	for k, v := range args {
		// type switch
		switch vl := v.(type) {
		case string:
			argsDict[k] = starlark.String(vl)
		case int:
			argsDict[k] = starlark.MakeInt(vl)
		case bool:
			argsDict[k] = starlark.Bool(vl)
		case []string:
			// convert to starlark list
			list := starlark.NewList([]starlark.Value{})
			for _, item := range vl {
				list.Append(starlark.String(item))
			}
			argsDict[k] = list
		case []int:
			// convert to starlark list
			list := starlark.NewList([]starlark.Value{})
			for _, item := range vl {
				list.Append(starlark.MakeInt(item))
			}
			argsDict[k] = list
		case []bool:
			// convert to starlark list
			list := starlark.NewList([]starlark.Value{})
			for _, item := range vl {
				list.Append(starlark.Bool(item))
			}
			argsDict[k] = list
		case map[string]string:
			// convert to starlark dict
			dict := starlark.NewDict(0)
			for key, value := range vl {
				dict.SetKey(starlark.String(key), starlark.String(value))
			}
			argsDict[k] = dict
		case map[string]int:
			// convert to starlark dict
			dict := starlark.NewDict(0)
			for key, value := range vl {
				dict.SetKey(starlark.String(key), starlark.MakeInt(value))
			}
			argsDict[k] = dict
		case map[string]bool:
			// convert to starlark dict
			dict := starlark.NewDict(0)
			for key, value := range vl {
				dict.SetKey(starlark.String(key), starlark.Bool(value))
			}
			argsDict[k] = dict

		case map[string][]string:
			// convert to starlark dict
			dict := starlark.NewDict(0)
			for key, value := range vl {
				// convert to starlark list
				list := starlark.NewList([]starlark.Value{})
				for _, item := range value {
					list.Append(starlark.String(item))
				}
				dict.SetKey(starlark.String(key), list)
			}
			argsDict[k] = dict
		case map[string][]int:
			// convert to starlark dict
			dict := starlark.NewDict(0)
			for key, value := range vl {
				// convert to starlark list
				list := starlark.NewList([]starlark.Value{})
				for _, item := range value {
					list.Append(starlark.MakeInt(item))
				}
				dict.SetKey(starlark.String(key), list)
			}
			argsDict[k] = dict
		case map[string][]bool:
			// convert to starlark dict
			dict := starlark.NewDict(0)
			for key, value := range vl {
				// convert to starlark list
				list := starlark.NewList([]starlark.Value{})
				for _, item := range value {
					list.Append(starlark.Bool(item))
				}
				dict.SetKey(starlark.String(key), list)

			}
			argsDict[k] = dict
		default:
			panic("unknown type")

		}
	}
	return argsDict
}

func (starsh *Starshell) Main() int {
	scriptPath := os.Args[1]
	var jsonArgs string
	if len(os.Args) > 2 {
		jsonArgs = os.Args[2]
	}
	if jsonArgs == "" {
		jsonArgs = "{}"
	}
	var args map[string]any

	// unmarshal
	err := json.Unmarshal([]byte(jsonArgs), &args)
	if err != nil {
		panic(err)
	}

	starArgs := convertToStarlark(args)

	thread := &starlark.Thread{}

	starsh.AddStdLib()

	_, err = starlark.ExecFile(thread, scriptPath, nil, starArgs)
	if err != nil {
		panic(err)
	}

	return 0
}
