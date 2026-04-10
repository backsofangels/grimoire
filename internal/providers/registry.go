package providers

import "fmt"

var registry = map[string]Provider{}

func Register(p Provider) {
	registry[p.Name()] = p
}

func Get(name string) (Provider, error) {
	if p, ok := registry[name]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("provider not found: %s", name)
}

func All() []Provider {
	out := make([]Provider, 0, len(registry))
	for _, p := range registry {
		out = append(out, p)
	}
	return out
}
