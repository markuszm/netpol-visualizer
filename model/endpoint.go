package model

import v1 "k8s.io/api/networking/v1"

type Pod struct {
	// Key string
	Name      string
	Namespace string
}

type PolicyEdge struct {
	PolicyType v1.PolicyType
	From, To Pod
	// if not specified, allow all ports
	Port int
}

type Policies = []PolicyEdge
