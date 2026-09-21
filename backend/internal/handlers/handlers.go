package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/azulrossini/budde-group-challenge/backend/internal/api"
	"github.com/azulrossini/budde-group-challenge/backend/internal/apperr"
	"github.com/azulrossini/budde-group-challenge/backend/internal/httperr"
	"github.com/azulrossini/budde-group-challenge/backend/internal/middleware"
	"github.com/azulrossini/budde-group-challenge/backend/internal/models"
	"github.com/azulrossini/budde-group-challenge/backend/internal/service"
)

type Handlers struct {
	service *service.Service
}

func New(svc *service.Service) *Handlers {
	return &Handlers{service: svc}
}

func (h *Handlers) GetBillShares(ctx context.Context, request api.GetBillSharesRequestObject) (api.GetBillSharesResponseObject, error) {
	result, err := h.service.GetBillShares(ctx, request.BillId)
	if err != nil {
		return getBillSharesError(err, middleware.RequestIDFromContext(ctx)), nil
	}
	return api.GetBillShares200JSONResponse(toAPIResponse(result)), nil
}

func (h *Handlers) ReplaceBillShares(ctx context.Context, request api.ReplaceBillSharesRequestObject) (api.ReplaceBillSharesResponseObject, error) {
	result, err := h.service.ReplaceShares(ctx, request.BillId, toShareInputs(request.Body.Shares))
	if err != nil {
		return replaceBillSharesError(err, middleware.RequestIDFromContext(ctx)), nil
	}
	return api.ReplaceBillShares200JSONResponse(toAPIResponse(result)), nil
}

func toShareInputs(inputs []api.ShareInput) []models.ShareInput {
	result := make([]models.ShareInput, len(inputs))
	for i, input := range inputs {
		result[i] = models.ShareInput{Name: input.Name, PercentageBasisPoints: input.PercentageBasisPoints}
	}
	return result
}

func toAPIResponse(result service.BillSharesResult) api.BillSharesResponse {
	shares := make([]api.Share, len(result.Shares))
	for i, share := range result.Shares {
		shares[i] = api.Share{
			Id:                    share.ID,
			Name:                  share.Name,
			PercentageBasisPoints: share.PercentageBasisPoints,
			AmountCents:           share.AmountCents,
		}
	}

	return api.BillSharesResponse{
		Bill: api.Bill{
			Id:          result.Bill.ID,
			Description: result.Bill.Description,
			TotalCents:  result.Bill.TotalCents,
		},
		Shares:           shares,
		TotalBasisPoints: result.TotalBasisPoints,
	}
}

func asAppError(err error) *apperr.Error {
	var appErr *apperr.Error
	if errors.As(err, &appErr) {
		return appErr
	}
	return apperr.Internal(err)
}

func getBillSharesError(err error, requestID string) api.GetBillSharesResponseObject {
	status, body := httperr.Body(asAppError(err), requestID)
	switch status {
	case http.StatusBadRequest:
		return api.GetBillShares400JSONResponse{BadRequestJSONResponse: api.BadRequestJSONResponse(body)}
	case http.StatusNotFound:
		return api.GetBillShares404JSONResponse{NotFoundJSONResponse: api.NotFoundJSONResponse(body)}
	case http.StatusUnprocessableEntity:
		return api.GetBillShares422JSONResponse{ValidationErrorJSONResponse: api.ValidationErrorJSONResponse(body)}
	default:
		return api.GetBillShares500JSONResponse{InternalErrorJSONResponse: api.InternalErrorJSONResponse(body)}
	}
}

func replaceBillSharesError(err error, requestID string) api.ReplaceBillSharesResponseObject {
	status, body := httperr.Body(asAppError(err), requestID)
	switch status {
	case http.StatusBadRequest:
		return api.ReplaceBillShares400JSONResponse{BadRequestJSONResponse: api.BadRequestJSONResponse(body)}
	case http.StatusNotFound:
		return api.ReplaceBillShares404JSONResponse{NotFoundJSONResponse: api.NotFoundJSONResponse(body)}
	case http.StatusUnprocessableEntity:
		return api.ReplaceBillShares422JSONResponse{ValidationErrorJSONResponse: api.ValidationErrorJSONResponse(body)}
	default:
		return api.ReplaceBillShares500JSONResponse{InternalErrorJSONResponse: api.InternalErrorJSONResponse(body)}
	}
}
