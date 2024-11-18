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
	"testing"
)

func TestEventClient_FixtureEvents(t *testing.T) {
	t.Run("calls event client and returns fixture events response", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoEventClient)
		client := statisticofootballdata.NewEventClient(m)

		req := mock.MatchedBy(func(r *statistico.FixtureRequest) bool {
			assert.Equal(t, uint64(78102), r.FixtureId)
			return true
		})

		ctx := context.Background()

		m.On("FixtureEvents", ctx, req, []grpc.CallOption(nil)).Return(&statistico.FixtureEventsResponse{}, nil)

		_, err := client.FixtureEvents(ctx, uint64(78102))

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		m.AssertExpectations(t)
	})

	t.Run("returns not found error if returned by client", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoEventClient)
		client := statisticofootballdata.NewEventClient(m)

		req := mock.MatchedBy(func(r *statistico.FixtureRequest) bool {
			assert.Equal(t, uint64(78102), r.FixtureId)
			return true
		})

		ctx := context.Background()

		e := status.Error(codes.NotFound, "not found")

		m.On("FixtureEvents", ctx, req, []grpc.CallOption(nil)).Return(&statistico.FixtureEventsResponse{}, e)

		_, err := client.FixtureEvents(ctx, uint64(78102))

		if err == nil {
			t.Fatal("Expected error got nil")
		}

		assert.Equal(t, "resource with ID '78102' does not exist. Error: rpc error: code = NotFound desc = not found", err.Error())
		m.AssertExpectations(t)
	})

	t.Run("returns internal server error if returned by client", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoEventClient)
		client := statisticofootballdata.NewEventClient(m)

		req := mock.MatchedBy(func(r *statistico.FixtureRequest) bool {
			assert.Equal(t, uint64(78102), r.FixtureId)
			return true
		})

		ctx := context.Background()

		e := status.Error(codes.Internal, "internal server error")

		m.On("FixtureEvents", ctx, req, []grpc.CallOption(nil)).Return(&statistico.FixtureEventsResponse{}, e)

		_, err := client.FixtureEvents(ctx, uint64(78102))

		if err == nil {
			t.Fatal("Expected error got nil")
		}

		assert.Equal(t, "internal server error returned from the data service: rpc error: code = Internal desc = internal server error", err.Error())
		m.AssertExpectations(t)
	})

	t.Run("returns bad gateway error", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoEventClient)
		client := statisticofootballdata.NewEventClient(m)

		req := mock.MatchedBy(func(r *statistico.FixtureRequest) bool {
			assert.Equal(t, uint64(78102), r.FixtureId)
			return true
		})

		ctx := context.Background()

		e := status.Error(codes.Aborted, "internal server error")

		m.On("FixtureEvents", ctx, req, []grpc.CallOption(nil)).Return(&statistico.FixtureEventsResponse{}, e)

		_, err := client.FixtureEvents(ctx, uint64(78102))

		if err == nil {
			t.Fatal("Expected error got nil")
		}

		assert.Equal(t, "error connecting to the data service: rpc error: code = Aborted desc = internal server error", err.Error())
		m.AssertExpectations(t)
	})
}

type MockProtoEventClient struct {
	mock.Mock
}

func (m *MockProtoEventClient) FixtureEvents(ctx context.Context, in *statistico.FixtureRequest, opts ...grpc.CallOption) (*statistico.FixtureEventsResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*statistico.FixtureEventsResponse), args.Error(1)
}
