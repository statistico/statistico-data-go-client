package statisticodata

import (
	"context"
	"github.com/statistico/statistico-proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type CompetitionClient interface {
	ByCountryID(ctx context.Context, countryId uint64) ([]*statisticoproto.Competition, error)
}

type competitionClient struct {
	competitionClient statisticoproto.CompetitionServiceClient
}

func (c *competitionClient) ByCountryID(ctx context.Context, countryId uint64) ([]*statisticoproto.Competition, error) {
	competitions := []*statisticoproto.Competition{}

	req := statisticoproto.CompetitionRequest{CountryIds: []uint64{countryId}}

	stream, err := c.competitionClient.ListCompetitions(ctx, &req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Internal:
				return competitions, ErrorInternalServerError{err}
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
			return competitions, ErrorInternalServerError{err}
		}

		competitions = append(competitions, competition)
	}

	return competitions, nil
}

func NewCompetitionClient(c statisticoproto.CompetitionServiceClient) CompetitionClient {
	return &competitionClient{competitionClient: c}
}
