package statisticodata

import (
	"context"
	"github.com/statistico/statistico-proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type ResultClient interface {
	ByID(ctx context.Context, fixtureID uint64) (*statistico.Result, error)
	ByTeam(ctx context.Context, req *statistico.TeamResultRequest) ([]*statistico.Result, error)
}

type resultClient struct {
	client statistico.ResultServiceClient
}

func (r resultClient) ByID(ctx context.Context, fixtureID uint64) (*statistico.Result, error) {
	request := statistico.ResultRequest{FixtureId: fixtureID}

	result, err := r.client.GetById(ctx, &request)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return nil, ErrorNotFound{fixtureID, err}
			case codes.Internal:
				return nil, ErrorInternalServerError{err}
			default:
				return nil, ErrorBadGateway{err}
			}
		}
	}

	return result, nil
}

func (r resultClient) ByTeam(ctx context.Context, req *statistico.TeamResultRequest) ([]*statistico.Result, error) {
	res := []*statistico.Result{}

	stream, err := r.client.GetResultsForTeam(ctx, req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				return res, ErrorInvalidArgument{err}
			case codes.Internal:
				return res, ErrorInternalServerError{err}
			default:
				return res, ErrorBadGateway{err}
			}
		}

		return res, err
	}

	for {
		result, err := stream.Recv()

		if err == io.EOF {
			return res, nil
		}

		if err != nil {
			return res, ErrorInternalServerError{err}
		}

		res = append(res, result)
	}
}

func NewResultClient(p statistico.ResultServiceClient) ResultClient {
	return resultClient{client: p}
}
