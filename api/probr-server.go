package api

import (	
	"sync"	
	"net/http"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	api "citihub.com/probr/api/probrapi"
	"citihub.com/probr/internal/coreengine"
)

//TestStore ...
type TestStore struct {
	Tests map[uuid.UUID]*api.ProbrTest
	Lock  sync.Mutex
}

//NewProbrAPI ...
func NewProbrAPI() *TestStore {
	return &TestStore{
		Tests: make(map[uuid.UUID]*api.ProbrTest),
	}
}

//GetTests - returns all available tests.
func (t *TestStore) GetTests(ctx echo.Context) error {

	p := coreengine.GetAvailableTests()
	return ctx.JSONPretty(http.StatusOK, p, " ")
}

