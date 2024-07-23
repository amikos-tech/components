package abstractions

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSystem(t *testing.T) {
	sys := NewMySystem(NewBaseSettings())
	instance, err := sys.Instance(TypeToFqdn(BaseComponent{}))
	require.NoError(t, err, "Error creating instance: %v", err)
	require.NotNil(t, instance, "Instance should not be nil")
	require.Equal(t, "base", instance.Name(), "Instance name should be 'base'")
}
