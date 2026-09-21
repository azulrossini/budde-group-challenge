package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/azulrossini/budde-group-challenge/backend/internal/api"
	"github.com/azulrossini/budde-group-challenge/backend/internal/apperr"
	"github.com/azulrossini/budde-group-challenge/backend/internal/httperr"
	"github.com/azulrossini/budde-group-challenge/backend/internal/middleware"
)

const (
	MsgMalformedJSON = "malformed JSON body"
	MsgInvalidBillID = "invalid billId"
)

// RejectUnknownFields enforces additionalProperties: false on PUT bodies
// (contract/openapi.yaml), which the generated strict handler's own decode
// step doesn't do on its own.
func RejectUnknownFields(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			next.ServeHTTP(w, r)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			httperr.Write(w, apperr.BadRequest(MsgMalformedJSON), middleware.RequestIDFromContext(r.Context()))
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))

		dec := json.NewDecoder(bytes.NewReader(body))
		dec.DisallowUnknownFields()
		var probe api.ReplaceSharesRequest
		if err := dec.Decode(&probe); err != nil {
			httperr.Write(w, apperr.BadRequest(MsgMalformedJSON), middleware.RequestIDFromContext(r.Context()))
			return
		}

		next.ServeHTTP(w, r)
	})
}
