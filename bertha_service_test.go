package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const etag = "W/\"75e-78600296\""
const berthaPath = "/view/publish/gss/123456XYZ/Authors"

var martinWolf = author{
	Name:          "Martin Wolf",
	Email:         "martin.wolf@ft.com",
	ImageUrl:      "https://next-geebee.ft.com/image/v1/images/raw/fthead:martin-wolf?source=next",
	Biography:     "Martin Wolf is chief economics commentator at the Financial Times, London. He was awarded the CBE (Commander of the British Empire) in 2000 “for services to financial journalism”",
	TwitterHandle: "@martinwolf_",
	Uuid:          "daf5fed2-013c-468d-85c4-aee779b8aa53",
	TmeIdentifier: "Q0ItMDAwMDkwMA==-QXV0aG9ycw==",
}

var lucyKellaway = author{
	Name:          "Lucy Kellaway",
	Email:         "lucy.kellaway@ft.com",
	ImageUrl:      "https://next-geebee.ft.com/image/v1/images/raw/fthead:lucy-kellaway?source=next",
	Biography:     "Lucy Kellaway is an Associate Editor and management columnist of the FT. For the past 15 years her weekly Monday column has poked fun at management fads and jargon and celebrated the ups and downs of office life.",
	Uuid:          "daf5fed2-013c-468d-85c4-aee779b8aa51",
	TmeIdentifier: "Q0ItMDAwMDkyNg==-QXV0aG9ycw==",
}

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

func TestShouldReturnAuthorsCount(t *testing.T) {
	startBerthaMock("happy")
	defer berthaMock.Close()
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs := &berthaService{berthaUrl: spreadSheetUrl}

	c, err := bs.getAuthorsCount()

	assert.Nil(t, err)
	assert.Equal(t, 2, c, "Bertha should return 2 authors")
}

func TestShouldReturnAuthorsUuids(t *testing.T) {
	startBerthaMock("happy")
	defer berthaMock.Close()
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs := &berthaService{berthaUrl: spreadSheetUrl}

	bs.getAuthorsCount()
	uuids := bs.getAuthorsUuids()
	assert.Equal(t, 2, len(uuids), "Bertha should return 2 authors")
	assert.Equal(t, true, contains(uuids, martinWolf.Uuid), "It should contain Martin Wolf's UUID")
	assert.Equal(t, true, contains(uuids, lucyKellaway.Uuid), "It should contain Lucy Kellaway's UUID")
}

func TestShouldReturnSingleAuthor(t *testing.T) {
	startBerthaMock("happy")
	defer berthaMock.Close()
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs := &berthaService{berthaUrl: spreadSheetUrl}

	bs.getAuthorsCount()
	a := bs.getAuthorByUuid(martinWolf.Uuid)
	assert.Equal(t, martinWolf, a, "The author should be Martin Wolf")
}

func TestShouldReturnEmptyAuthorsUuidsWhenAuthorsCountIsNotCalled(t *testing.T) {
	startBerthaMock("happy")
	defer berthaMock.Close()
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs := &berthaService{berthaUrl: spreadSheetUrl}

	uuids := bs.getAuthorsUuids()
	assert.Equal(t, 0, len(uuids), "Bertha should return 0 authors")
}

func TestShouldReturnEmptyAuthorWhenAuthorsCountIsNotCalled(t *testing.T) {
	startBerthaMock("happy")
	defer berthaMock.Close()
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs := &berthaService{berthaUrl: spreadSheetUrl}

	a := bs.getAuthorByUuid(martinWolf.Uuid)
	assert.Equal(t, author{}, a, "The author should be empty")
}

func TestShouldReturnEmptyAuthorWhenAuthorIsNotAvailable(t *testing.T) {
	startBerthaMock("happy")
	defer berthaMock.Close()
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs := &berthaService{berthaUrl: spreadSheetUrl}

	bs.getAuthorsCount()
	a := bs.getAuthorByUuid("7f8bd61a-3575-4d32-a758-0fa41cbcc826")
	assert.Equal(t, author{}, a, "The author should be empty")
}

func TestShouldReturnErrorWhenBerthaIsUnhappy(t *testing.T) {
	startBerthaMock("unhappy")
	defer berthaMock.Close()
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs := &berthaService{berthaUrl: spreadSheetUrl}

	c, err := bs.getAuthorsCount()
	assert.NotNil(t, err)
	assert.Equal(t, -1, c, "It should return -1")

	authors := bs.getAuthorsUuids()
	assert.Equal(t, 0, len(authors), "It should return 0 authors")

	a := bs.getAuthorByUuid(martinWolf.Uuid)
	assert.Equal(t, author{}, a, "The author should be empty")
}

func TestCheckConnectivityOfHappyBerta(t *testing.T) {
	startBerthaMock("happy")
	defer berthaMock.Close()
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs := &berthaService{berthaUrl: spreadSheetUrl}

	c := bs.checkConnectivity()
	assert.Nil(t, c)
}

func TestCheckConnectivityOfUnhappyBertha(t *testing.T) {
	startBerthaMock("unhappy")
	defer berthaMock.Close()
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs := &berthaService{berthaUrl: spreadSheetUrl}

	c := bs.checkConnectivity()
	assert.NotNil(t, c)
}

func TestCheckConnectivityBerthaOffline(t *testing.T) {
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs := &berthaService{berthaUrl: spreadSheetUrl}

	c := bs.checkConnectivity()
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
