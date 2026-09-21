package service

import (
	"fmt"
	"strings"

	"github.com/azulrossini/budde-group-challenge/backend/internal/apperr"
	"github.com/azulrossini/budde-group-challenge/backend/internal/models"
	"github.com/azulrossini/budde-group-challenge/backend/internal/money"
)

const (
	MaxNameLength = 100
	MaxShares     = 50
)

const (
	fieldShares                = "shares"
	fieldName                  = "name"
	fieldPercentageBasisPoints = "percentageBasisPoints"
)

const (
	msgAtLeastOneShare      = "at least one share is required"
	msgTooManyShares        = "at most %d shares are allowed"
	msgNameRequired         = "name is required"
	msgNameTooLong          = "name must be at most %d characters"
	msgDuplicateName        = "name must be unique"
	msgPercentageOutOfRange = "percentage must be greater than 0 and at most 100.00"
	msgSumMismatch          = "percentages must sum to 100.00 (got %s)"
	msgValidationFailed     = "validation failed"
)

// ValidateReplaceShares checks the rules from docs/SPEC.md §7 and returns a
// single *apperr.Error carrying every violation found, or nil if the input
// is valid.
func ValidateReplaceShares(inputs []models.ShareInput) *apperr.Error {
	var details []apperr.FieldError

	if len(inputs) == 0 {
		details = append(details, apperr.FieldError{Field: fieldShares, Message: msgAtLeastOneShare})
	}
	if len(inputs) > MaxShares {
		details = append(details, apperr.FieldError{Field: fieldShares, Message: fmt.Sprintf(msgTooManyShares, MaxShares)})
	}

	seenNames := make(map[string]bool, len(inputs))
	var totalBasisPoints int32

	for i, input := range inputs {
		name := strings.TrimSpace(input.Name)

		switch {
		case name == "":
			details = append(details, fieldError(i, fieldName, msgNameRequired))
		case len(name) > MaxNameLength:
			details = append(details, fieldError(i, fieldName, fmt.Sprintf(msgNameTooLong, MaxNameLength)))
		default:
			key := strings.ToLower(name)
			if seenNames[key] {
				details = append(details, fieldError(i, fieldName, msgDuplicateName))
			}
			seenNames[key] = true
		}

		if input.PercentageBasisPoints <= 0 || input.PercentageBasisPoints > money.FullPercentBasisPoints {
			details = append(details, fieldError(i, fieldPercentageBasisPoints, msgPercentageOutOfRange))
		}

		totalBasisPoints += input.PercentageBasisPoints
	}

	if len(inputs) > 0 && totalBasisPoints != money.FullPercentBasisPoints {
		details = append(details, apperr.FieldError{
			Field:   fieldShares,
			Message: fmt.Sprintf(msgSumMismatch, formatBasisPointsAsPercent(totalBasisPoints)),
		})
	}

	if len(details) == 0 {
		return nil
	}
	return apperr.Validation(msgValidationFailed, details)
}

func fieldError(index int, field, message string) apperr.FieldError {
	return apperr.FieldError{
		Field:   fmt.Sprintf("%s[%d].%s", fieldShares, index, field),
		Message: message,
	}
}

func formatBasisPointsAsPercent(bp int32) string {
	sign := ""
	if bp < 0 {
		sign = "-"
		bp = -bp
	}
	return fmt.Sprintf("%s%d.%02d", sign, bp/100, bp%100)
}
