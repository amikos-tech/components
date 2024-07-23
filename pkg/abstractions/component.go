package abstractions

import (
	"fmt"
)

type Component interface {
	GetSystem() System
	Name() string
	Start() error
	Stop() error
	Dependencies() []Component
	Require(components ...FQDN) error
}

type BaseComponent struct {
	name         string
	system       System
	dependencies []Component
}

func (c *BaseComponent) GetSystem() System {
	return c.system
}

// NewBaseComponent is a utility constructor for BaseComponent and should be reused from your own components.
func NewBaseComponent(system System) (Component, error) {
	return &BaseComponent{
		name:         "base",
		system:       system,
		dependencies: make([]Component, 0),
	}, nil
}

func (c *BaseComponent) Name() string {
	return c.name
}

func (c *BaseComponent) Start() error {
	return fmt.Errorf("not implemented")
}

func (c *BaseComponent) Stop() error {
	return fmt.Errorf("not implemented")
}

func (c *BaseComponent) Dependencies() []Component {
	return c.dependencies
}

func (c *BaseComponent) Require(components ...FQDN) error {
	return CommonRequire(&c.dependencies, c.system, components...)
}

func CommonRequire(dependencies *[]Component, system System, components ...FQDN) error {
	if dependencies == nil {
		return fmt.Errorf("dependencies is nil")
	}
	if system == nil {
		return fmt.Errorf("system is nil")
	}

	for _, component := range components {
		instance, err := system.Instance(component)
		if err != nil {
			return err
		}
		*dependencies = append(*dependencies, instance)
	}
	return nil
}

func init() {
	RegisterComponentConstructor(TypeToFqdn(BaseComponent{}), NewBaseComponent)
}
