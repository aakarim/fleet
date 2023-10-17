package playbook

type Playbook struct {
	Name  string
	Hosts string
	Tasks []Task
}

type Action struct {
}

type Task struct {
	Name   string
	Action Action
}
