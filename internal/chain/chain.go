package chain

import (
	"errors"
	"fmt"
)

// Layer represents a named set of environment variables.
type Layer struct {
	Name string
	Vars map[string]string
}

// Chain holds an ordered list of layers. Variables in later layers
// override those defined in earlier layers.
type Chain struct {
	Layers []*Layer
}

// New creates an empty Chain.
func New() *Chain {
	return &Chain{}
}

// AddLayer appends a layer to the chain. Returns an error if a layer
// with the same name already exists.
func (c *Chain) AddLayer(name string, vars map[string]string) error {
	for _, l := range c.Layers {
		if l.Name == name {
			return fmt.Errorf("layer %q already exists in chain", name)
		}
	}
	c.Layers = append(c.Layers, &Layer{Name: name, Vars: vars})
	return nil
}

// RemoveLayer removes a layer by name. Returns an error if not found.
func (c *Chain) RemoveLayer(name string) error {
	for i, l := range c.Layers {
		if l.Name == name {
			c.Layers = append(c.Layers[:i], c.Layers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("layer %q not found", name)
}

// Resolve merges all layers in order and returns the final environment
// variable map. Later layers take precedence.
func (c *Chain) Resolve() (map[string]string, error) {
	if len(c.Layers) == 0 {
		return nil, errors.New("chain has no layers")
	}
	result := make(map[string]string)
	for _, layer := range c.Layers {
		for k, v := range layer.Vars {
			result[k] = v
		}
	}
	return result, nil
}

// GetLayer returns a layer by name or an error if not found.
func (c *Chain) GetLayer(name string) (*Layer, error) {
	for _, l := range c.Layers {
		if l.Name == name {
			return l, nil
		}
	}
	return nil, fmt.Errorf("layer %q not found", name)
}
