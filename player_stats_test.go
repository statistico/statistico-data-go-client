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

func TestPlayerStatsClient_FixtureStats(t *testing.T) {
	t.Run("calls player stats client and returns a player stats response struct", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoPlayerStatsClient)
		client := statisticofootballdata.NewPlayerStatsClient(m)

		request := statistico.FixtureRequest{
			FixtureId: uint64(5),
		}

		ctx := context.Background()

		playerOne := newProtoPlayerStat(10)
		playerTwo := newProtoPlayerStat(20)

		res := statistico.PlayerStatsResponse{
			HomeTeam: []*statistico.PlayerStats{playerOne},
			AwayTeam: []*statistico.PlayerStats{playerTwo},
		}

		m.On("GetPlayerStatsForFixture", ctx, &request, []grpc.CallOption(nil)).Return(&res, nil)

		stats, err := client.FixtureStats(ctx, &request)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Equal(t, playerOne, stats.HomeTeam[0])
		assert.Equal(t, playerTwo, stats.AwayTeam[0])
		m.AssertExpectations(t)
	})

	t.Run("returns an error if error returned by client", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoPlayerStatsClient)
		client := statisticofootballdata.NewPlayerStatsClient(m)

		request := statistico.FixtureRequest{
			FixtureId: uint64(5),
		}

		ctx := context.Background()

		res := statistico.PlayerStatsResponse{}

		m.On("GetPlayerStatsForFixture", ctx, &request, []grpc.CallOption(nil)).
			Return(&res, status.Error(codes.InvalidArgument, "invalid argument"))

		_, err := client.FixtureStats(ctx, &request)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "invalid argument provided: rpc error: code = InvalidArgument desc = invalid argument", err.Error())
		m.AssertExpectations(t)
	})

	t.Run("returns an error if internal server error occurs", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoPlayerStatsClient)
		client := statisticofootballdata.NewPlayerStatsClient(m)

		request := statistico.FixtureRequest{
			FixtureId: uint64(5),
		}

		ctx := context.Background()

		res := statistico.PlayerStatsResponse{}

		m.On("GetPlayerStatsForFixture", ctx, &request, []grpc.CallOption(nil)).
			Return(&res, status.Error(codes.Internal, "internal error"))

		_, err := client.FixtureStats(ctx, &request)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "internal server error returned from the data service: rpc error: code = Internal desc = internal error", err.Error())
		m.AssertExpectations(t)
	})
}

func newProtoPlayerStat(playerID uint64) *statistico.PlayerStats {
	return &statistico.PlayerStats{
		PlayerId: playerID,
		Assists:  &wrapperspb.Int32Value{Value: 2},
	}
}

type MockProtoPlayerStatsClient struct {
	mock.Mock
}

func (m *MockProtoPlayerStatsClient) GetPlayerStatsForFixture(ctx context.Context, in *statistico.FixtureRequest, opts ...grpc.CallOption) (*statistico.PlayerStatsResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*statistico.PlayerStatsResponse), args.Error(1)
}

func (m *MockProtoPlayerStatsClient) GetLineUpForFixture(ctx context.Context, in *statistico.FixtureRequest, opts ...grpc.CallOption) (*statistico.LineupResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*statistico.LineupResponse), args.Error(1)
}

func (m *MockProtoPlayerStatsClient) GetTeamSeasonPlayerStats(ctx context.Context, in *statistico.TeamSeasonPlayStatsRequest, opts ...grpc.CallOption) (statistico.PlayerStatsService_GetTeamSeasonPlayerStatsClient, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(statistico.PlayerStatsService_GetTeamSeasonPlayerStatsClient), args.Error(1)
}
