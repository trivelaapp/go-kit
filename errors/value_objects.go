package errors

var (
	// ErrResourceNotFound indicates that a desired resource was not found.
	ErrResourceNotFound error = New("resource not found").WithKind(KindNotFound).WithCode("RESOURCE_NOT_FOUND")

	// ErrNotImplemented indicates that a given feature is not implemented yet.
	ErrNotImplemented error = New("feature not implemented yet").WithCode("FEATURE_NOT_IMPLEMENTED")

	// ErrMock is a fake mocked that should be used in test scenarios.
	ErrMock error = New("mocked error").WithCode("MOCKED_ERROR")
)
