package abstractions

import (
	"fmt"
	"reflect"
)

type System interface {
	Start() error
	Stop() error
	GetComponentByFQDN(name FQDN) (Component, error)
	GetSettings() Settings
	Instance(componentFqdn FQDN) (Component, error)
}

type FQDN string

func TypeToFqdn[T any](t T) FQDN {
	return FQDN(fmt.Sprintf("%T", t))
}

type MySystem struct {
	components map[FQDN]Component
	settings   Settings
}

// Start starts all components in the system
func (ms *MySystem) Start() error {
	for _, component := range ms.components {
		if err := component.Start(); err != nil {
			// TODO does it make sense to stop the components that have already started?
			return err
		}
	}
	return nil
}

func (ms *MySystem) Stop() error {
	errors := make(map[string]error, 0)
	// TODO we could invert the component stop order
	for _, component := range ms.components {
		if err := component.Stop(); err != nil {
			errors[component.Name()] = err
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("errors stopping components: %v", errors)
	}
	return nil
}

func (ms *MySystem) GetComponentByFQDN(name FQDN) (Component, error) {
	if component, ok := ms.components[name]; ok {
		return component, nil
	}
	return nil, fmt.Errorf("component not found")
}

func (ms *MySystem) GetSettings() Settings {
	return ms.settings
}

// Instance creates a new instance of a component without starting it
func (ms *MySystem) Instance(componentFqdn FQDN) (Component, error) {
	if typ, exists := typeRegistry[componentFqdn]; exists {
		if component, ok := ms.components[componentFqdn]; ok {
			return component, nil
		}
		args := []reflect.Value{reflect.ValueOf(ms)}
		newInstance := reflect.ValueOf(typ).Call(args)
		err := newInstance[1].Interface()
		if err != nil {
			return nil, err.(error)
		}
		component := newInstance[0].Interface().(Component)
		ms.components[componentFqdn] = component
		return component, nil
	} else {
		return nil, fmt.Errorf("type not found")
	}
}

func NewMySystem(settings Settings) System {
	return &MySystem{
		components: make(map[FQDN]Component),
		settings:   settings,
	}
}

var typeRegistry = map[FQDN]func(System) (Component, error){}

func RegisterComponentConstructor(fqdn FQDN, constructor func(System) (Component, error)) {
	typeRegistry[fqdn] = constructor
}
