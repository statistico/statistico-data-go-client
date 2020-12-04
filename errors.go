package statisticodata

import "errors"

var ErrorBadGateway = errors.New("error response returned from external service")
var ErrorInternalServerError = errors.New("internal server error")
var ErrorNotFound = errors.New("the resource requested does not exist")
