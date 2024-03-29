package errors

import (
	e "errors"
	"fmt"
)

// CodeType is a string that contains error's code description.
type CodeType string

// KindType is a string that contains error's kind description.
type KindType string

// CustomError is a structure that encodes useful information about a given error.
// It's supposed to flow within the application in detriment of the the default golang error,
// since its Kind and Code attributes are the keys to express its semantic and uniqueness, respectively.
// It should be generated once by the peace of code that found the error (because it's where we have more context about the error),
// and be by passed to the upper layers of the application.
type CustomError struct {
	kind    KindType
	code    CodeType
	rootErr error
	message string
}

const (
	// CodeUnknown is the default code returned when the application doesn't attach any code into the error.
	CodeUnknown CodeType = "UNKNOWN"
	// KindUnexpected is the default kind returned when the application doesn't attach any kind into the error.
	KindUnexpected KindType = "UNEXPECTED"
	// KindConflict are errors caused by requests with data that conflicts with the current state of the system.
	KindConflict KindType = "CONFLICT"
	// KindInternal are errors caused by some internal fail like failed IO calls or invalid memory states.
	KindInternal KindType = "INTERNAL"
	// KindInvalidInput are errors caused by some invalid values on the input.
	KindInvalidInput KindType = "INVALID_INPUT"
	// KindNotFound are errors caused by any required resources that not exists on the data repository.
	KindNotFound KindType = "NOT_FOUND"
	// KindUnauthenticated are errors caused by an unauthenticated call.
	KindUnauthenticated KindType = "UNAUTHENTICATED"
	// KindUnauthorized are errors caused by an unauthorized call.
	KindUnauthorized KindType = "UNAUTHORIZED"
	// KindResourceExhausted indicates some resource has been exhausted, perhaps a per-user quota, or perhaps the entire file system is out of space.
	KindResourceExhausted KindType = "RESOURCE_EXHAUSTED"
)

// New returns a new instance of CustomError with the given message.
func New(message string, args ...interface{}) CustomError {
	return CustomError{
		kind:    KindUnexpected,
		code:    CodeUnknown,
		message: fmt.Sprintf(message, args...),
	}
}

// NewMissingRequiredDependency creates a new error that indicates a missing required dependency.
// It should be producing at struct constructors.
func NewMissingRequiredDependency(name string) error {
	return New("Missing required dependency: %s", name).WithKind(KindInvalidInput).WithCode("MISSING_REQUIRED_DEPENDENCY")
}

// NewValidationError creates a Validation error.
func NewValidationError(desc string) error {
	return New(desc).WithKind(KindInvalidInput).WithCode("VALIDATION_ERROR")
}

// WithKind return a copy of the CustomError with the given KindType filled.
func (ce CustomError) WithKind(kind KindType) CustomError {
	ce.kind = kind
	return ce
}

// WithCode return a copy of the CustomError with the given CodeType filled.
func (ce CustomError) WithCode(code CodeType) CustomError {
	ce.code = code
	return ce
}

// WithRootError returns a copy of the CustomError with the RootError filled.
func (ce CustomError) WithRootError(err error) CustomError {
	ce.rootErr = err
	return ce
}

// Error returns CustomError message.
func (ce CustomError) Error() string {
	msg := ce.message
	if msg == "" && ce.rootErr != nil {
		msg = ce.rootErr.Error()
	}
	return msg
}

// RootError tries to convert the given error into a CustomError.
// If so, it recursively tries to find the root error (non CustomError) in a CustomError RootError chain and returns its message.
func RootError(err error) string {
	if err == nil {
		return ""
	}

	var customError CustomError
	if e.As(err, &customError) {
		if customError.rootErr == nil {
			return customError.Error()
		}

		return RootError(customError.rootErr)
	}

	return err.Error()
}

// Kind this method receives an error, then compares its interface type with the CustomError interface.
// If the interfaces types matches, returns its kind.
func Kind(err error) KindType {
	var customError CustomError
	if e.As(err, &customError) {
		return customError.kind
	}
	return KindUnexpected
}

// Kind this method receives an error, then compares its interface type with the CustomError interface.
// If the interfaces types matches, returns its Code.
func Code(err error) CodeType {
	var customError CustomError
	if e.As(err, &customError) {
		return customError.code
	}
	return CodeUnknown
}
