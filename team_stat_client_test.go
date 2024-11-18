package statisticofootballdata_test

import (
	"context"
	"github.com/statistico/statistico-football-data-go-grpc-client"
	"github.com/statistico/statistico-proto/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"testing"
)

func TestTeamStatClient_Stats(t *testing.T) {
	t.Run("calls team stat client and returns a team stats response struct", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoTeamStatsClient)
		client := statisticofootballdata.NewTeamStatClient(m)

		request := statistico.FixtureRequest{
			FixtureId: uint64(5),
		}

		ctx := context.Background()

		statOne := newProtoTeamStat(42)
		statTwo := newProtoTeamStat(43)

		res := statistico.TeamStatsResponse{
			HomeTeam: statOne,
			AwayTeam: statTwo,
		}

		m.On("GetTeamStatsForFixture", ctx, &request, []grpc.CallOption(nil)).Return(&res, nil)

		stats, err := client.Stats(ctx, &request)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Equal(t, statOne, stats.HomeTeam)
		assert.Equal(t, statTwo, stats.AwayTeam)
		m.AssertExpectations(t)
	})

	t.Run("returns an error if error returned by client", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoTeamStatsClient)
		client := statisticofootballdata.NewTeamStatClient(m)

		request := statistico.FixtureRequest{
			FixtureId: uint64(5),
		}

		ctx := context.Background()

		res := statistico.TeamStatsResponse{}

		m.On("GetTeamStatsForFixture", ctx, &request, []grpc.CallOption(nil)).
			Return(&res, status.Error(codes.InvalidArgument, "invalid argument"))

		_, err := client.Stats(ctx, &request)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "invalid argument provided: rpc error: code = InvalidArgument desc = invalid argument", err.Error())
		m.AssertExpectations(t)
	})
}

func newProtoTeamStat(fixtureID uint64) *statistico.TeamStats {
	return &statistico.TeamStats{FixtureId: fixtureID, Saves: &wrapperspb.Int32Value{Value: 2}}
}

type MockProtoTeamStatsClient struct {
	mock.Mock
}

func (m *MockProtoTeamStatsClient) GetTeamStatsForFixture(ctx context.Context, in *statistico.FixtureRequest, opts ...grpc.CallOption) (*statistico.TeamStatsResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*statistico.TeamStatsResponse), args.Error(1)
}
