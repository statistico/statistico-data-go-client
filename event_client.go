package statisticofootballdata

import (
	"context"
	"github.com/statistico/statistico-proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventClient interface {
	FixtureEvents(ctx context.Context, fixtureID uint64) (*statistico.FixtureEventsResponse, error)
}

type eventClient struct {
	client statistico.EventServiceClient
}

func (e eventClient) FixtureEvents(ctx context.Context, fixtureID uint64) (*statistico.FixtureEventsResponse, error) {
	req := statistico.FixtureRequest{FixtureId: fixtureID}

	res, err := e.client.FixtureEvents(ctx, &req)

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

	return res, nil
}

func NewEventClient(c statistico.EventServiceClient) EventClient {
	return &eventClient{client: c}
}
