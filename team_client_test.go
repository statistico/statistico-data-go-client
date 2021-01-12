package statisticodata_test

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/statistico/statistico-data-go-grpc-client"
	"github.com/statistico/statistico-proto/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"testing"
)

func TestTeamClient_ByID(t *testing.T) {
	t.Run("calls fixture client and return team struct", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoTeamClient)
		client := statisticodata.NewTeamClient(m)

		request := statisticoproto.TeamRequest{TeamId: 1}

		response := statisticoproto.Team{
			Id:             1,
			Name:           "West Ham United",
			ShortCode:      &wrappers.StringValue{Value: "WHU"},
			CountryId:      8,
			VenueId:        214,
			IsNationalTeam: &wrappers.BoolValue{Value: false},
			Founded:        &wrappers.UInt64Value{Value: 1895},
			Logo:           &wrappers.StringValue{Value: "logo"},
		}

		ctx := context.Background()

		m.On("GetTeamByID", ctx, &request, []grpc.CallOption(nil)).Return(&response, nil)

		team, err := client.ByID(ctx, uint64(1))

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(uint64(1), team.GetId())
		a.Equal("West Ham United", team.GetName())
		a.Equal("WHU", team.GetShortCode().GetValue())
		a.Equal(uint64(8), team.GetCountryId())
		a.Equal(uint64(214), team.GetVenueId())
		a.Equal(false, team.GetIsNationalTeam().GetValue())
		a.Equal(uint64(1895), team.GetFounded().GetValue())
		a.Equal("logo", team.GetLogo().GetValue())
		m.AssertExpectations(t)
	})

	t.Run("parses nullable fields from team returned in response", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoTeamClient)
		client := statisticodata.NewTeamClient(m)

		request := statisticoproto.TeamRequest{TeamId: 1}

		response := statisticoproto.Team{
			Id:        1,
			Name:      "West Ham United",
			CountryId: 8,
			VenueId:   214,
		}

		ctx := context.Background()

		m.On("GetTeamByID", ctx, &request, []grpc.CallOption(nil)).Return(&response, nil)

		team, err := client.ByID(ctx, uint64(1))

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(uint64(1), team.Id)
		a.Equal("West Ham United", team.Name)
		a.Nil(team.ShortCode)
		a.Equal(uint64(8), team.CountryId)
		a.Equal(uint64(214), team.VenueId)
		a.Equal(false, team.IsNationalTeam.GetValue())
		a.Nil(team.Founded)
		a.Nil(team.Logo)
		m.AssertExpectations(t)
	})

	t.Run("returns a not found if not found error is returned by grpc client", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoTeamClient)
		client := statisticodata.NewTeamClient(m)

		request := statisticoproto.TeamRequest{TeamId: 1}

		ctx := context.Background()

		e := status.Error(codes.NotFound, "not found")

		m.On("GetTeamByID", ctx, &request, []grpc.CallOption(nil)).Return(&statisticoproto.Team{}, e)

		_, err := client.ByID(ctx, uint64(1))

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "resource with ID '1' does not exist. Error: rpc error: code = NotFound desc = not found", err.Error())
		m.AssertExpectations(t)
	})

	t.Run("returns a bad gateway error", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoTeamClient)
		client := statisticodata.NewTeamClient(m)

		request := statisticoproto.TeamRequest{TeamId: 1}

		ctx := context.Background()

		e := status.Error(codes.Aborted, "aborted")

		m.On("GetTeamByID", ctx, &request, []grpc.CallOption(nil)).Return(&statisticoproto.Team{}, e)

		_, err := client.ByID(ctx, uint64(1))

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "error connecting to external service: rpc error: code = Aborted desc = aborted", err.Error())
	})

	t.Run("returns an internal error", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoTeamClient)
		client := statisticodata.NewTeamClient(m)

		request := statisticoproto.TeamRequest{TeamId: 1}

		ctx := context.Background()

		e := errors.New("internal server error")

		m.On("GetTeamByID", ctx, &request, []grpc.CallOption(nil)).Return(&statisticoproto.Team{}, e)

		_, err := client.ByID(ctx, uint64(1))

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "internal server error returned from external service: internal server error", err.Error())
	})
}

func TestTeamClient_BySeasonID(t *testing.T) {
	t.Run("calls team client and returns a slice of team struct", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoTeamClient)
		client := statisticodata.NewTeamClient(m)

		stream := new(MockTeamStream)

		team := statisticoproto.Team{
			Id:        1,
			Name:      "West Ham United",
			CountryId: 8,
			VenueId:   214,
		}

		ctx := context.Background()

		request := statisticoproto.SeasonTeamsRequest{SeasonId: 16036}

		m.On("GetTeamsBySeasonId", ctx, &request, []grpc.CallOption(nil)).Return(stream, nil)
		stream.On("Recv").Twice().Return(&team, nil)
		stream.On("Recv").Once().Return(&statisticoproto.Team{}, io.EOF)

		teams, err := client.BySeasonID(ctx, 16036)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Equal(t, 2, len(teams))
		m.AssertExpectations(t)
		stream.AssertExpectations(t)
	})

	t.Run("logs error and returns internal server error if internal server error is returned by client", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoTeamClient)
		client := statisticodata.NewTeamClient(m)

		stream := new(MockTeamStream)

		ctx := context.Background()

		request := statisticoproto.SeasonTeamsRequest{SeasonId: 16036}

		e := status.Error(codes.Internal, "internal error")

		m.On("GetTeamsBySeasonId", ctx, &request, []grpc.CallOption(nil)).Return(stream, e)

		_, err := client.BySeasonID(ctx, 16036)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "internal server error returned from external service: rpc error: code = Internal desc = internal error", err.Error())
		m.AssertExpectations(t)
		stream.AssertNotCalled(t, "Recv")
	})

	t.Run("logs error and returns bad gateway error for non internal server error returned by client", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoTeamClient)
		client := statisticodata.NewTeamClient(m)

		stream := new(MockTeamStream)

		ctx := context.Background()

		request := statisticoproto.SeasonTeamsRequest{SeasonId: 16036}

		e := status.Error(codes.Unavailable, "service unavailable")

		m.On("GetTeamsBySeasonId", ctx, &request, []grpc.CallOption(nil)).Return(stream, e)

		_, err := client.BySeasonID(ctx, 16036)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "error connecting to external service: rpc error: code = Unavailable desc = service unavailable", err.Error())
		m.AssertExpectations(t)
		stream.AssertNotCalled(t, "Recv")
	})

	t.Run("logs error and returns internal server error if error reading from stream", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoTeamClient)
		client := statisticodata.NewTeamClient(m)

		stream := new(MockTeamStream)

		team := statisticoproto.Team{
			Id:        1,
			Name:      "West Ham United",
			CountryId: 8,
			VenueId:   214,
		}

		ctx := context.Background()

		request := statisticoproto.SeasonTeamsRequest{SeasonId: 16036}

		e := errors.New("oh damn")

		m.On("GetTeamsBySeasonId", ctx, &request, []grpc.CallOption(nil)).Return(stream, nil)
		stream.On("Recv").Twice().Return(&team, nil)
		stream.On("Recv").Once().Return(&statisticoproto.Team{}, e)

		_, err := client.BySeasonID(ctx, 16036)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "internal server error returned from external service: oh damn", err.Error())
		m.AssertExpectations(t)
		stream.AssertExpectations(t)
	})
}

type MockProtoTeamClient struct {
	mock.Mock
}

func (t *MockProtoTeamClient) GetTeamByID(ctx context.Context, in *statisticoproto.TeamRequest, opts ...grpc.CallOption) (*statisticoproto.Team, error) {
	args := t.Called(ctx, in, opts)
	return args.Get(0).(*statisticoproto.Team), args.Error(1)
}

func (t *MockProtoTeamClient) GetTeamsBySeasonId(ctx context.Context, in *statisticoproto.SeasonTeamsRequest, opts ...grpc.CallOption) (statisticoproto.TeamService_GetTeamsBySeasonIdClient, error) {
	args := t.Called(ctx, in, opts)
	return args.Get(0).(statisticoproto.TeamService_GetTeamsBySeasonIdClient), args.Error(1)
}

type MockTeamStream struct {
	mock.Mock
	grpc.ClientStream
}

func (t *MockTeamStream) Recv() (*statisticoproto.Team, error) {
	args := t.Called()
	return args.Get(0).(*statisticoproto.Team), args.Error(1)
}
