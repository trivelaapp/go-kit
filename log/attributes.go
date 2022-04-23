package log

// LogAttribute represents an information to be extracted from the context and included into the log.
type LogAttribute string

// LogAttributeSet is a set of LogAttributes.
type LogAttributeSet map[LogAttribute]bool

const (
	// LogAttributeRootError defines the name of the RooError attribute attached into logs.
	LogAttributeRootError LogAttribute = "root_error"

	// LogAttributeErrorKind defines the name of the ErrorKind attribute attached into logs.
	LogAttributeErrorKind LogAttribute = "err_kind"

	// LogAttributeErrorCode defines the name of the ErrorCode attribute attached into logs.
	LogAttributeErrorCode LogAttribute = "err_code"
)
