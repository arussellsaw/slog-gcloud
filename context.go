package sloggcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/monzo/typhon"
)

// CloudContextFilter adds data to the context for the Google Cloud Run environment
func CloudContextFilter(r typhon.Request, s typhon.Service) typhon.Response {
	ctx := WithTrace(r.Context, r)
	r.Context = ctx

	return s(r)
}

type traceKey string

func WithTrace(ctx context.Context, r typhon.Request) context.Context {
	var trace string

	traceHeader := r.Header.Get("X-Cloud-Trace-Context")

	traceParts := strings.Split(traceHeader, "/")
	if len(traceParts) > 0 && len(traceParts[0]) > 0 {
		trace = fmt.Sprintf("projects/russellsaw/traces/%s", traceParts[0])
	}

	return context.WithValue(ctx, traceKey("trace"), trace)
}

func Trace(ctx context.Context) string {
	v, ok := ctx.Value(traceKey("trace")).(string)
	if !ok {
		return "NOT_FOUND"
	}
	return v
}
