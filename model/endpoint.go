package model

type Pod struct {
	// Key string
	Name      string
	Namespace string
}

type Allow struct {
	From, To Pod
	// if not specified, allow all ports
	Port int
}

type Policies = []Allow
