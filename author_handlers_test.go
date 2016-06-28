package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var curatedAuthorsTransformer *httptest.Server

var uuids = []string{martinWolf.Uuid, lucyKellaway.Uuid}

type MockedBerthaService struct {
	mock.Mock
}

func (m *MockedBerthaService) getAuthorsUuids() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockedBerthaService) getAuthorByUuid(uuid string) author {
	args := m.Called(uuid)
	return args.Get(0).(author)
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

func TestShouldReturn200AndTrasformedAuthors(t *testing.T) {

	mbs := new(MockedBerthaService)
	mbs.On("getAuthorsUuids").Return(uuids, nil)
	startCuratedAuthorsTransformer(mbs)
	defer curatedAuthorsTransformer.Close()

	resp, err := http.Get(curatedAuthorsTransformer.URL + "/transformers/authors/__ids")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")

	//TODO Implement propert test acording to expected output
	//  file, _ := os.Open("test-resources/transformer-output.json")
	// 	defer file.Close()
	//
	// 	expectedOutput := getStringFromReader(file)
	// 	actualOutput := getStringFromReader(resp.Body)
	//
	// 	assert.Equal(t, expectedOutput, actualOutput, "Response body shoud be equal to transformer response body")
}

func getStringFromReader(r io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return buf.String()
}

func TestShouldReturn200AndTrasformedAuthor(t *testing.T) {

	mbs := new(MockedBerthaService)
	mbs.On("getAuthorByUuid", martinWolf.Uuid).Return(martinWolf)
	startCuratedAuthorsTransformer(mbs)
	defer curatedAuthorsTransformer.Close()

	resp, err := http.Get(curatedAuthorsTransformer.URL + "/transformers/authors/" + martinWolf.Uuid)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")

	//TODO Implement propert test acording to expected output
	//  file, _ := os.Open("test-resources/transformer-output.json")
	// 	defer file.Close()
	//
	// 	expectedOutput := getStringFromReader(file)
	// 	actualOutput := getStringFromReader(resp.Body)
	//
	// 	assert.Equal(t, expectedOutput, actualOutput, "Response body shoud be equal to transformer response body")
}
