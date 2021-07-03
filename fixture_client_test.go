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

func TestFixtureClient_Search(t *testing.T) {
	t.Run("calls fixture proto client and returns a slice of fixture struct", func(t *testing.T) {
		t.Helper()

		pc := new(MockFixtureProtoClient)
		client := statisticodata.NewFixtureClient(pc)

		request := statistico.FixtureSearchRequest{}

		stream := new(MockFixtureStream)
		ctx := context.Background()

		pc.On("Search", ctx, &request, []grpc.CallOption(nil)).Return(stream, nil)

		stream.On("Recv").Twice().Return(&statistico.Fixture{}, nil)
		stream.On("Recv").Once().Return(&statistico.Fixture{}, io.EOF)

		fixtures, err := client.Search(ctx, &request)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Equal(t, 2, len(fixtures))
		pc.AssertExpectations(t)
	})

	t.Run("returns error if invalid argument error returned by result client", func(t *testing.T) {
		t.Helper()

		pc := new(MockFixtureProtoClient)
		client := statisticodata.NewFixtureClient(pc)

		request := statistico.FixtureSearchRequest{}

		stream := new(MockFixtureStream)
		ctx := context.Background()

		e := status.Error(codes.InvalidArgument, "incorrect format")

		pc.On("Search", ctx, &request, []grpc.CallOption(nil)).Return(stream, e)

		_, err := client.Search(ctx, &request)

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "invalid argument provided: rpc error: code = InvalidArgument desc = incorrect format", err.Error())
		pc.AssertExpectations(t)
	})

	t.Run("returns internal server error", func(t *testing.T) {
		t.Helper()

		pc := new(MockFixtureProtoClient)
		client := statisticodata.NewFixtureClient(pc)

		request := statistico.FixtureSearchRequest{}

		stream := new(MockFixtureStream)
		ctx := context.Background()

		e := status.Error(codes.Internal, "incorrect format")

		pc.On("Search", ctx, &request, []grpc.CallOption(nil)).Return(stream, e)

		_, err := client.Search(ctx, &request)

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "internal server error returned from the data service: rpc error: code = Internal desc = incorrect format", err.Error())
		pc.AssertExpectations(t)
	})

	t.Run("returns bad gateway error", func(t *testing.T) {
		t.Helper()

		pc := new(MockFixtureProtoClient)
		client := statisticodata.NewFixtureClient(pc)

		request := statistico.FixtureSearchRequest{}

		stream := new(MockFixtureStream)
		ctx := context.Background()

		e := status.Error(codes.Aborted, "aborted")

		pc.On("Search", ctx, &request, []grpc.CallOption(nil)).Return(stream, e)

		_, err := client.Search(ctx, &request)

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "error connecting to the data service: rpc error: code = Aborted desc = aborted", err.Error())
		pc.AssertExpectations(t)
	})

	t.Run("returns internal server error if error parsing stream", func(t *testing.T) {
		t.Helper()

		pc := new(MockFixtureProtoClient)
		client := statisticodata.NewFixtureClient(pc)

		request := statistico.FixtureSearchRequest{}

		stream := new(MockFixtureStream)
		ctx := context.Background()

		pc.On("Search", ctx, &request, []grpc.CallOption(nil)).Return(stream, nil)

		e := errors.New("oh damn")

		stream.On("Recv").Twice().Return(&statistico.Fixture{}, nil)
		stream.On("Recv").Once().Return(&statistico.Fixture{}, e)

		_, err := client.Search(ctx, &request)

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "internal server error returned from the data service: oh damn", err.Error())
		pc.AssertExpectations(t)
	})
}

func TestFixtureClient_ByID(t *testing.T) {
	t.Run("returns a fixture struct", func(t *testing.T) {
		t.Helper()

		m := new(MockFixtureProtoClient)
		client := statisticodata.NewFixtureClient(m)

		req := mock.MatchedBy(func (r *statistico.FixtureRequest) bool {
			assert.Equal(t, uint64(78102), r.FixtureId)
			return true
		})

		ctx := context.Background()

		m.On("FixtureByID", ctx, req, []grpc.CallOption(nil)).Return(newProtoFixture(78102), nil)

		fixture, err := client.ByID(ctx, uint64(78102))

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		m.AssertExpectations(t)
		assert.Equal(t, int64(78102), fixture.Id)
	})

	t.Run("returns not found error if returned by client", func(t *testing.T) {
		t.Helper()

		m := new(MockFixtureProtoClient)
		client := statisticodata.NewFixtureClient(m)

		req := mock.MatchedBy(func (r *statistico.FixtureRequest) bool {
			assert.Equal(t, uint64(78102), r.FixtureId)
			return true
		})

		ctx := context.Background()

		e := status.Error(codes.NotFound, "not found")

		m.On("FixtureByID", ctx, req, []grpc.CallOption(nil)).Return(&statistico.Fixture{}, e)

		_, err := client.ByID(ctx, uint64(78102))

		if err == nil {
			t.Fatal("Expected error got nil")
		}

		assert.Equal(t, "resource with ID '78102' does not exist. Error: rpc error: code = NotFound desc = not found", err.Error())
		m.AssertExpectations(t)
	})

	t.Run("logs error and returns internal server error if returned by client", func(t *testing.T) {
		t.Helper()

		m := new(MockFixtureProtoClient)
		client := statisticodata.NewFixtureClient(m)

		req := mock.MatchedBy(func (r *statistico.FixtureRequest) bool {
			assert.Equal(t, uint64(78102), r.FixtureId)
			return true
		})

		ctx := context.Background()

		e := status.Error(codes.Internal, "internal server error")

		m.On("FixtureByID", ctx, req, []grpc.CallOption(nil)).Return(&statistico.Fixture{}, e)

		_, err := client.ByID(ctx, uint64(78102))

		if err == nil {
			t.Fatal("Expected error got nil")
		}

		assert.Equal(t, "internal server error returned from the data service: rpc error: code = Internal desc = internal server error", err.Error())
		m.AssertExpectations(t)
	})

	t.Run("logs error and returns bad gateway error", func(t *testing.T) {
		t.Helper()

		m := new(MockFixtureProtoClient)
		client := statisticodata.NewFixtureClient(m)

		req := mock.MatchedBy(func (r *statistico.FixtureRequest) bool {
			assert.Equal(t, uint64(78102), r.FixtureId)
			return true
		})

		ctx := context.Background()

		e := status.Error(codes.Aborted, "internal server error")

		m.On("FixtureByID", ctx, req, []grpc.CallOption(nil)).Return(&statistico.Fixture{}, e)

		_, err := client.ByID(ctx, uint64(78102))

		if err == nil {
			t.Fatal("Expected error got nil")
		}

		assert.Equal(t, "error connecting to the data service: rpc error: code = Aborted desc = internal server error", err.Error())
		m.AssertExpectations(t)
	})
}

func newProtoFixture(id int64) *statistico.Fixture {
	return &statistico.Fixture{Id: id}
}

type MockFixtureStream struct {
	mock.Mock
	grpc.ClientStream
}

func (f *MockFixtureStream) Recv() (*statistico.Fixture, error) {
	args := f.Called()
	return args.Get(0).(*statistico.Fixture), args.Error(1)
}

type MockFixtureProtoClient struct {
	mock.Mock
}

func (f *MockFixtureProtoClient) ListSeasonFixtures(ctx context.Context, in *statistico.SeasonFixtureRequest, opts ...grpc.CallOption) (statistico.FixtureService_ListSeasonFixturesClient, error) {
	args := f.Called(ctx, in, opts)
	return args.Get(0).(statistico.FixtureService_ListSeasonFixturesClient), args.Error(1)
}

func (f *MockFixtureProtoClient) FixtureByID(ctx context.Context, in *statistico.FixtureRequest, opts ...grpc.CallOption) (*statistico.Fixture, error) {
	args := f.Called(ctx, in, opts)
	return args.Get(0).(*statistico.Fixture), args.Error(1)
}

func (f *MockFixtureProtoClient) Search(ctx context.Context, in *statistico.FixtureSearchRequest, opts ...grpc.CallOption) (statistico.FixtureService_SearchClient, error) {
	args := f.Called(ctx, in, opts)
	return args.Get(0).(statistico.FixtureService_SearchClient), args.Error(1)
}
