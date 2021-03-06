package medium

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/factly/bindu-server/util/test"
	"github.com/gavv/httpexpect/v2"
	"gopkg.in/h2non/gock.v1"
)

func TestMediumCreate(t *testing.T) {

	mock := test.SetupMockDB()

	testServer := httptest.NewServer(Routes())
	gock.New(testServer.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()
	defer testServer.Close()

	// create httpexpect instance
	e := httpexpect.New(t, testServer.URL)

	t.Run("cannot decode medium", func(t *testing.T) {

		e.POST(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)

	})

	t.Run("Unprocessable medium", func(t *testing.T) {

		e.POST(basePath).
			WithHeaders(headers).
			WithJSON(invalidData).
			Expect().
			Status(http.StatusUnprocessableEntity)

	})

	t.Run("create medium", func(t *testing.T) {

		slugCheckMock(mock)

		mediumInsertMock(mock)

		e.POST(basePath).
			WithHeaders(headers).
			WithJSON(data).
			Expect().
			Status(http.StatusCreated).JSON().Object().ContainsMap(data)
		test.ExpectationsMet(t, mock)

	})

	t.Run("create medium with slug is empty", func(t *testing.T) {

		slugCheckMock(mock)

		mediumInsertMock(mock)

		e.POST(basePath).
			WithHeaders(headers).
			WithJSON(mediumWithoutSlug).
			Expect().
			Status(http.StatusCreated).JSON().Object().ContainsMap(data)

		test.ExpectationsMet(t, mock)
	})

}
