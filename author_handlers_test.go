package main

import (
	"bytes"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var curatedAuthorsTransformer *httptest.Server

var uuids = []string{martinWolf.Uuid, lucyKellaway.Uuid}
var expectedStreamOutput = `{"id":"` + martinWolf.Uuid + `"} {"id":"` + lucyKellaway.Uuid + `"} `

var martinWolfIdentifier = identifier{tmeAuthority, "Q0ItMDAwMDkwMA==-QXV0aG9ycw=="}

var transformedMartinWolf = person{
	Uuid:           "daf5fed2-013c-468d-85c4-aee779b8aa53",
	Name:           "Martin Wolf",
	EmailAddress:   "martin.wolf@ft.com",
	TwitterHandle:  "@martinwolf_",
	Description:    "Martin Wolf is chief economics commentator at the Financial Times, London.",
	DescriptionXML: `<p>Martin Wolf is chief economics commentator at the Financial Times, London.</p>`,
	ImageUrl:       "https://next-geebee.ft.com/image/v1/images/raw/fthead:martin-wolf?source=next",
	Identifiers:    []identifier{martinWolfIdentifier},
}

type MockedBerthaService struct {
	mock.Mock
}

func (m *MockedBerthaService) getAuthorsUuids() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockedBerthaService) getAuthorByUuid(uuid string) author {
	args := m.Called(uuid)
	return args.Get(0).(author)
}

func (m *MockedBerthaService) getAuthorsCount() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockedBerthaService) checkConnectivity() error {
	args := m.Called()
	return args.Error(0)
}

type MockedTransformer struct {
	mock.Mock
}

func (m *MockedTransformer) authorToPerson(a author) (person, error) {
	args := m.Called(a)
	return args.Get(0).(person), args.Error(1)
}

func startCuratedAuthorsTransformer(bs *MockedBerthaService, mt *MockedTransformer) {
	ah := authorHandler{
		authorsService: bs,
		transformer:    mt,
	}
	r := mux.NewRouter()
	r.HandleFunc("/transformers/authors/__count", ah.getAuthorsCount).Methods("GET")
	r.HandleFunc("/transformers/authors/__ids", ah.getAuthorsUuids).Methods("GET")
	r.HandleFunc("/transformers/authors/{uuid}", ah.getAuthorByUuid).Methods("GET")
	curatedAuthorsTransformer = httptest.NewServer(r)
}

func TestShouldReturn200AndAuthorsCount(t *testing.T) {

	mbs := new(MockedBerthaService)
	mbs.On("getAuthorsCount").Return(2, nil)
	mt := new(MockedTransformer)
	startCuratedAuthorsTransformer(mbs, mt)
	defer curatedAuthorsTransformer.Close()

	resp, err := http.Get(curatedAuthorsTransformer.URL + "/transformers/authors/__count")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")

	actualOutput := getStringFromReader(resp.Body)
	assert.Equal(t, "2\n", actualOutput, "Response body should contain the count of available authors")
}

func TestShouldReturn200AndAuthorsUuids(t *testing.T) {

	mbs := new(MockedBerthaService)
	mbs.On("getAuthorsUuids").Return(uuids)
	mt := new(MockedTransformer)
	startCuratedAuthorsTransformer(mbs, mt)
	defer curatedAuthorsTransformer.Close()

	resp, err := http.Get(curatedAuthorsTransformer.URL + "/transformers/authors/__ids")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")

	actualOutput := getStringFromReader(resp.Body)
	assert.Equal(t, expectedStreamOutput, actualOutput, "Response body should be a sequence of ids")
}

func getStringFromReader(r io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return buf.String()
}

func TestShouldReturn200AndTrasformedAuthor(t *testing.T) {

	mbs := new(MockedBerthaService)
	mbs.On("getAuthorByUuid", martinWolf.Uuid).Return(martinWolf)
	mt := new(MockedTransformer)
	mt.On("authorToPerson", martinWolf).Return(transformedMartinWolf, nil)
	startCuratedAuthorsTransformer(mbs, mt)
	defer curatedAuthorsTransformer.Close()

	resp, err := http.Get(curatedAuthorsTransformer.URL + "/transformers/authors/" + martinWolf.Uuid)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")

	file, _ := os.Open("test-resources/martin-wolf-transformed-output.json")
	defer file.Close()

	expectedOutput := getStringFromReader(file)
	actualOutput := getStringFromReader(resp.Body)

	assert.Equal(t, expectedOutput, actualOutput, "Response body should be Martin Wolf")
}

func TestShouldReturn404WhenAuthorIsNotFound(t *testing.T) {
	mbs := new(MockedBerthaService)
	mbs.On("getAuthorByUuid", martinWolf.Uuid).Return(author{})
	mt := new(MockedTransformer)
	mt.On("authorToPerson", author{}).Return(person{}, nil)
	startCuratedAuthorsTransformer(mbs, mt)
	defer curatedAuthorsTransformer.Close()

	resp, err := http.Get(curatedAuthorsTransformer.URL + "/transformers/authors/" + martinWolf.Uuid)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Response status should be 404")
}

func TestShouldReturn500WhenBerthaReturnsError(t *testing.T) {
	mbs := new(MockedBerthaService)
	mbs.On("getAuthorsCount").Return(-1, errors.New("I am a zobie"))
	mt := new(MockedTransformer)
	startCuratedAuthorsTransformer(mbs, mt)
	defer curatedAuthorsTransformer.Close()

	resp, err := http.Get(curatedAuthorsTransformer.URL + "/transformers/authors/__count")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Response status should be 500")
}

func TestShouldReturn500WhenTransformerReturnsError(t *testing.T) {
	mbs := new(MockedBerthaService)
	mbs.On("getAuthorByUuid", martinWolf.Uuid).Return(martinWolf)
	mt := new(MockedTransformer)
	mt.On("authorToPerson", martinWolf).Return(person{}, errors.New("I hate Luca!!!"))
	startCuratedAuthorsTransformer(mbs, mt)
	defer curatedAuthorsTransformer.Close()

	resp, err := http.Get(curatedAuthorsTransformer.URL + "/transformers/authors/" + martinWolf.Uuid)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Response status should be 500")
}
