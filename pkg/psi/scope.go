package psi

type Scope interface {
	GetItems() map[string]ScopeItem
}

type ScopeItem interface {
	GetName() string
	GetType() string
	// TODO: get icon or something
}
