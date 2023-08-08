package main

import (
	"net/http"
	"net/http/httptest"
	"score/model"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

type MockDB struct {
	mock.Mock
}

var _ model.Repo = (*MockDB)(nil)

func (m *MockDB) GetAllPlayerRankwise() ([]model.Player, error) {
	args := m.Called()
	return args.Get(0).([]model.Player), args.Error(1)
}

func (m *MockDB) GetPlayerByRank(r int) (model.Player, error) {
	args := m.Called(r)
	return args.Get(0).(model.Player), args.Error(1)
}

func (m *MockDB) GetRandomPlayer() (model.Player, error) {
	args := m.Called()
	return args.Get(0).(model.Player), args.Error(1)
}

func (m *MockDB) DeletePlayer(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDB) CreateNewPlayer(p *model.Player) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockDB) UpdatePlayer(id int, name string, score int, check map[string]interface{}) (model.Player, error) {
	args := m.Called(id, name, score, check)
	return args.Get(0).(model.Player), args.Error(1)
}

var (
	Db = &MockDB{}
)

func TestGetAllPlayerRankwise(t *testing.T) {
	tcs := []struct {
		name       string
		want       int
		players    []model.Player
		serviceErr error
	}{
		{
			name:       "empty players",
			want:       http.StatusOK,
			players:    make([]model.Player, 0),
			serviceErr: model.ErrNoPlayer,
		},
		{
			name: "database failing",
			want: http.StatusInternalServerError,
			players: []model.Player{
				{
					Id:      1,
					Name:    "abhiraj ranjan",
					Score:   1212,
					Country: "IN",
				},
			},
			serviceErr: errors.New("database failing"),
		},
		{
			name: "no error",
			want: http.StatusOK,
			players: []model.Player{
				{
					Id:      1,
					Name:    "abhiraj ranjan",
					Score:   1212,
					Country: "IN",
				},
			},
		},
	}

	for _, tc := range tcs {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		t.Run(tc.name, func(t *testing.T) {

			call := Db.On("GetAllPlayerRankwise").Return(tc.players, tc.serviceErr)
			GetAllPlayerRankwise(c, Db)
			call.Unset()

			assert.Equal(t, tc.want, w.Code)
		})
	}
}

func TestGetPlayerByRank(t *testing.T) {
	tcs := []struct {
		name              string
		rank              any
		want              int
		isPostServiceTest bool
		players           model.Player
		serviceErr        error
	}{
		{
			name: "parameter not int",
			rank: "not int",
			want: http.StatusBadRequest,
		},
		{
			name: "parameter less than 0",
			rank: -1,
			want: http.StatusBadRequest,
		},
		{
			name: "parameter is 0",
			rank: 0,
			want: http.StatusBadRequest,
		},
		{
			name:              "no players",
			isPostServiceTest: true,
			rank:              1,
			want:              http.StatusNotFound,
			players:           model.Player{},
			serviceErr:        model.ErrRankDoensnotExist,
		},
		{
			name:              "service error",
			isPostServiceTest: true,
			rank:              1,
			want:              http.StatusInternalServerError,
			players:           model.Player{},
			serviceErr:        errors.New("service error"),
		},
		{
			name:              "no error",
			isPostServiceTest: true,
			rank:              1,
			want:              http.StatusOK,
			players: model.Player{
				Id:      1,
				Name:    "abhiraj",
				Country: "IN",
				Score:   133,
			},
		},
	}

	for _, tc := range tcs {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		s, ok := tc.rank.(string)
		if ok {
			c.AddParam("val", s)
		}

		i, ok := tc.rank.(int)
		if ok {
			s := strconv.Itoa(i)
			c.AddParam("val", s)
		}
		t.Log(c.Params)

		t.Run(tc.name, func(t *testing.T) {
			if tc.isPostServiceTest {
				call := Db.On("GetPlayerByRank", tc.rank).Return(tc.players, tc.serviceErr)
				GetPlayerByRank(c, Db)
				call.Unset()
			} else {
				GetPlayerByRank(c, Db)
			}

			assert.Equal(t, tc.want, w.Code)
		})
	}
}

func TestGetRandomPlayer(t *testing.T) {
	tcs := []struct {
		name       string
		player     model.Player
		serviceErr error
		want       int
	}{
		{
			name:       "service error",
			player:     model.Player{},
			serviceErr: errors.New("service error"),
			want:       http.StatusInternalServerError,
		},
		{
			name: "no error",
			player: model.Player{
				Id:      1,
				Name:    "abhiraj",
				Country: "IN",
				Score:   1221,
			},
			want: http.StatusOK,
		},
	}

	for _, tc := range tcs {

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			call := Db.On("GetRandomPlayer").Return(tc.player, tc.serviceErr)
			GetRandomPlayer(c, Db)
			call.Unset()

			assert.Equal(t, tc.want, w.Code)
		})
	}
}
