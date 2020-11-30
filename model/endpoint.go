package model

type Pod struct {
	// Key string
	Name      string
	Namespace string
}

type Allow struct {
	from, to Pod
	// if not specified, allow all ports
	port int
}

type Policies = []Allow
