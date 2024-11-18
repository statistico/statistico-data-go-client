package statisticofootballdata

import (
	"fmt"
)

type ErrorBadGateway struct {
	err error
}

func (e ErrorBadGateway) Error() string {
	return fmt.Sprintf("error connecting to the data service: %s", e.err.Error())
}

type ErrorExternalServer struct {
	err error
}

func (e ErrorExternalServer) Error() string {
	return fmt.Sprintf("internal server error returned from the data service: %s", e.err.Error())
}

type ErrorInvalidArgument struct {
	err error
}

func (e ErrorInvalidArgument) Error() string {
	return fmt.Sprintf("invalid argument provided: %s", e.err.Error())
}

type ErrorNotFound struct {
	ID  uint64
	err error
}

func (e ErrorNotFound) Error() string {
	return fmt.Sprintf("resource with ID '%d' does not exist. Error: %s", e.ID, e.err.Error())
}
