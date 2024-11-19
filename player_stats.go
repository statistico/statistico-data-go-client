package statisticofootballdata

import (
	"context"
	statistico "github.com/statistico/statistico-proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PlayerStatsClient interface {
	FixtureStats(ctx context.Context, req *statistico.FixtureRequest) (*statistico.PlayerStatsResponse, error)
}

type playerStatsClient struct {
	client statistico.PlayerStatsServiceClient
}

func (p *playerStatsClient) FixtureStats(ctx context.Context, req *statistico.FixtureRequest) (*statistico.PlayerStatsResponse, error) {
	res, err := p.client.GetPlayerStatsForFixture(ctx, req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				return nil, ErrorInvalidArgument{err}
			case codes.Internal:
				return nil, ErrorExternalServer{err}
			default:
				return nil, ErrorBadGateway{err}
			}
		}

		return nil, err
	}

	return res, nil
}

func NewPlayerStatsClient(p statistico.PlayerStatsServiceClient) PlayerStatsClient {
	return &playerStatsClient{client: p}
}
