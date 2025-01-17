package statisticofootballdata

import (
	"context"
	"github.com/statistico/statistico-proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type CompetitionClient interface {
	ByCountryID(ctx context.Context, countryId uint64) ([]*statistico.Competition, error)
}

type competitionClient struct {
	competitionClient statistico.CompetitionServiceClient
}

func (c *competitionClient) ByCountryID(ctx context.Context, countryId uint64) ([]*statistico.Competition, error) {
	competitions := []*statistico.Competition{}

	req := statistico.CompetitionRequest{CountryIds: []uint64{countryId}}

	stream, err := c.competitionClient.ListCompetitions(ctx, &req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Internal:
				return competitions, ErrorExternalServer{err}
			default:
				return competitions, ErrorBadGateway{err}
			}
		}
	}

	for {
		competition, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return competitions, ErrorExternalServer{err}
		}

		competitions = append(competitions, competition)
	}

	return competitions, nil
}

func NewCompetitionClient(c statistico.CompetitionServiceClient) CompetitionClient {
	return &competitionClient{competitionClient: c}
}
