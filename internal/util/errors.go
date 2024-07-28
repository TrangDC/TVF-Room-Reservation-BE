package util

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type ErrorFlag string
type ErrorMSg string

const (
	ErrorFlagUncategorized ErrorFlag = "UNCATEGORIZED"
	ErrorFlagCanNotDelete  ErrorFlag = "CAN_NOT_DELETE"
	ErrorFlagCanNotUpdate  ErrorFlag = "CAN_NOT_UPDATE"
	ErrorFlagCanNotDisable ErrorFlag = "CAN_NOT_DISABLE"
	ErrorFlagCanNotEnable  ErrorFlag = "CAN_NOT_ENABLE"
	ErrorFlagCanNotCreate  ErrorFlag = "CAN_NOT_CREATE"
	ErrorFlagUpdateFail    ErrorFlag = "UPDATE_FAIL"
	ErrorFlagCreateFail    ErrorFlag = "CREATE_FAIL"
	ErrorFlagDeleteFail    ErrorFlag = "DELETE_FAIL"
	ErrorFlagNotFound      ErrorFlag = "NOT_FOUND"
	ErrorFlagValidateFail  ErrorFlag = "VALIDATE_FAIL"
	ErrorFlagInternalError ErrorFlag = "INTERNAL_ERROR"
)

func WrapGQLError(ctx context.Context, message string, code int, errorFlag ErrorFlag) *gqlerror.Error {
	e := &gqlerror.Error{
		Message: strings.ToUpper(string(message[0])) + message[1:],
		Extensions: map[string]interface{}{
			"code":       code,
			"error_flag": errorFlag,
		},
	}
	if ctx != nil {
		e.Path = graphql.GetPath(ctx)
	}

	return e
}

func WrapGQLInternalError(ctx context.Context) *gqlerror.Error {
	return WrapGQLError(ctx, "Internal Server Error", http.StatusInternalServerError, ErrorFlagUncategorized)
}

func WrapGQLBadRequestError(ctx context.Context, format string, args ...interface{}) *gqlerror.Error {
	return WrapGQLError(ctx, fmt.Sprintf(format, args...), http.StatusBadRequest, ErrorFlagUncategorized)
}

func WrapGQLUnauthorizedError(ctx context.Context) *gqlerror.Error {
	return WrapGQLError(ctx, "Unauthorized Request", http.StatusUnauthorized, ErrorFlagUncategorized)
}

func WrapGQLNotFoundError(ctx context.Context) *gqlerror.Error {
	return WrapGQLError(ctx, "Not Found", http.StatusNotFound, ErrorFlagNotFound)
}
