package util

import (
	"context"
	"net/http"
)

type CommonParameters struct {
	Page    int
	PerPage int
	Keyword *string
}

func ValidateCommonParameters(ctx context.Context, pageParam, perPageParam *int, keywordParam *string) (*CommonParameters, error) {
	// Set default values for page and perPage if not provided
	page := 1
	perPage := 10
	if pageParam != nil {
		page = *pageParam
	}
	if perPageParam != nil {
		perPage = *perPageParam
	}

	// Validate page and perPage numbers
	if page < 1 {
		return nil, WrapGQLError(ctx, "page number must be greater than 0", http.StatusBadRequest, ErrorFlagValidateFail)
	}
	if perPage < 10 {
		return nil, WrapGQLError(ctx, "per page number must be greater than 20", http.StatusBadRequest, ErrorFlagValidateFail)
	}

	// Validate keyword parameters
	if keywordParam != nil && len(*keywordParam) > 50 {
		return nil, WrapGQLError(ctx, "keyword length exceeds maximum allowed limit", http.StatusBadRequest, ErrorFlagValidateFail)
	}

	return &CommonParameters{
		Page:    page,
		PerPage: perPage,
		Keyword: keywordParam,
	}, nil
}
