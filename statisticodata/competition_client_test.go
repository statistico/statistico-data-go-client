package statisticodata_test

import (
	"context"
	"errors"
	"github.com/statistico/statistico-data-go-grpc-client/statisticodata"
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
		client := statisticodata.NewCompetitionClient(m)

		stream := new(MockCompetitionStream)

		request := statisticoproto.CompetitionRequest{
			CountryIds: []uint64{462},
			Sort:       nil,
			IsCup:      nil,
		}

		ctx := context.Background()

		m.On("ListCompetitions", ctx, &request, []grpc.CallOption(nil)).Return(stream, nil)
		stream.On("Recv").Twice().Return(newProtoCompetition(), nil)
		stream.On("Recv").Once().Return(&statisticoproto.Competition{}, io.EOF)

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
		client := statisticodata.NewCompetitionClient(m)

		stream := new(MockCompetitionStream)

		request := statisticoproto.CompetitionRequest{
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

		assert.Equal(t, "internal server error returned from external service: rpc error: code = Internal desc = internal error", err.Error())
		m.AssertExpectations(t)
		stream.AssertNotCalled(t, "Recv")
	})

	t.Run("logs error and returns bad gateway error for non internal server error returned by client", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoCompetitionClient)
		client := statisticodata.NewCompetitionClient(m)

		stream := new(MockCompetitionStream)

		request := statisticoproto.CompetitionRequest{
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

		assert.Equal(t, "error connecting to external service: rpc error: code = Unavailable desc = service unavailable", err.Error())
		m.AssertExpectations(t)
		stream.AssertNotCalled(t, "Recv")
	})

	t.Run("logs error and returns internal server error if error reading from stream", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoCompetitionClient)
		client := statisticodata.NewCompetitionClient(m)

		stream := new(MockCompetitionStream)

		request := statisticoproto.CompetitionRequest{
			CountryIds: []uint64{462},
			Sort:       nil,
			IsCup:      nil,
		}

		ctx := context.Background()

		e := errors.New("oh damn")

		m.On("ListCompetitions", ctx, &request, []grpc.CallOption(nil)).Return(stream, nil)
		stream.On("Recv").Twice().Return(newProtoCompetition(), nil)
		stream.On("Recv").Once().Return(&statisticoproto.Competition{}, e)

		_, err := client.ByCountryID(ctx, 462)

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "internal server error returned from external service: oh damn", err.Error())
		m.AssertExpectations(t)
		stream.AssertExpectations(t)
	})
}

func newProtoCompetition() *statisticoproto.Competition {
	return &statisticoproto.Competition{
		Id:        8,
		Name:      "Premier League",
		IsCup:     false,
		CountryId: 462,
	}
}

type MockProtoCompetitionClient struct {
	mock.Mock
}

func (c *MockProtoCompetitionClient) ListCompetitions(ctx context.Context, in *statisticoproto.CompetitionRequest, opts ...grpc.CallOption) (statisticoproto.CompetitionService_ListCompetitionsClient, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(statisticoproto.CompetitionService_ListCompetitionsClient), args.Error(1)
}

type MockCompetitionStream struct {
	mock.Mock
	grpc.ClientStream
}

func (c *MockCompetitionStream) Recv() (*statisticoproto.Competition, error) {
	args := c.Called()
	return args.Get(0).(*statisticoproto.Competition), args.Error(1)
}
