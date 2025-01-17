package statisticofootballdata_test

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

func TestCompetitionClient_ByCountryID(t *testing.T) {
	t.Run("calls competition client and returns a slice of competition struct", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoCompetitionClient)
		client := statisticofootballdata.NewCompetitionClient(m)

		stream := new(MockCompetitionStream)

		request := statistico.CompetitionRequest{
			CountryIds: []uint64{462},
			Sort:       nil,
			IsCup:      nil,
		}

		ctx := context.Background()

		m.On("ListCompetitions", ctx, &request, []grpc.CallOption(nil)).Return(stream, nil)
		stream.On("Recv").Twice().Return(newProtoCompetition(), nil)
		stream.On("Recv").Once().Return(&statistico.Competition{}, io.EOF)

		competitions, err := client.ByCountryID(ctx, 462)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Equal(t, 2, len(competitions))
		m.AssertExpectations(t)
		stream.AssertExpectations(t)
	})

	t.Run("logs error and returns internal server error if internal server error returned by client", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoCompetitionClient)
		client := statisticofootballdata.NewCompetitionClient(m)

		stream := new(MockCompetitionStream)

		request := statistico.CompetitionRequest{
			CountryIds: []uint64{462},
			Sort:       nil,
			IsCup:      nil,
		}

		ctx := context.Background()

		e := status.Error(codes.Internal, "internal error")

		m.On("ListCompetitions", ctx, &request, []grpc.CallOption(nil)).Return(stream, e)

		_, err := client.ByCountryID(ctx, 462)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "internal server error returned from the data service: rpc error: code = Internal desc = internal error", err.Error())
		m.AssertExpectations(t)
		stream.AssertNotCalled(t, "Recv")
	})

	t.Run("logs error and returns bad gateway error for non internal server error returned by client", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoCompetitionClient)
		client := statisticofootballdata.NewCompetitionClient(m)

		stream := new(MockCompetitionStream)

		request := statistico.CompetitionRequest{
			CountryIds: []uint64{462},
			Sort:       nil,
			IsCup:      nil,
		}

		ctx := context.Background()

		e := status.Error(codes.Unavailable, "service unavailable")

		m.On("ListCompetitions", ctx, &request, []grpc.CallOption(nil)).Return(stream, e)

		_, err := client.ByCountryID(ctx, 462)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "error connecting to the data service: rpc error: code = Unavailable desc = service unavailable", err.Error())
		m.AssertExpectations(t)
		stream.AssertNotCalled(t, "Recv")
	})

	t.Run("logs error and returns internal server error if error reading from stream", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoCompetitionClient)
		client := statisticofootballdata.NewCompetitionClient(m)

		stream := new(MockCompetitionStream)

		request := statistico.CompetitionRequest{
			CountryIds: []uint64{462},
			Sort:       nil,
			IsCup:      nil,
		}

		ctx := context.Background()

		e := errors.New("oh damn")

		m.On("ListCompetitions", ctx, &request, []grpc.CallOption(nil)).Return(stream, nil)
		stream.On("Recv").Twice().Return(newProtoCompetition(), nil)
		stream.On("Recv").Once().Return(&statistico.Competition{}, e)

		_, err := client.ByCountryID(ctx, 462)

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "internal server error returned from the data service: oh damn", err.Error())
		m.AssertExpectations(t)
		stream.AssertExpectations(t)
	})
}

func newProtoCompetition() *statistico.Competition {
	return &statistico.Competition{
		Id:        8,
		Name:      "Premier League",
		CountryId: 462,
	}
}

type MockProtoCompetitionClient struct {
	mock.Mock
}

func (c *MockProtoCompetitionClient) ListCompetitions(ctx context.Context, in *statistico.CompetitionRequest, opts ...grpc.CallOption) (statistico.CompetitionService_ListCompetitionsClient, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(statistico.CompetitionService_ListCompetitionsClient), args.Error(1)
}

type MockCompetitionStream struct {
	mock.Mock
	grpc.ClientStream
}

func (c *MockCompetitionStream) Recv() (*statistico.Competition, error) {
	args := c.Called()
	return args.Get(0).(*statistico.Competition), args.Error(1)
}
