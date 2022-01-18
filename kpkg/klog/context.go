package klog

import (
	"context"
)

type ctxMarker struct{}

var ctxMarkerKey = &ctxMarker{}

// WithCommonLog ...
func WithCommonLog(ctx context.Context, v interface{}) context.Context {
	return context.WithValue(ctx, ctxMarkerKey, v)
}

//FromContext ...
func FromContext(ctx context.Context) (Common, bool) {
	if ctx == nil {
		return Common{}, false
	}
	v, ok := ctx.Value(ctxMarkerKey).(Common)
	return v, ok
}
