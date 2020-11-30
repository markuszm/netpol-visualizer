package model

type Endpoint struct {
	// Key string
	Name      string
	Namespace string
	Port      int
}

type Allow struct {
	e1, e2 Endpoint
}

type Policies struct {
	_ []Allow
}
