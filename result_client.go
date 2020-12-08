package statisticodata

import (
	"context"
	"github.com/statistico/statistico-proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type ResultClient interface {
	ByID(ctx context.Context, fixtureID uint64) (*statisticoproto.Result, error)
	ByTeam(ctx context.Context, req *statisticoproto.TeamResultRequest) (<-chan *statisticoproto.Result, <-chan error)
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

func (r resultClient) ByTeam(ctx context.Context, req *statisticoproto.TeamResultRequest) (<-chan *statisticoproto.Result, <-chan error) {
	ch := make(chan *statisticoproto.Result, req.GetLimit().Value)
	errCh := make(chan error)

	go r.streamResults(ctx, req, ch, errCh)

	return ch, errCh
}

func (r resultClient) streamResults(ctx context.Context, req *statisticoproto.TeamResultRequest, ch chan<- *statisticoproto.Result, errChan chan<- error) {
	stream, err := r.client.GetResultsForTeam(ctx, req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				errChan <- ErrorInvalidArgument{err}
				break
			case codes.Internal:
				errChan <- ErrorInternalServerError{err}
				break
			default:
				errChan <- ErrorBadGateway{err}
				break
			}
		}

		closeChannels(ch, errChan)
		return
	}

	for {
		result, err := stream.Recv()

		if err == io.EOF {
			closeChannels(ch, errChan)
			return
		}

		if err != nil {
			errChan <- ErrorInternalServerError{err}
			closeChannels(ch, errChan)
			return
		}

		ch <- result
	}
}

func closeChannels(ch chan<- *statisticoproto.Result, errChan chan<- error) {
	close(ch)
	close(errChan)
}

func NewResultClient(p statisticoproto.ResultServiceClient) ResultClient {
	return resultClient{client: p}
}
