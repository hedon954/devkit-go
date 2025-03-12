package responsibility

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockOutboundFactory is a mock implementation of OutboundFactory
type MockOutboundFactory struct{}

func (f *MockOutboundFactory) New() *string {
	s := ""
	return &s
}

// MockHandler is a mock implementation of Handler
type MockHandler struct {
	HandlerBase[string, string]
	name      string
	handled   bool
	shouldErr bool
	stopChain bool
}

func NewMockHandler(name string, shouldErr, stopChain bool) *MockHandler {
	return &MockHandler{
		name:      name,
		shouldErr: shouldErr,
		stopChain: stopChain,
	}
}

func (h *MockHandler) Name() string {
	return h.name
}

func (h *MockHandler) Handle(ctx *ChainCtx[string, string]) (bool, error) {
	h.handled = true
	if h.shouldErr {
		return h.stopChain, errors.New("mock error")
	}
	*ctx.Response += h.name
	return h.stopChain, nil
}

func (h *MockHandler) Rollback(ctx *ChainCtx[string, string]) {
	h.handled = false
	h.HandlerBase.Rollback(ctx)
}

// nolint: dupl
func TestBuilder_Execute_Success(t *testing.T) {
	builder := NewBuilder("input", &MockOutboundFactory{})

	h1 := NewMockHandler("h1", false, false)
	h2 := NewMockHandler("h2", false, false)
	h3 := NewMockHandler("h3", false, false)

	builder.Link(h1).Link(h2).Link(h3)

	ctx, err := builder.Execute()

	assert.NoError(t, err)
	assert.Equal(t, "h1h2h3", *ctx.Response)
	assert.True(t, h1.handled)
	assert.True(t, h2.handled)
	assert.True(t, h3.handled)
}

// nolint: dupl
func TestBuilder_Execute_StopChain(t *testing.T) {
	builder := NewBuilder("input", &MockOutboundFactory{})

	h1 := NewMockHandler("h1", false, false)
	h2 := NewMockHandler("h2", false, true) // This handler stops the chain
	h3 := NewMockHandler("h3", false, false)

	builder.Link(h1).Link(h2).Link(h3)

	ctx, err := builder.Execute()

	assert.NoError(t, err)
	assert.Equal(t, "h1h2", *ctx.Response)
	assert.True(t, h1.handled)
	assert.True(t, h2.handled)
	assert.False(t, h3.handled)
}

// nolint: dupl
func TestBuilder_Execute_ErrorWithRollback(t *testing.T) {
	builder := NewBuilder("input",
		&MockOutboundFactory{},
		RollbackOnError[string, string](true))

	h1 := NewMockHandler("h1", false, false)
	h2 := NewMockHandler("h2", true, false) // This handler returns error
	h3 := NewMockHandler("h3", false, false)

	builder.Link(h1).Link(h2).Link(h3)

	ctx, err := builder.Execute()

	assert.Error(t, err)
	assert.Nil(t, ctx)
	assert.False(t, h3.handled)
	// After rollback
	assert.False(t, h1.handled)
	assert.False(t, h2.handled)
}
