package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const etag = "W/\"75e-78600296\""
const berthaPath = "/view/publish/gss/123456XYZ/Authors"

var berthaMock *httptest.Server

func startBerthaMock(status string) {
	r := mux.NewRouter()
	if status == "happy" {
		r.Path(berthaPath).Handler(handlers.MethodHandler{"GET": http.HandlerFunc(berthaHandlerMock)})
	} else {
		r.Path(berthaPath).Handler(handlers.MethodHandler{"GET": http.HandlerFunc(unhappyHandler)})
	}
	berthaMock = httptest.NewServer(r)
}

func berthaHandlerMock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	ifNoneMatch := r.Header.Get("If-None-Match")

	if ifNoneMatch == etag {
		w.WriteHeader(http.StatusNotModified)
	} else {
		w.Header().Set("ETag", etag)

		file, err := os.Open("test-resources/bertha-output.json")
		if err != nil {
			return
		}
		defer file.Close()
		io.Copy(w, file)
	}
}

func unhappyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

type MockedTransformer struct {
	mock.Mock
}

func (m *MockedTransformer) authorToPerson(a author) (person, error) {
	args := m.Called(a)
	return args.Get(0).(person), args.Error(1)
}

func TestShouldReturnAuthorsCount(t *testing.T) {
	startBerthaMock("happy")
	defer berthaMock.Close()
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs, err := newBerthaService(spreadSheetUrl)

	c := bs.getAuthorsCount()

	assert.Nil(t, err)
	assert.Equal(t, 2, c, "Bertha should return 2 authors")
}

func TestShouldReturnAuthorsUuids(t *testing.T) {
	startBerthaMock("happy")
	defer berthaMock.Close()
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs, err := newBerthaService(spreadSheetUrl)

	uuids := bs.getAuthorsUuids()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(uuids), "Bertha should return 2 authors")
	assert.Equal(t, true, contains(uuids, martinWolfUuid), "It should contain Martin Wolf's UUID")
	assert.Equal(t, true, contains(uuids, lucyKellawayUuid), "It should contain Lucy Kellaway's UUID")
}

func TestShouldReturnSingleAuthor(t *testing.T) {
	startBerthaMock("happy")
	defer berthaMock.Close()
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs, err := newBerthaService(spreadSheetUrl)

	a := bs.getAuthorByUuid(martinWolfUuid)

	assert.Nil(t, err)
	assert.Equal(t, transformedMartinWolf, a, "The author should be Martin Wolf")
}

func TestShouldReturnEmptyAuthorWhenAuthorIsNotAvailable(t *testing.T) {
	startBerthaMock("happy")
	defer berthaMock.Close()
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs, err := newBerthaService(spreadSheetUrl)

	bs.getAuthorsCount()

	assert.Nil(t, err)
	a := bs.getAuthorByUuid("7f8bd61a-3575-4d32-a758-0fa41cbcc826")
	assert.Equal(t, person{}, a, "The author should be empty")
}

func TestShouldReturnErrorWhenBerthaIsUnhappy(t *testing.T) {
	startBerthaMock("unhappy")
	defer berthaMock.Close()
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs, err := newBerthaService(spreadSheetUrl)
	assert.NotNil(t, err)

	c := bs.getAuthorsCount()
	assert.Equal(t, 0, c, "It should return 0")

	authors := bs.getAuthorsUuids()
	assert.Equal(t, 0, len(authors), "It should return 0 authors")

	a := bs.getAuthorByUuid(martinWolfUuid)
	assert.Equal(t, person{}, a, "The author should be empty")
}

func TestCheckConnectivityOfHappyBerta(t *testing.T) {
	startBerthaMock("happy")
	defer berthaMock.Close()
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs, err := newBerthaService(spreadSheetUrl)

	c := bs.checkConnectivity()
	assert.Nil(t, err)
	assert.Nil(t, c)
}

func TestCheckConnectivityOfUnhappyBertha(t *testing.T) {
	startBerthaMock("unhappy")
	defer berthaMock.Close()
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs, err := newBerthaService(spreadSheetUrl)

	c := bs.checkConnectivity()
	assert.NotNil(t, err)
	assert.NotNil(t, c)
}

func TestCheckConnectivityBerthaOffline(t *testing.T) {
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs, err := newBerthaService(spreadSheetUrl)

	c := bs.checkConnectivity()
	assert.NotNil(t, err)
	assert.NotNil(t, c)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
