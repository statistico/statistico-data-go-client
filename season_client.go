package statisticodata

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/statistico/statistico-proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type SeasonClient interface {
	ByTeamID(ctx context.Context, teamId uint64, sort string) ([]*statisticoproto.Season, error)
	ByCompetitionID(ctx context.Context, competitionId uint64, sort string) ([]*statisticoproto.Season, error)
}

type seasonClient struct {
	client statisticoproto.SeasonServiceClient
}

func (s *seasonClient) ByTeamID(ctx context.Context, teamId uint64, sort string) ([]*statisticoproto.Season, error) {
	seasons := []*statisticoproto.Season{}

	req := statisticoproto.TeamSeasonsRequest{
		TeamId: teamId,
		Sort:   &wrappers.StringValue{Value: sort},
	}

	stream, err := s.client.GetSeasonsForTeam(ctx, &req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Internal:
				return seasons, ErrorInternalServerError{err}
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
			return seasons, ErrorInternalServerError{err}
		}

		seasons = append(seasons, season)
	}

	return seasons, nil
}

func (s *seasonClient) ByCompetitionID(ctx context.Context, competitionId uint64, sort string) ([]*statisticoproto.Season, error) {
	seasons := []*statisticoproto.Season{}

	req := statisticoproto.SeasonCompetitionRequest{CompetitionId: competitionId, Sort: &wrappers.StringValue{Value: sort}}

	stream, err := s.client.GetSeasonsForCompetition(ctx, &req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Internal:
				return seasons, ErrorInternalServerError{err}
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
			return seasons, ErrorInternalServerError{err}
		}

		seasons = append(seasons, season)
	}

	return seasons, nil
}

func NewSeasonClient(c statisticoproto.SeasonServiceClient) SeasonClient {
	return &seasonClient{client: c}
}
