package starshell

import "go.starlark.net/starlark"

// List returns all the available modules to this starshell instance.
func (s *Starshell) List() []string {
	var modules []string
	for name := range starlark.Universe {
		modules = append(modules, name)
	}
	return modules
}
