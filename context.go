package sloggcloud

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/monzo/typhon"
)

// CloudContextFilter adds data to the context for the Google Cloud Run environment
func CloudContextFilter(r typhon.Request, s typhon.Service, projectID string) typhon.Response {
	ctx := WithTrace(r.Context, &r.Request, projectID)
	r.Context = ctx

	return s(r)
}

func CloudContextMiddleware(h http.Handler, projectID string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := WithTrace(r.Context(), r, projectID)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

type traceKey string

func WithTrace(ctx context.Context, r *http.Request, projectID string) context.Context {
	var trace string

	traceHeader := r.Header.Get("X-Cloud-Trace-Context")

	traceParts := strings.Split(traceHeader, "/")
	if len(traceParts) > 0 && len(traceParts[0]) > 0 {
		trace = fmt.Sprintf("projects/%s/traces/%s", projectID, traceParts[0])
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
