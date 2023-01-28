package trace

import (
	"context"
	"github.com/lackone/gin-ext/framework"
	"github.com/lackone/gin-ext/framework/contract"
	"github.com/lackone/gin-ext/framework/gin"
	"net/http"
	"time"
)

type TraceKey string

var ContextKey = TraceKey("trace-key")

type ExtTrace struct {
	idService contract.IDService

	traceIDGenerator contract.IDService
	spanIDGenerator  contract.IDService
}

func NewExtTrace(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	id := c.MustMake(contract.IdKey).(contract.IDService)
	return &ExtTrace{idService: id}, nil
}

// WithTrace register new trace to context
func (e *ExtTrace) WithTrace(c context.Context, trace *contract.TraceContext) context.Context {
	if ginC, ok := c.(*gin.Context); ok {
		ginC.Set(string(ContextKey), trace)
		return ginC
	} else {
		newC := context.WithValue(c, ContextKey, trace)
		return newC
	}
}

// GetTrace From trace context
func (e *ExtTrace) GetTrace(c context.Context) *contract.TraceContext {
	if ginC, ok := c.(*gin.Context); ok {
		if val, ok2 := ginC.Get(string(ContextKey)); ok2 {
			return val.(*contract.TraceContext)
		}
	}

	if tc, ok := c.Value(ContextKey).(*contract.TraceContext); ok {
		return tc
	}
	return nil
}

// NewTrace generate a new trace
func (e *ExtTrace) NewTrace() *contract.TraceContext {
	var traceID, spanID string
	if e.traceIDGenerator != nil {
		traceID = e.traceIDGenerator.NewID()
	} else {
		traceID = e.idService.NewID()
	}

	if e.spanIDGenerator != nil {
		spanID = e.spanIDGenerator.NewID()
	} else {
		spanID = e.idService.NewID()
	}
	tc := &contract.TraceContext{
		TraceID:    traceID,
		ParentID:   "",
		SpanID:     spanID,
		CspanID:    "",
		Annotation: map[string]string{},
	}
	return tc
}

// ChildSpan instance a sub trace with new span id
func (e *ExtTrace) StartSpan(tc *contract.TraceContext) *contract.TraceContext {
	var childSpanID string
	if e.spanIDGenerator != nil {
		childSpanID = e.spanIDGenerator.NewID()
	} else {
		childSpanID = e.idService.NewID()
	}
	childSpan := &contract.TraceContext{
		TraceID:  tc.TraceID,
		ParentID: "",
		SpanID:   tc.SpanID,
		CspanID:  childSpanID,
		Annotation: map[string]string{
			contract.TraceKeyTime: time.Now().String(),
		},
	}
	return childSpan
}

// GetTrace By Http
func (e *ExtTrace) ExtractHTTP(req *http.Request) *contract.TraceContext {
	tc := &contract.TraceContext{}
	tc.TraceID = req.Header.Get(contract.TraceKeyTraceID)
	tc.ParentID = req.Header.Get(contract.TraceKeySpanID)
	tc.SpanID = req.Header.Get(contract.TraceKeyCspanID)
	tc.CspanID = ""

	if tc.TraceID == "" {
		tc.TraceID = e.idService.NewID()
	}

	if tc.SpanID == "" {
		tc.SpanID = e.idService.NewID()
	}

	return tc
}

// Set Trace to Http
func (e *ExtTrace) InjectHTTP(req *http.Request, tc *contract.TraceContext) *http.Request {
	req.Header.Add(contract.TraceKeyTraceID, tc.TraceID)
	req.Header.Add(contract.TraceKeySpanID, tc.SpanID)
	req.Header.Add(contract.TraceKeyCspanID, tc.CspanID)
	req.Header.Add(contract.TraceKeyParentID, tc.ParentID)
	return req
}

func (e *ExtTrace) ToMap(tc *contract.TraceContext) map[string]string {
	m := map[string]string{}
	if tc == nil {
		return m
	}
	m[contract.TraceKeyTraceID] = tc.TraceID
	m[contract.TraceKeySpanID] = tc.SpanID
	m[contract.TraceKeyCspanID] = tc.CspanID
	m[contract.TraceKeyParentID] = tc.ParentID

	if tc.Annotation != nil {
		for k, v := range tc.Annotation {
			m[k] = v
		}
	}
	return m
}

// func (e *ExtTrace) SetTraceIDService(service contract.IDService) {
// 	e.traceIDGenerator = service
// }

// func (e *ExtTrace) SetSpanIDService(service contract.IDService) {
// 	e.spanIDGenerator = service
// }
