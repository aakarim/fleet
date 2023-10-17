# Motivations & Goals
Build a simple and easy to use configuration management tool that has:

- small deployment footprint, no requirement for a full python environment on the host. Agentless
- the capability to easily package up playbooks as binaries so they can be shared with others
- a secure, sandboxed environment
- testable and composable playbooks
- a straightforward onboarding process

## Why Starlark?
- Good syntax for both declarative and imperative programming.
- Composable. Playbooks in YAML resulted in large files that were hard to manage. It was hard to 'import' other playbooks and run them all as one group. With Starlark you can keep file sizes small and manageable, then import other playbooks and run them as one group.
- Import Go code and use it in your playbooks.

# Components

## Starctl
Starctl is the CLI tool that is used to run playbooks on hosts and manage Starshells.

## Starshell
Starshell is the binary that is deployed on a host machine, it is the environment on which every Playbook is run.

## Playbooks
Playbooks are configuration files that define what series of tasks are run on a server to configure it.

Check out the examples/ directory for some examples.


# How to build your Starshell
A Starshell is a binary that contains all the dependencies needed to run your
scripts. It is a Go Binary that can be copied and run on any machine. Check out cmd/starshell for an example.

## Repl
You can run the Starshell as a Repl so you can test out modules.

# Embedding Playbooks
You can embed playbooks in your Starshell binary. This is useful if you want to share your Starshell and playbooks with others. Check out the example in cmd/embedded