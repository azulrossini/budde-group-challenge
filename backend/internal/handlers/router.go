package handlers

import (
	"net/http"

	"github.com/azulrossini/budde-group-challenge/backend/internal/api"
	"github.com/azulrossini/budde-group-challenge/backend/internal/apperr"
	"github.com/azulrossini/budde-group-challenge/backend/internal/httperr"
	"github.com/azulrossini/budde-group-challenge/backend/internal/middleware"
)

// NewRouter wires the generated strict server to net/http. Our own
// middleware wraps the whole mux, rather than going through
// StdHTTPServerOptions.Middlewares: that hook only runs after the
// generated wrapper's own path-parameter binding succeeds, so a request
// with an invalid billId would otherwise skip RequestID/Logging/Recover
// entirely (no X-Request-ID header, no log line).
func NewRouter(h *Handlers) http.Handler {
	strictHandler := api.NewStrictHandlerWithOptions(h, nil, api.StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			httperr.Write(w, apperr.BadRequest(MsgMalformedJSON), middleware.RequestIDFromContext(r.Context()))
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			httperr.Write(w, apperr.Internal(err), middleware.RequestIDFromContext(r.Context()))
		},
	})

	mux := api.HandlerWithOptions(strictHandler, api.StdHTTPServerOptions{
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			httperr.Write(w, apperr.BadRequest(MsgInvalidBillID), middleware.RequestIDFromContext(r.Context()))
		},
	})

	handler := RejectUnknownFields(mux)
	handler = middleware.Auth(handler)
	handler = middleware.Logging(handler)
	handler = middleware.RequestID(handler)
	handler = middleware.Recover(handler)
	return handler
}
