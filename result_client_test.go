package statisticodata_test

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/wrappers"
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

func TestResultClient_ByTeam(t *testing.T) {
	t.Run("calls result client and returns a slice of result struct", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoResultClient)
		client := statisticodata.NewResultClient(m)

		stream := new(MockResultStream)

		request := statistico.TeamResultRequest{
			TeamId: 1,
			Limit:  &wrappers.UInt64Value{Value: 8},
		}

		res1 := newProtoResult(78102)
		res2 := newProtoResult(78103)

		ctx := context.Background()

		m.On("GetResultsForTeam", ctx, &request, []grpc.CallOption(nil)).Return(stream, nil)
		stream.On("Recv").Once().Return(res1, nil)
		stream.On("Recv").Once().Return(res2, nil)
		stream.On("Recv").Once().Return(&statistico.Result{}, io.EOF)

		results, err := client.ByTeam(ctx, &request)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Equal(t, res1, results[0])
		assert.Equal(t, res2, results[1])
		m.AssertExpectations(t)
	})

	t.Run("returns error if invalid argument error returned by result client", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoResultClient)
		client := statisticodata.NewResultClient(m)

		stream := new(MockResultStream)

		request := statistico.TeamResultRequest{
			TeamId: 1,
			Limit:  &wrappers.UInt64Value{Value: 8},
		}

		ctx := context.Background()

		e := status.Error(codes.InvalidArgument, "incorrect format")

		m.On("GetResultsForTeam", ctx, &request, []grpc.CallOption(nil)).Return(stream, e)

		_, err := client.ByTeam(ctx, &request)

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "invalid argument provided: rpc error: code = InvalidArgument desc = incorrect format", err.Error())
		m.AssertExpectations(t)
	})

	t.Run("returns internal server error", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoResultClient)
		client := statisticodata.NewResultClient(m)

		stream := new(MockResultStream)

		request := statistico.TeamResultRequest{
			TeamId: 1,
			Limit:  &wrappers.UInt64Value{Value: 8},
		}

		ctx := context.Background()

		e := status.Error(codes.Internal, "internal error")

		m.On("GetResultsForTeam", ctx, &request, []grpc.CallOption(nil)).Return(stream, e)

		_, err := client.ByTeam(ctx, &request)

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "internal server error returned from the data service: rpc error: code = Internal desc = internal error", err.Error())
		m.AssertExpectations(t)
	})

	t.Run("returns bad gateway error", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoResultClient)
		client := statisticodata.NewResultClient(m)

		stream := new(MockResultStream)

		request := statistico.TeamResultRequest{
			TeamId: 1,
			Limit:  &wrappers.UInt64Value{Value: 8},
		}

		ctx := context.Background()

		e := status.Error(codes.Aborted, "aborted")

		m.On("GetResultsForTeam", ctx, &request, []grpc.CallOption(nil)).Return(stream, e)

		_, err := client.ByTeam(ctx, &request)

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "error connecting to the data service: rpc error: code = Aborted desc = aborted", err.Error())
		m.AssertExpectations(t)
	})

	t.Run("returns internal server error if error parsing stream", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoResultClient)
		client := statisticodata.NewResultClient(m)

		stream := new(MockResultStream)

		request := statistico.TeamResultRequest{
			TeamId: 1,
			Limit:  &wrappers.UInt64Value{Value: 8},
		}

		ctx := context.Background()

		e := errors.New("oh damn")

		m.On("GetResultsForTeam", ctx, &request, []grpc.CallOption(nil)).Return(stream, nil)
		stream.On("Recv").Twice().Return(newProtoResult(17801), nil)
		stream.On("Recv").Once().Return(&statistico.Result{}, e)

		_, err := client.ByTeam(ctx, &request)

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "internal server error returned from the data service: oh damn", err.Error())
		m.AssertExpectations(t)
	})
}

func TestResultClient_ByID(t *testing.T) {
	t.Run("returns a result struct", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoResultClient)
		client := statisticodata.NewResultClient(m)

		req := mock.MatchedBy(func(r *statistico.ResultRequest) bool {
			assert.Equal(t, uint64(78102), r.FixtureId)
			return true
		})

		ctx := context.Background()

		m.On("GetById", ctx, req, []grpc.CallOption(nil)).Return(newProtoResult(78102), nil)

		result, err := client.ByID(ctx, uint64(78102))

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		m.AssertExpectations(t)
		assertResult(t, result)
	})

	t.Run("returns not found error if returned by client", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoResultClient)
		client := statisticodata.NewResultClient(m)

		req := mock.MatchedBy(func(r *statistico.ResultRequest) bool {
			assert.Equal(t, uint64(78102), r.FixtureId)
			return true
		})

		ctx := context.Background()

		e := status.Error(codes.NotFound, "not found")

		m.On("GetById", ctx, req, []grpc.CallOption(nil)).Return(&statistico.Result{}, e)

		_, err := client.ByID(ctx, uint64(78102))

		if err == nil {
			t.Fatal("Expected error got nil")
		}

		assert.Equal(t, "resource with ID '78102' does not exist. Error: rpc error: code = NotFound desc = not found", err.Error())
		m.AssertExpectations(t)
	})

	t.Run("logs error and returns internal server error if returned by client", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoResultClient)
		client := statisticodata.NewResultClient(m)

		req := mock.MatchedBy(func(r *statistico.ResultRequest) bool {
			assert.Equal(t, uint64(78102), r.FixtureId)
			return true
		})

		ctx := context.Background()

		e := status.Error(codes.Internal, "internal server error")

		m.On("GetById", ctx, req, []grpc.CallOption(nil)).Return(&statistico.Result{}, e)

		_, err := client.ByID(ctx, uint64(78102))

		if err == nil {
			t.Fatal("Expected error got nil")
		}

		assert.Equal(t, "internal server error returned from the data service: rpc error: code = Internal desc = internal server error", err.Error())
		m.AssertExpectations(t)
	})

	t.Run("logs error and returns bad gateway error", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoResultClient)
		client := statisticodata.NewResultClient(m)

		req := mock.MatchedBy(func(r *statistico.ResultRequest) bool {
			assert.Equal(t, uint64(78102), r.FixtureId)
			return true
		})

		ctx := context.Background()

		e := status.Error(codes.Aborted, "internal server error")

		m.On("GetById", ctx, req, []grpc.CallOption(nil)).Return(&statistico.Result{}, e)

		_, err := client.ByID(ctx, uint64(78102))

		if err == nil {
			t.Fatal("Expected error got nil")
		}

		assert.Equal(t, "error connecting to the data service: rpc error: code = Aborted desc = internal server error", err.Error())
		m.AssertExpectations(t)
	})
}

func assertResult(t *testing.T, result *statistico.Result) {
	home := statistico.Team{
		Id:             1,
		Name:           "West Ham United",
		ShortCode:      &wrappers.StringValue{Value: "WHU"},
		CountryId:      8,
		VenueId:        214,
		IsNationalTeam: &wrappers.BoolValue{Value: false},
		Founded:        &wrappers.UInt64Value{Value: 1895},
		Logo:           &wrappers.StringValue{Value: "logo"},
	}

	away := statistico.Team{
		Id:             10,
		Name:           "Nottingham Forest",
		ShortCode:      &wrappers.StringValue{Value: "NOT"},
		CountryId:      8,
		VenueId:        300,
		IsNationalTeam: &wrappers.BoolValue{Value: true},
		Founded:        &wrappers.UInt64Value{Value: 1895},
		Logo:           &wrappers.StringValue{Value: "logo"},
	}

	season := statistico.Season{
		Id:        16036,
		Name:      "2019/2020",
		IsCurrent: &wrappers.BoolValue{Value: true},
	}

	round := statistico.Round{
		Id:        38,
		Name:      "38",
		SeasonId:  16036,
		StartDate: "2020-07-07T12:00:00+00:00",
		EndDate:   "2020-07-23T23:59:59+00:00",
	}

	venue := statistico.Venue{
		Id:   214,
		Name: "London Stadium",
	}

	date := statistico.Date{
		Utc: 1594132077,
		Rfc: "2020-07-07T15:00:00+00:00",
	}

	stats := statistico.MatchStats{
		HomeScore: &wrappers.UInt32Value{Value: 5},
		AwayScore: &wrappers.UInt32Value{Value: 2},
	}

	assert.Equal(t, uint64(78102), result.Id)
	assert.Equal(t, &home, result.GetHomeTeam())
	assert.Equal(t, &away, result.GetAwayTeam())
	assert.Equal(t, &season, result.GetSeason())
	assert.Equal(t, &round, result.GetRound())
	assert.Equal(t, &venue, result.GetVenue())
	assert.Equal(t, &stats, result.GetStats())
	assert.Equal(t, &date, result.GetDateTime())
}

func newProtoResult(id uint64) *statistico.Result {
	home := statistico.Team{
		Id:             1,
		Name:           "West Ham United",
		ShortCode:      &wrappers.StringValue{Value: "WHU"},
		CountryId:      8,
		VenueId:        214,
		IsNationalTeam: &wrappers.BoolValue{Value: false},
		Founded:        &wrappers.UInt64Value{Value: 1895},
		Logo:           &wrappers.StringValue{Value: "logo"},
	}

	away := statistico.Team{
		Id:             10,
		Name:           "Nottingham Forest",
		ShortCode:      &wrappers.StringValue{Value: "NOT"},
		CountryId:      8,
		VenueId:        300,
		IsNationalTeam: &wrappers.BoolValue{Value: true},
		Founded:        &wrappers.UInt64Value{Value: 1895},
		Logo:           &wrappers.StringValue{Value: "logo"},
	}

	season := statistico.Season{
		Id:        16036,
		Name:      "2019/2020",
		IsCurrent: &wrappers.BoolValue{Value: true},
	}

	round := statistico.Round{
		Id:        38,
		Name:      "38",
		SeasonId:  16036,
		StartDate: "2020-07-07T12:00:00+00:00",
		EndDate:   "2020-07-23T23:59:59+00:00",
	}

	venue := statistico.Venue{
		Id:   214,
		Name: "London Stadium",
	}

	date := statistico.Date{
		Utc: 1594132077,
		Rfc: "2020-07-07T15:00:00+00:00",
	}

	stats := statistico.MatchStats{
		HomeScore: &wrappers.UInt32Value{Value: 5},
		AwayScore: &wrappers.UInt32Value{Value: 2},
	}

	return &statistico.Result{
		Id:       id,
		HomeTeam: &home,
		AwayTeam: &away,
		Season:   &season,
		Round:    &round,
		Venue:    &venue,
		DateTime: &date,
		Stats:    &stats,
	}
}

type MockProtoResultClient struct {
	mock.Mock
}

func (m *MockProtoResultClient) GetResultsForTeam(ctx context.Context, in *statistico.TeamResultRequest, opts ...grpc.CallOption) (statistico.ResultService_GetResultsForTeamClient, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(statistico.ResultService_GetResultsForTeamClient), args.Error(1)
}

func (m *MockProtoResultClient) GetById(ctx context.Context, in *statistico.ResultRequest, opts ...grpc.CallOption) (*statistico.Result, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*statistico.Result), args.Error(1)
}

func (m *MockProtoResultClient) GetResultsForSeason(ctx context.Context, in *statistico.SeasonRequest, opts ...grpc.CallOption) (statistico.ResultService_GetResultsForSeasonClient, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(statistico.ResultService_GetResultsForSeasonClient), args.Error(1)
}

func (m *MockProtoResultClient) GetHistoricalResultsForFixture(ctx context.Context, in *statistico.HistoricalResultRequest, opts ...grpc.CallOption) (statistico.ResultService_GetHistoricalResultsForFixtureClient, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(statistico.ResultService_GetResultsForSeasonClient), args.Error(1)
}

type MockResultStream struct {
	mock.Mock
	grpc.ClientStream
}

func (r *MockResultStream) Recv() (*statistico.Result, error) {
	args := r.Called()
	return args.Get(0).(*statistico.Result), args.Error(1)
}
