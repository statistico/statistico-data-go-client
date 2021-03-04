package statisticodata

import (
	"context"
	"github.com/statistico/statistico-proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type TeamClient interface {
	ByID(ctx context.Context, teamID uint64) (*statistico.Team, error)
	BySeasonID(ctx context.Context, seasonId uint64) ([]*statistico.Team, error)
}

type teamClient struct {
	client statistico.TeamServiceClient
}

func (t *teamClient) ByID(ctx context.Context, teamID uint64) (*statistico.Team, error) {
	req := statistico.TeamRequest{TeamId: teamID}

	team, err := t.client.GetTeamByID(ctx, &req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return nil, ErrorNotFound{ID: teamID, err: err}
			default:
				return nil, ErrorBadGateway{err}
			}
		}

		return nil, err
	}

	return team, nil
}

func (t *teamClient) BySeasonID(ctx context.Context, seasonId uint64) ([]*statistico.Team, error) {
	teams := []*statistico.Team{}

	req := statistico.SeasonTeamsRequest{SeasonId: seasonId}

	stream, err := t.client.GetTeamsBySeasonId(ctx, &req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Internal:
				return teams, ErrorExternalServer{err}
			default:
				return teams, ErrorBadGateway{err}
			}
		}

		return nil, err
	}

	for {
		team, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return teams, ErrorExternalServer{err}
		}

		teams = append(teams, team)
	}

	return teams, nil
}

func NewTeamClient(p statistico.TeamServiceClient) TeamClient {
	return &teamClient{client: p}
}
