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
	"testing"
)

func TestPlayerClient_ByID(t *testing.T) {
	t.Run("calls fixture client and return team struct", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoPlayerClient)
		client := statisticofootballdata.NewPlayerClient(m)

		request := statistico.PlayerRequest{PlayerId: 1}

		response := statistico.Player{
			Id:            1,
			CountryId:     8,
			NationalityId: 8,
			CommonName:    "Joe Sweeny",
			FirstName:     "Joe",
			LastName:      "Sweeny",
			Name:          "Joseph Sweeny",
			DisplayName:   "Joseph Sweeny",
			ImagePath:     "https://image.png",
			Height:        181,
			Weight:        78,
			DateOfBirth:   "1984-03-12",
			Gender:        "male",
		}

		ctx := context.Background()

		m.On("GetPlayerByID", ctx, &request, []grpc.CallOption(nil)).Return(&response, nil)

		player, err := client.ByID(ctx, uint64(1))

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Equal(t, uint64(1), player.GetId())
		assert.Equal(t, uint64(8), player.GetCountryId())
		assert.Equal(t, uint64(8), player.GetNationalityId())
		assert.Equal(t, "Joe Sweeny", player.GetCommonName())
		assert.Equal(t, "Joe", player.GetFirstName())
		assert.Equal(t, "Sweeny", player.GetLastName())
		assert.Equal(t, "Joseph Sweeny", player.GetName())
		assert.Equal(t, "Joseph Sweeny", player.GetDisplayName())
		assert.Equal(t, "https://image.png", player.GetImagePath())
		assert.Equal(t, "male", player.GetGender())
		assert.Equal(t, int32(181), player.GetHeight())
		assert.Equal(t, int32(78), player.GetWeight())
		assert.Equal(t, "1984-03-12", player.GetDateOfBirth())
		m.AssertExpectations(t)
	})

	t.Run("returns a not found if not found error is returned by grpc client", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoPlayerClient)
		client := statisticofootballdata.NewPlayerClient(m)

		request := statistico.PlayerRequest{PlayerId: 1}

		ctx := context.Background()

		e := status.Error(codes.NotFound, "not found")

		m.On("GetPlayerByID", ctx, &request, []grpc.CallOption(nil)).Return(&statistico.Player{}, e)

		_, err := client.ByID(ctx, uint64(1))

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "resource with ID '1' does not exist. Error: rpc error: code = NotFound desc = not found", err.Error())
		m.AssertExpectations(t)
	})

	t.Run("returns a bad gateway error", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoPlayerClient)
		client := statisticofootballdata.NewPlayerClient(m)

		request := statistico.PlayerRequest{PlayerId: 1}

		ctx := context.Background()

		e := status.Error(codes.Aborted, "aborted")

		m.On("GetPlayerByID", ctx, &request, []grpc.CallOption(nil)).Return(&statistico.Player{}, e)

		_, err := client.ByID(ctx, uint64(1))

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "error connecting to the data service: rpc error: code = Aborted desc = aborted", err.Error())
	})

	t.Run("returns an internal error", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoPlayerClient)
		client := statisticofootballdata.NewPlayerClient(m)

		request := statistico.PlayerRequest{PlayerId: 1}

		ctx := context.Background()

		e := errors.New("internal server error")

		m.On("GetPlayerByID", ctx, &request, []grpc.CallOption(nil)).Return(&statistico.Player{}, e)

		_, err := client.ByID(ctx, uint64(1))

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "internal server error", err.Error())
	})
}

type MockProtoPlayerClient struct {
	mock.Mock
}

func (t *MockProtoPlayerClient) GetPlayerByID(ctx context.Context, in *statistico.PlayerRequest, opts ...grpc.CallOption) (*statistico.Player, error) {
	args := t.Called(ctx, in, opts)
	return args.Get(0).(*statistico.Player), args.Error(1)
}
