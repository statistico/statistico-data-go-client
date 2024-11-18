package statisticofootballdata

import (
	"context"
	statistico "github.com/statistico/statistico-proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
)

type SeasonClient interface {
	ByTeamID(ctx context.Context, teamId uint64, sort string) ([]*statistico.Season, error)
	ByCompetitionID(ctx context.Context, competitionId uint64, sort string) ([]*statistico.Season, error)
}

type seasonClient struct {
	client statistico.SeasonServiceClient
}

func (s *seasonClient) ByTeamID(ctx context.Context, teamId uint64, sort string) ([]*statistico.Season, error) {
	seasons := []*statistico.Season{}

	req := statistico.TeamSeasonsRequest{
		TeamId: teamId,
		Sort:   &wrapperspb.StringValue{Value: sort},
	}

	response, err := s.client.GetSeasonsForTeam(ctx, &req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Internal:
				return seasons, ErrorExternalServer{err}
			default:
				return seasons, ErrorBadGateway{err}
			}
		}
	}

	return response.Seasons, nil
}

func (s *seasonClient) ByCompetitionID(ctx context.Context, competitionId uint64, sort string) ([]*statistico.Season, error) {
	seasons := []*statistico.Season{}

	req := statistico.SeasonCompetitionRequest{CompetitionId: competitionId, Sort: &wrapperspb.StringValue{Value: sort}}

	stream, err := s.client.GetSeasonsForCompetition(ctx, &req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Internal:
				return seasons, ErrorExternalServer{err}
			default:
				return seasons, ErrorBadGateway{err}
			}
		}
	}

	for {
		season, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return seasons, ErrorExternalServer{err}
		}

		seasons = append(seasons, season)
	}

	return seasons, nil
}

func NewSeasonClient(c statistico.SeasonServiceClient) SeasonClient {
	return &seasonClient{client: c}
}
