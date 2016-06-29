package main

import (
	"bytes"
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
	DescriptionXML: `<p>Martin Wolf is chief economics commentator at the Financial Times, London.<p>`,
	ImageUrl:       "https://upload.wikimedia.org/wikipedia/en/7/77/EricCartman.png",
	Identifiers:    []identifier{martinWolfIdentifier},
}

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
	h := setupServiceHandlers(ah)
	curatedAuthorsTransformer = httptest.NewServer(h)
}

func TestShouldReturn200AndTrasformedAuthors(t *testing.T) {

	mbs := new(MockedBerthaService)
	mbs.On("getAuthorsUuids").Return(uuids, nil)
	mt := new(MockedTransformer)
	mt.On("authorToPerson", martinWolf).Return(transformedMartinWolf, nil)
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
