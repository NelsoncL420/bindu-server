package organisation

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/factly/bindu-server/util"
	"github.com/factly/bindu-server/util/test"
	"github.com/factly/x/loggerx"
	"github.com/gavv/httpexpect/v2"
	"github.com/go-chi/chi"
	"gopkg.in/h2non/gock.v1"
)

var basePath = "/organisations"

var orgList = []map[string]interface{}{{
	"id":         1,
	"deleted_at": nil,
	"title":      "test org",
	"slug":       "tesing",
	"permission": map[string]interface{}{
		"id":         1,
		"deleted_at": nil,
		"role":       "owner",
	}},
}

func routes() http.Handler {
	r := chi.NewRouter()
	r.Use(loggerx.Init())
	r.With(util.CheckUser, util.CheckOrganisation).Mount(basePath, Router())
	return r
}

var headers = map[string]string{
	"X-Organisation": "1",
	"X-User":         "1",
}

func TestOrganisationList(t *testing.T) {
	test.SetEnv()
	defer gock.Disable()
	test.MockServer()
	defer gock.DisableNetworking()

	testServer := httptest.NewServer(routes())
	gock.New(testServer.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()
	defer testServer.Close()

	// create httpexpect instance
	e := httpexpect.New(t, testServer.URL)

	t.Run("get empty list of organisations", func(t *testing.T) {

		e.GET(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 1}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(orgList[0])

	})
}
