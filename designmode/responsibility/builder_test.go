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
	name          string
	handleSuccess bool
	shouldErr     bool
	stopChain     bool
	canHandle     bool
}

func NewMockHandler(name string, shouldErr, stopChain bool) *MockHandler {
	return &MockHandler{
		name:      name,
		shouldErr: shouldErr,
		stopChain: stopChain,
		canHandle: true, // By default handlers can handle requests
	}
}

func (h *MockHandler) Name() string {
	return h.name
}

func (h *MockHandler) CanHandle(ctx *ChainCtx[string, string]) bool {
	return h.canHandle
}

func (h *MockHandler) Handle(ctx *ChainCtx[string, string]) (bool, error) {
	if h.shouldErr {
		return h.stopChain, errors.New("mock error")
	}
	h.handleSuccess = true
	*ctx.Response += h.name
	return h.stopChain, nil
}

func (h *MockHandler) Rollback(ctx *ChainCtx[string, string]) {
	h.handleSuccess = false
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
	assert.True(t, h1.handleSuccess)
	assert.True(t, h2.handleSuccess)
	assert.True(t, h3.handleSuccess)
}

// nolint: dupl
func TestBuilder_Execute_SkipHandler(t *testing.T) {
	builder := NewBuilder("input", &MockOutboundFactory{})

	h1 := NewMockHandler("h1", false, false)
	h2 := NewMockHandler("h2", false, false)
	h3 := NewMockHandler("h3", false, false)

	// h2 cannot handle the request
	h2.canHandle = false

	builder.Link(h1).Link(h2).Link(h3)

	ctx, err := builder.Execute()

	assert.NoError(t, err)
	assert.Equal(t, "h1h3", *ctx.Response)
	assert.True(t, h1.handleSuccess)
	assert.False(t, h2.handleSuccess)
	assert.True(t, h3.handleSuccess)
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
	assert.True(t, h1.handleSuccess)
	assert.True(t, h2.handleSuccess)
	assert.False(t, h3.handleSuccess)
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
	assert.False(t, h3.handleSuccess)
	// After rollback
	assert.False(t, h1.handleSuccess)
	assert.False(t, h2.handleSuccess)
}

func TestBuilder_Execute_ErrorWithoutRollback(t *testing.T) {
	builder := NewBuilder("input", &MockOutboundFactory{})

	h1 := NewMockHandler("h1", false, false)
	h2 := NewMockHandler("h2", true, false) // This handler returns error
	h3 := NewMockHandler("h3", false, false)

	builder.Link(h1).Link(h2).Link(h3)

	ctx, err := builder.Execute()

	assert.NoError(t, err)
	assert.Contains(t, ctx.Metadata, "h2_error")
	assert.True(t, h1.handleSuccess)
	assert.False(t, h2.handleSuccess)
	assert.True(t, h3.handleSuccess)
}

func TestBuilder_Execute_ErrorWithSelectiveRollback(t *testing.T) {
	builder := NewBuilder("input", &MockOutboundFactory{}, RollbackOnError[string, string](true))

	h1 := NewMockHandler("h1", false, false)
	h2 := NewMockHandler("h2", false, false)
	h3 := NewMockHandler("h3", true, false) // This handler returns error
	h4 := NewMockHandler("h4", false, false)

	// h2 cannot handle the request
	h2.canHandle = false

	builder.Link(h1).Link(h2).Link(h3).Link(h4)

	ctx, err := builder.Execute()

	assert.Error(t, err)
	assert.Nil(t, ctx)

	// Verify that h2 was not handled and not rolled back (because canHandle = false)
	assert.False(t, h2.handleSuccess)

	// Verify that h1 and h3 were handled and rolled back
	assert.False(t, h1.handleSuccess) // Rolled back from true
	assert.False(t, h3.handleSuccess) // Rolled back from true

	// Verify that h4 was never handled (due to error in h3)
	assert.False(t, h4.handleSuccess)
}
