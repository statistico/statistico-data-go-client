package statisticodata_test

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/statistico/statistico-data-go-grpc-client"
	"github.com/statistico/statistico-proto/data/go"
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

		request := statisticoproto.TeamResultRequest{
			TeamId: 1,
			Limit:  &wrappers.UInt64Value{Value: 8},
		}

		res1 := newProtoResult(78102)
		res2 := newProtoResult(78103)

		ctx := context.Background()

		m.On("GetResultsForTeam", ctx, &request, []grpc.CallOption(nil)).Return(stream, nil)
		stream.On("Recv").Once().Return(res1, nil)
		stream.On("Recv").Once().Return(res2, nil)
		stream.On("Recv").Once().Return(&statisticoproto.Result{}, io.EOF)

		results, err := client.ByTeam(ctx, &request)

		if len(err) != 0 {
			t.Fatal("Expected nil, got errors on channel")
		}

		one :=<- results
		two :=<- results

		assert.Equal(t, res1, one)
		assert.Equal(t, res2, two)
		m.AssertExpectations(t)
	})

	t.Run("returns error if invalid argument error returned by result client", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoResultClient)
		client := statisticodata.NewResultClient(m)

		stream := new(MockResultStream)

		request := statisticoproto.TeamResultRequest{
			TeamId: 1,
			Limit:  &wrappers.UInt64Value{Value: 8},
		}

		ctx := context.Background()

		e := status.Error(codes.InvalidArgument, "incorrect format")

		m.On("GetResultsForTeam", ctx, &request, []grpc.CallOption(nil)).Return(stream, e)

		_, err := client.ByTeam(ctx, &request)

		e = <-err

		if e == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "invalid argument provided: rpc error: code = InvalidArgument desc = incorrect format", e.Error())
		m.AssertExpectations(t)
	})

	t.Run("logs error and returns internal server error", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoResultClient)
		client := statisticodata.NewResultClient(m)

		stream := new(MockResultStream)

		request := statisticoproto.TeamResultRequest{
			TeamId: 1,
			Limit:  &wrappers.UInt64Value{Value: 8},
		}

		ctx := context.Background()

		e := status.Error(codes.Internal, "internal error")

		m.On("GetResultsForTeam", ctx, &request, []grpc.CallOption(nil)).Return(stream, e)

		_, err := client.ByTeam(ctx, &request)

		e = <-err

		if e == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "internal server error returned from external service: rpc error: code = Internal desc = internal error", e.Error())
		m.AssertExpectations(t)
	})

	t.Run("logs error and returns bad gateway error", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoResultClient)
		client := statisticodata.NewResultClient(m)

		stream := new(MockResultStream)

		request := statisticoproto.TeamResultRequest{
			TeamId: 1,
			Limit:  &wrappers.UInt64Value{Value: 8},
		}

		ctx := context.Background()

		e := status.Error(codes.Aborted, "aborted")

		m.On("GetResultsForTeam", ctx, &request, []grpc.CallOption(nil)).Return(stream, e)

		_, err := client.ByTeam(ctx, &request)

		e = <-err

		if e == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "error connecting to external service: rpc error: code = Aborted desc = aborted", e.Error())
		m.AssertExpectations(t)
	})

	t.Run("logs error and returns internal server error if error parsing stream", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoResultClient)
		client := statisticodata.NewResultClient(m)

		stream := new(MockResultStream)

		request := statisticoproto.TeamResultRequest{
			TeamId: 1,
			Limit:  &wrappers.UInt64Value{Value: 8},
		}

		ctx := context.Background()

		e := errors.New("oh damn")

		m.On("GetResultsForTeam", ctx, &request, []grpc.CallOption(nil)).Return(stream, nil)
		stream.On("Recv").Twice().Return(newProtoResult(17801), nil)
		stream.On("Recv").Once().Return(&statisticoproto.Result{}, e)

		_, err := client.ByTeam(ctx, &request)

		e = <-err

		if e == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "internal server error returned from external service: oh damn", e.Error())
		m.AssertExpectations(t)
	})
}

func TestResultClient_ByID(t *testing.T) {
	t.Run("returns a result struct", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoResultClient)
		client := statisticodata.NewResultClient(m)

		req := mock.MatchedBy(func (r *statisticoproto.ResultRequest) bool {
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

		req := mock.MatchedBy(func (r *statisticoproto.ResultRequest) bool {
			assert.Equal(t, uint64(78102), r.FixtureId)
			return true
		})

		ctx := context.Background()

		e := status.Error(codes.NotFound, "not found")

		m.On("GetById", ctx, req, []grpc.CallOption(nil)).Return(&statisticoproto.Result{}, e)

		_, err := client.ByID(ctx, uint64(78102))

		if err == nil {
			t.Fatal("Expected error got nil")
		}

		a := assert.New(t)
		a.Equal("resource with is '78102' does not exist. Error: rpc error: code = NotFound desc = not found", err.Error())
		m.AssertExpectations(t)
	})

	t.Run("logs error and returns internal server error if returned by client", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoResultClient)
		client := statisticodata.NewResultClient(m)

		req := mock.MatchedBy(func (r *statisticoproto.ResultRequest) bool {
			assert.Equal(t, uint64(78102), r.FixtureId)
			return true
		})

		ctx := context.Background()

		e := status.Error(codes.Internal, "internal server error")

		m.On("GetById", ctx, req, []grpc.CallOption(nil)).Return(&statisticoproto.Result{}, e)

		_, err := client.ByID(ctx, uint64(78102))

		if err == nil {
			t.Fatal("Expected error got nil")
		}

		assert.Equal(t, "internal server error returned from external service: rpc error: code = Internal desc = internal server error", err.Error())
		m.AssertExpectations(t)
	})

	t.Run("logs error and returns bad gateway error", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoResultClient)
		client := statisticodata.NewResultClient(m)

		req := mock.MatchedBy(func (r *statisticoproto.ResultRequest) bool {
			assert.Equal(t, uint64(78102), r.FixtureId)
			return true
		})

		ctx := context.Background()

		e := status.Error(codes.Aborted, "internal server error")

		m.On("GetById", ctx, req, []grpc.CallOption(nil)).Return(&statisticoproto.Result{}, e)

		_, err := client.ByID(ctx, uint64(78102))

		if err == nil {
			t.Fatal("Expected error got nil")
		}

		assert.Equal(t, "error connecting to external service: rpc error: code = Aborted desc = internal server error", err.Error())
		m.AssertExpectations(t)
	})
}

func assertResult(t *testing.T, result *statisticoproto.Result) {
	a := assert.New(t)

	home := statisticoproto.Team{
		Id:             1,
		Name:           "West Ham United",
		ShortCode:      &wrappers.StringValue{Value: "WHU"},
		CountryId:      8,
		VenueId:        214,
		IsNationalTeam: &wrappers.BoolValue{Value: false},
		Founded:        &wrappers.UInt64Value{Value: 1895},
		Logo:           &wrappers.StringValue{Value: "logo"},
	}

	away := statisticoproto.Team{
		Id:             10,
		Name:           "Nottingham Forest",
		ShortCode:      &wrappers.StringValue{Value: "NOT"},
		CountryId:      8,
		VenueId:        300,
		IsNationalTeam: &wrappers.BoolValue{Value: true},
		Founded:        &wrappers.UInt64Value{Value: 1895},
		Logo:           &wrappers.StringValue{Value: "logo"},
	}

	season := statisticoproto.Season{
		Id:        16036,
		Name:      "2019/2020",
		IsCurrent: &wrappers.BoolValue{Value: true},
	}

	round := statisticoproto.Round{
		Id:        38,
		Name:      "38",
		SeasonId:  16036,
		StartDate: "2020-07-07T12:00:00+00:00",
		EndDate:   "2020-07-23T23:59:59+00:00",
	}

	venue := statisticoproto.Venue{
		Id:   214,
		Name: "London Stadium",
	}

	date := statisticoproto.Date{
		Utc: 1594132077,
		Rfc: "2020-07-07T15:00:00+00:00",
	}

	stats := statisticoproto.MatchStats{
		HomeScore: &wrappers.UInt32Value{Value: 5},
		AwayScore: &wrappers.UInt32Value{Value: 2},
	}

	a.Equal(uint64(78102), result.Id)
	a.Equal(home, *result.HomeTeam)
	a.Equal(away, *result.AwayTeam)
	a.Equal(season, *result.Season)
	a.Equal(round, *result.Round)
	a.Equal(venue, *result.Venue)
	a.Equal(stats, *result.Stats)
	a.Equal(date, *result.DateTime)
}

func newProtoResult(id uint64) *statisticoproto.Result {
	home := statisticoproto.Team{
		Id:             1,
		Name:           "West Ham United",
		ShortCode:      &wrappers.StringValue{Value: "WHU"},
		CountryId:      8,
		VenueId:        214,
		IsNationalTeam: &wrappers.BoolValue{Value: false},
		Founded:        &wrappers.UInt64Value{Value: 1895},
		Logo:           &wrappers.StringValue{Value: "logo"},
	}

	away := statisticoproto.Team{
		Id:             10,
		Name:           "Nottingham Forest",
		ShortCode:      &wrappers.StringValue{Value: "NOT"},
		CountryId:      8,
		VenueId:        300,
		IsNationalTeam: &wrappers.BoolValue{Value: true},
		Founded:        &wrappers.UInt64Value{Value: 1895},
		Logo:           &wrappers.StringValue{Value: "logo"},
	}

	season := statisticoproto.Season{
		Id:        16036,
		Name:      "2019/2020",
		IsCurrent: &wrappers.BoolValue{Value: true},
	}

	round := statisticoproto.Round{
		Id:        38,
		Name:      "38",
		SeasonId:  16036,
		StartDate: "2020-07-07T12:00:00+00:00",
		EndDate:   "2020-07-23T23:59:59+00:00",
	}

	venue := statisticoproto.Venue{
		Id:   214,
		Name: "London Stadium",
	}

	date := statisticoproto.Date{
		Utc: 1594132077,
		Rfc: "2020-07-07T15:00:00+00:00",
	}

	stats := statisticoproto.MatchStats{
		HomeScore: &wrappers.UInt32Value{Value: 5},
		AwayScore: &wrappers.UInt32Value{Value: 2},
	}

	return &statisticoproto.Result{
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

func (m *MockProtoResultClient) GetResultsForTeam(ctx context.Context, in *statisticoproto.TeamResultRequest, opts ...grpc.CallOption) (statisticoproto.ResultService_GetResultsForTeamClient, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(statisticoproto.ResultService_GetResultsForTeamClient), args.Error(1)
}

func (m *MockProtoResultClient) GetById(ctx context.Context, in *statisticoproto.ResultRequest, opts ...grpc.CallOption) (*statisticoproto.Result, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*statisticoproto.Result), args.Error(1)
}

func (m *MockProtoResultClient) GetResultsForSeason(ctx context.Context, in *statisticoproto.SeasonRequest, opts ...grpc.CallOption) (statisticoproto.ResultService_GetResultsForSeasonClient, error) {
	args := m.Called(ctx, in ,opts)
	return args.Get(0).(statisticoproto.ResultService_GetResultsForSeasonClient), args.Error(1)
}

func (m *MockProtoResultClient) GetHistoricalResultsForFixture(ctx context.Context, in *statisticoproto.HistoricalResultRequest, opts ...grpc.CallOption) (statisticoproto.ResultService_GetHistoricalResultsForFixtureClient, error) {
	args := m.Called(ctx, in ,opts)
	return args.Get(0).(statisticoproto.ResultService_GetResultsForSeasonClient), args.Error(1)
}

type MockResultStream struct {
	mock.Mock
	grpc.ClientStream
}

func (r *MockResultStream) Recv() (*statisticoproto.Result, error) {
	args := r.Called()
	return args.Get(0).(*statisticoproto.Result), args.Error(1)
}
