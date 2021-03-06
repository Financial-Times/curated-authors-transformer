package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var curatedAuthorsTransformer *httptest.Server
var expectedStreamOutput = `{"id":"` + martinWolfUuid + `"} {"id":"` + lucyKellawayUuid + `"} `

type MockedBerthaService struct {
	mock.Mock
}

func (m *MockedBerthaService) refreshCache() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockedBerthaService) getAuthorsUuids() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockedBerthaService) getAuthorByUuid(uuid string) person {
	args := m.Called(uuid)
	return args.Get(0).(person)
}

func (m *MockedBerthaService) getAuthorsCount() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockedBerthaService) checkConnectivity() error {
	args := m.Called()
	return args.Error(0)
}

func startCuratedAuthorsTransformer(bs *MockedBerthaService) {
	ah := newAuthorHandler(bs)
	h := setupServiceHandlers(ah)
	curatedAuthorsTransformer = httptest.NewServer(h)
}

func TestShouldReturn200AndAuthorsCount(t *testing.T) {
	mbs := new(MockedBerthaService)
	mbs.On("getAuthorsCount").Return(2)
	mbs.On("refreshCache").Return(nil)
	startCuratedAuthorsTransformer(mbs)
	defer curatedAuthorsTransformer.Close()

	resp, err := http.Get(curatedAuthorsTransformer.URL + "/transformers/authors/__count")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")
	assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), "Content-Type should be text/plain")
	actualOutput := getStringFromReader(resp.Body)
	assert.Equal(t, "2", actualOutput, "Response body should contain the count of available authors")
}

func TestShouldReturn500WhenAuthorsCountIsCalledAndCacheRefreshFails(t *testing.T) {
	mbs := new(MockedBerthaService)
	mbs.On("getAuthorsCount").Return(2)
	mbs.On("refreshCache").Return(errors.New("I hate Luca!"))
	startCuratedAuthorsTransformer(mbs)
	defer curatedAuthorsTransformer.Close()

	resp, err := http.Get(curatedAuthorsTransformer.URL + "/transformers/authors/__count")
	defer resp.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Response status should be 500")
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type should be application/json")
	actualOutput := getStringFromReader(resp.Body)
	assert.Equal(t, "{\"message\": \"I hate Luca!\"}\n", actualOutput, "Response body should contain the error message as JSON")
}

func TestShouldReturn200AndAuthorsUuids(t *testing.T) {
	mbs := new(MockedBerthaService)
	mbs.On("getAuthorsUuids").Return(expectedUuids)
	startCuratedAuthorsTransformer(mbs)
	defer curatedAuthorsTransformer.Close()

	resp, err := http.Get(curatedAuthorsTransformer.URL + "/transformers/authors/__ids")
	defer resp.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")
	assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), "Content-Type should be text/plain")
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
	mbs.On("getAuthorByUuid", martinWolfUuid).Return(transformedMartinWolf)
	startCuratedAuthorsTransformer(mbs)
	defer curatedAuthorsTransformer.Close()

	resp, err := http.Get(curatedAuthorsTransformer.URL + "/transformers/authors/" + martinWolfUuid)
	defer resp.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type should be application/json")

	file, _ := os.Open("test-resources/martin-wolf-transformed-output.json")
	defer file.Close()

	expectedOutput := getStringFromReader(file)
	actualOutput := getStringFromReader(resp.Body)

	assert.Equal(t, expectedOutput, actualOutput, "Response body should be Martin Wolf")
}

func TestShouldReturn404WhenAuthorIsNotFound(t *testing.T) {
	mbs := new(MockedBerthaService)
	mbs.On("getAuthorByUuid", martinWolfUuid).Return(person{})
	startCuratedAuthorsTransformer(mbs)
	defer curatedAuthorsTransformer.Close()

	resp, err := http.Get(curatedAuthorsTransformer.URL + "/transformers/authors/" + martinWolfUuid)
	defer resp.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Response status should be 404")
}
