package statisticofootballdata

import (
	"context"
	"github.com/statistico/statistico-proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TeamStatClient interface {
	Stats(ctx context.Context, req *statistico.FixtureRequest) (*statistico.TeamStatsResponse, error)
}

type teamStatClient struct {
	client statistico.TeamStatsServiceClient
}

func (t *teamStatClient) Stats(ctx context.Context, req *statistico.FixtureRequest) (*statistico.TeamStatsResponse, error) {
	res, err := t.client.GetTeamStatsForFixture(ctx, req)

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

func NewTeamStatClient(p statistico.TeamStatsServiceClient) TeamStatClient {

	return &teamStatClient{client: p}
}
