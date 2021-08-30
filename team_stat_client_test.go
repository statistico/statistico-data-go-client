package statisticodata_test

import (
	"context"
	"errors"
	"github.com/statistico/statistico-football-data-go-grpc-client"
	"github.com/statistico/statistico-proto/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"testing"
)

func TestTeamStatClient_Stats(t *testing.T) {
	t.Run("calls team stat client and returns a channel of team stat struct", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoTeamStatsClient)
		client := statisticodata.NewTeamStatClient(m)

		stream := new(MockTeamStatStream)

		request := statistico.TeamStatRequest{
			Stat:      "shots_total",
			TeamId:    5,
			SeasonIds: []uint64{16036},
		}

		ctx := context.Background()

		statOne := newProtoTeamStat(42)
		statTwo := newProtoTeamStat(43)

		m.On("GetStatForTeam", ctx, &request, []grpc.CallOption(nil)).Return(stream, nil)
		stream.On("Recv").Once().Return(statOne, nil)
		stream.On("Recv").Once().Return(statTwo, nil)
		stream.On("Recv").Once().Return(&statistico.TeamStat{}, io.EOF)

		stats, err := client.Stats(ctx, &request)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Equal(t, statOne, stats[0])
		assert.Equal(t, statTwo, stats[1])
		m.AssertExpectations(t)
	})

	t.Run("returns error in error channel if invalid argument error returned by team stat client", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoTeamStatsClient)
		client := statisticodata.NewTeamStatClient(m)

		stream := new(MockTeamStatStream)

		request := statistico.TeamStatRequest{
			Stat:      "shots_total",
			TeamId:    5,
			SeasonIds: []uint64{16036},
		}

		ctx := context.Background()

		e := status.Error(codes.InvalidArgument, "incorrect format")

		m.On("GetStatForTeam", ctx, &request, []grpc.CallOption(nil)).Return(stream, e)

		_, err := client.Stats(ctx, &request)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "invalid argument provided: rpc error: code = InvalidArgument desc = incorrect format", err.Error())
		m.AssertExpectations(t)
	})

	t.Run("logs error and returns internal server error in error channel", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoTeamStatsClient)
		client := statisticodata.NewTeamStatClient(m)

		stream := new(MockTeamStatStream)

		request := statistico.TeamStatRequest{
			Stat:      "shots_total",
			TeamId:    5,
			SeasonIds: []uint64{16036},
		}

		ctx := context.Background()

		e := status.Error(codes.Internal, "internal error")

		m.On("GetStatForTeam", ctx, &request, []grpc.CallOption(nil)).Return(stream, e)

		_, err := client.Stats(ctx, &request)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "internal server error returned from the data service: rpc error: code = Internal desc = internal error", err.Error())
		m.AssertExpectations(t)
	})

	t.Run("logs error and returns bad gateway error in error channel", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoTeamStatsClient)
		client := statisticodata.NewTeamStatClient(m)

		stream := new(MockTeamStatStream)

		request := statistico.TeamStatRequest{
			Stat:      "shots_total",
			TeamId:    5,
			SeasonIds: []uint64{16036},
		}

		ctx := context.Background()

		e := status.Error(codes.Aborted, "aborted")

		m.On("GetStatForTeam", ctx, &request, []grpc.CallOption(nil)).Return(stream, e)

		_, err := client.Stats(ctx, &request)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "error connecting to the data service: rpc error: code = Aborted desc = aborted", err.Error())
		m.AssertExpectations(t)
	})

	t.Run("logs error and returns internal server error in error channel if error parsing stream", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoTeamStatsClient)
		client := statisticodata.NewTeamStatClient(m)

		stream := new(MockTeamStatStream)

		request := statistico.TeamStatRequest{
			Stat:      "shots_total",
			TeamId:    5,
			SeasonIds: []uint64{16036},
		}

		ctx := context.Background()

		e := errors.New("oh damn")

		m.On("GetStatForTeam", ctx, &request, []grpc.CallOption(nil)).Return(stream, nil)
		stream.On("Recv").Once().Return(&statistico.TeamStat{}, e)

		_, err := client.Stats(ctx, &request)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "internal server error returned from the data service: oh damn", err.Error())
		m.AssertExpectations(t)
	})
}

func newProtoTeamStat(fixtureID uint64) *statistico.TeamStat {
	return &statistico.TeamStat{FixtureId: fixtureID, Stat: "shots_total"}
}

type MockProtoTeamStatsClient struct {
	mock.Mock
}

func (m *MockProtoTeamStatsClient) GetTeamStatsForFixture(ctx context.Context, in *statistico.FixtureRequest, opts ...grpc.CallOption) (*statistico.TeamStatsResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*statistico.TeamStatsResponse), args.Error(1)
}

func (m *MockProtoTeamStatsClient) GetStatForTeam(ctx context.Context, in *statistico.TeamStatRequest, opts ...grpc.CallOption) (statistico.TeamStatsService_GetStatForTeamClient, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(statistico.TeamStatsService_GetStatForTeamClient), args.Error(1)
}

type MockTeamStatStream struct {
	mock.Mock
	grpc.ClientStream
}

func (m *MockTeamStatStream) Recv() (*statistico.TeamStat, error) {
	args := m.Called()
	return args.Get(0).(*statistico.TeamStat), args.Error(1)
}
