package statisticodata

import (
	"context"
	"github.com/statistico/statistico-proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type FixtureClient interface {
	Search(ctx context.Context, req *statistico.FixtureSearchRequest) ([]*statistico.Fixture, error)
	ByID(ctx context.Context, fixtureID uint64) (*statistico.Fixture, error)
}

type fixtureClient struct {
	client statistico.FixtureServiceClient
}

func (f *fixtureClient) ByID(ctx context.Context, fixtureID uint64) (*statistico.Fixture, error) {
	request := statistico.FixtureRequest{FixtureId: fixtureID}

	fixture, err := f.client.FixtureByID(ctx, &request)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return nil, ErrorNotFound{fixtureID, err}
			case codes.Internal:
				return nil, ErrorExternalServer{err}
			default:
				return nil, ErrorBadGateway{err}
			}
		}
	}

	return fixture, nil
}

func (f *fixtureClient) Search(ctx context.Context, req *statistico.FixtureSearchRequest) ([]*statistico.Fixture, error) {
	fixtures := []*statistico.Fixture{}

	stream, err := f.client.Search(ctx, req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				return fixtures, ErrorInvalidArgument{err}
			case codes.Internal:
				return fixtures, ErrorExternalServer{err}
			default:
				return fixtures, ErrorBadGateway{err}
			}
		}

		return fixtures, err
	}

	for {
		fixture, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fixtures, ErrorExternalServer{err: err}
		}

		fixtures = append(fixtures, fixture)
	}

	return fixtures, nil
}

func NewFixtureClient(p statistico.FixtureServiceClient) FixtureClient {
	return &fixtureClient{client: p}
}
