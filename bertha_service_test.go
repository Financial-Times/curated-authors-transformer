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

func TestShouldReturnAuthors(t *testing.T) {
	startBerthaMock("happy")
	defer berthaMock.Close()
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs := &berthaService{berthaUrl: spreadSheetUrl}

	uuids, err := bs.getAuthorsUuids()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(uuids), "Bertha should return 2 authors")
	assert.Equal(t, martinWolf.Uuid, uuids[0], "First authors should be Martin Wolf")
	assert.Equal(t, lucyKellaway.Uuid, uuids[1], "First authors should be Lucy Kellaway")
}

func TestShouldReturnSingleAuthor(t *testing.T) {
	startBerthaMock("happy")
	defer berthaMock.Close()
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs := &berthaService{berthaUrl: spreadSheetUrl}

	bs.getAuthorsUuids()
	a := bs.getAuthorByUuid(martinWolf.Uuid)
	assert.Equal(t, martinWolf, a, "The author should be Martin Wolf")
}

func TestShouldReturnEmptyAuthorWhenAuthorsUuidsAreNotFetched(t *testing.T) {
	startBerthaMock("happy")
	defer berthaMock.Close()
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs := &berthaService{berthaUrl: spreadSheetUrl}

	a := bs.getAuthorByUuid(martinWolf.Uuid)
	assert.Equal(t, author{}, a, "The author should be empty")
}

func TestShouldReturnErrorWhenBerthaIsUnhappy(t *testing.T) {
	startBerthaMock("unhappy")
	defer berthaMock.Close()
	spreadSheetUrl := berthaMock.URL + berthaPath
	bs := &berthaService{berthaUrl: spreadSheetUrl}

	authors, err := bs.getAuthorsUuids()
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(authors), "Bertha should return 0 authors")

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
