package statisticodata

import (
	"context"
	"github.com/statistico/statistico-proto/data/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type ResultClient interface {
	ByID(ctx context.Context, fixtureID uint64) (*statisticoproto.Result, error)
	ByTeam(ctx context.Context, req *statisticoproto.TeamResultRequest) ([]*statisticoproto.Result, error)
}

type resultClient struct {
	client statisticoproto.ResultServiceClient
}

func (r resultClient) ByID(ctx context.Context, fixtureID uint64) (*statisticoproto.Result, error) {
	request := statisticoproto.ResultRequest{FixtureId: fixtureID}

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

func (r resultClient) ByTeam(ctx context.Context, req *statisticoproto.TeamResultRequest) ([]*statisticoproto.Result, error) {
	results := []*statisticoproto.Result{}

	stream, err := r.client.GetResultsForTeam(ctx, req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				return results, ErrorInvalidArgument{err}
			case codes.Internal:
				return results, ErrorInternalServerError{err}
			default:
				return results, ErrorBadGateway{err}
			}
		}
	}

	for {
		result, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return results, ErrorInternalServerError{err}
		}

		results = append(results, result)
	}

	return results, nil
}

func NewResultClient(p statisticoproto.ResultServiceClient) ResultClient {
	return resultClient{client: p}
}
