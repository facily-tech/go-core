package telemetry

type ddSpan struct {
	context ddSpanContext
}

type ddSpanContext struct {
	spanID  uint64
	traceID uint64
}

func (dd *ddSpan) Context() SpanContext {
	return &dd.context
}

func (dd *ddSpanContext) SpanID() uint64 {
	return dd.spanID
}
func (dd *ddSpanContext) TraceID() uint64 {
	return dd.traceID
}

func (dd *ddSpanContext) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"traceID": dd.traceID,
		"spanID":  dd.spanID,
	}
}
