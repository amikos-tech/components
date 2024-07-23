package abstractions

import "fmt"

type Settings interface {
	Get(key string) any
	GetWithDefault(key string, defaultValue any) any
	Require(key string) (any, error)
	FromEnv()
}

type BaseSettings struct {
	settings map[string]any
}

func (bs *BaseSettings) Require(key string) (any, error) {
	if val, ok := bs.settings[key]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("required setting [%s] not found", key)
}

func (bs *BaseSettings) Get(key string) any {
	return bs.settings[key]
}

func (bs *BaseSettings) FromEnv() {
	// Read settings from environment
}

func (bs *BaseSettings) GetWithDefault(key string, defaultValue any) any {
	if val, ok := bs.settings[key]; ok {
		return val
	}
	return defaultValue
}

func NewBaseSettings(keysAndValues ...interface{}) Settings {
	s := &BaseSettings{
		settings: map[string]any{},
	}
	for i := 0; i < len(keysAndValues); i += 2 {
		s.settings[keysAndValues[i].(string)] = keysAndValues[i+1]
	}
	return s
}
