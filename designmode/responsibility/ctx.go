package responsibility

// ChainCtx is the context for a chain of responsibility.
type ChainCtx[I any, O any] struct {
	// Request is the request for the chain.
	Request I
	// Response is the response for the chain.
	Response *O
	// Metadata is the metadata for the chain.
	// It is used to store the data for following handlers.
	Metadata map[string]any
}

func NewChainCtx[I any, O any](request I, outboundFactory OutboundFactory[O]) *ChainCtx[I, O] {
	return &ChainCtx[I, O]{
		Request:  request,
		Response: outboundFactory.New(),
		Metadata: make(map[string]any),
	}
}
