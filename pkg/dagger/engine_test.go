package dagger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEngine(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		ctx := context.Background()
		engine, err := NewEngine(ctx)

		require.NoError(t, err)
		assert.NotNil(t, engine)
		assert.Equal(t, ctx, engine.Context())
		assert.False(t, engine.IsClosed())

		// Cleanup
		engine.Close()
	})

	t.Run("nil context", func(t *testing.T) {
		engine, err := NewEngine(nil)

		assert.Error(t, err)
		assert.Nil(t, engine)
		assert.Contains(t, err.Error(), "context cannot be nil")
	})
}

func TestEngine_Container(t *testing.T) {
	ctx := context.Background()
	engine, err := NewEngine(ctx)
	require.NoError(t, err)
	defer engine.Close()

	t.Run("create container builder", func(t *testing.T) {
		container := engine.Container()

		assert.NotNil(t, container)
		assert.Equal(t, engine, container.engine)
		assert.Nil(t, container.err)
	})

	t.Run("container builder on closed engine", func(t *testing.T) {
		engine.Close()
		container := engine.Container()

		assert.NotNil(t, container)
		assert.Error(t, container.err)
		assert.Contains(t, container.err.Error(), "engine is closed")
	})
}

func TestEngine_Host(t *testing.T) {
	ctx := context.Background()
	engine, err := NewEngine(ctx)
	require.NoError(t, err)
	defer engine.Close()

	host := engine.Host()

	assert.NotNil(t, host)
	assert.Equal(t, engine, host.engine)
}

func TestEngine_Close(t *testing.T) {
	ctx := context.Background()
	engine, err := NewEngine(ctx)
	require.NoError(t, err)

	t.Run("close engine", func(t *testing.T) {
		assert.False(t, engine.IsClosed())

		err := engine.Close()
		assert.NoError(t, err)
		assert.True(t, engine.IsClosed())
	})

	t.Run("close already closed engine", func(t *testing.T) {
		err := engine.Close()
		assert.NoError(t, err)
		assert.True(t, engine.IsClosed())
	})
}

func TestEngine_Context(t *testing.T) {
	ctx := context.WithValue(context.Background(), "test", "value")
	engine, err := NewEngine(ctx)
	require.NoError(t, err)
	defer engine.Close()

	engineCtx := engine.Context()

	assert.Equal(t, ctx, engineCtx)
	assert.Equal(t, "value", engineCtx.Value("test"))
}
