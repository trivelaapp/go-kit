package error_handler

import (
	"google.golang.org/grpc/codes"

	"github.com/trivelaapp/go-kit/errors"
)

func kindToGRPCStatusCode(kind errors.KindType) codes.Code {
	switch kind {
	case errors.KindInvalidInput:
		return codes.InvalidArgument
	case errors.KindUnauthenticated:
		return codes.Unauthenticated
	case errors.KindUnauthorized:
		return codes.PermissionDenied
	case errors.KindNotFound:
		return codes.NotFound
	case errors.KindConflict:
		return codes.FailedPrecondition
	case errors.KindUnexpected:
		return codes.Unknown
	case errors.KindInternal:
		return codes.Internal
	default:
		return codes.Unknown
	}
}
