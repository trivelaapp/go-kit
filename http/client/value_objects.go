package client

// HTTPHeaders is a map containing the relation key=value of the headers used on the http rest request.
type HTTPHeaders map[string]string

// HTTPQueryParams is a map containing the relation key=value of the query params used on the http rest request
type HTTPQueryParams map[string]string

// HTTPRequest are the params used to build a new http rest request
type HTTPRequest struct {
	URL         string
	Body        []byte
	Headers     HTTPHeaders
	QueryParams HTTPQueryParams
}

// HTTPResult are the params returned from the client HTTP request
type HTTPResult struct {
	StatusCode int
	Response   []byte
}
