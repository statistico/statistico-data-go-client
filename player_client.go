package statisticofootballdata

import (
	"context"
	"github.com/statistico/statistico-proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PlayerClient interface {
	ByID(ctx context.Context, id uint64) (*statistico.Player, error)
}

type playerClient struct {
	client statistico.PlayerServiceClient
}

func (t *playerClient) ByID(ctx context.Context, id uint64) (*statistico.Player, error) {
	req := statistico.PlayerRequest{PlayerId: id}

	player, err := t.client.GetPlayerByID(ctx, &req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return nil, ErrorNotFound{ID: id, err: err}
			default:
				return nil, ErrorBadGateway{err}
			}
		}

		return nil, err
	}

	return player, nil
}

func NewPlayerClient(p statistico.PlayerServiceClient) PlayerClient {
	return &playerClient{client: p}
}
