package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Financial-Times/go-fthealth/v1a"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
)

type authorHandler struct {
	authorsService authorsService
}

func newAuthorHandler(as authorsService) authorHandler {
	return authorHandler{
		authorsService: as,
	}
}

func (ah *authorHandler) getAuthorsCount(writer http.ResponseWriter, req *http.Request) {
	c, err := ah.authorsService.getAuthorsCount()
	if err != nil {
		writeJSONError(writer, err.Error(), http.StatusInternalServerError)
	} else {
		var buffer bytes.Buffer
		buffer.WriteString(fmt.Sprintf(`%v`, c))
		buffer.WriteTo(writer)
	}
}

func (ah *authorHandler) getAuthorsUuids(writer http.ResponseWriter, req *http.Request) {
	uuids := ah.authorsService.getAuthorsUuids()
	writeStreamResponse(uuids, writer)
}

func (ah *authorHandler) getAuthorByUuid(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]

	a := ah.authorsService.getAuthorByUuid(uuid)
	writeJSONResponse(a, !reflect.DeepEqual(a, person{}), writer)
}

func (ah *authorHandler) HealthCheck() v1a.Check {
	return v1a.Check{
		BusinessImpact:   "Unable to respond to request for curated author data from Bertha",
		Name:             "Check connectivity to Bertha",
		PanicGuide:       "https://sites.google.com/a/ft.com/ft-technology-service-transition/home/run-book-library/curated-authors-transformer",
		Severity:         1,
		TechnicalSummary: "Cannot connect to Bertha to be able to supply curated authors",
		Checker:          ah.checker,
	}
}

func (ah *authorHandler) checker() (string, error) {
	err := ah.authorsService.checkConnectivity()
	if err == nil {
		return "Connectivity to Bertha is ok", err
	}
	return "Error connecting to Bertha", err
}

func (ah *authorHandler) GoodToGo(writer http.ResponseWriter, req *http.Request) {
	if _, err := ah.checker(); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
	}
}

func writeJSONResponse(obj interface{}, found bool, writer http.ResponseWriter) {
	writer.Header().Add("Content-Type", "application/json")

	if !found {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	enc := json.NewEncoder(writer)
	if err := enc.Encode(obj); err != nil {
		log.Errorf("Error on json encoding=%v\n", err)
		writeJSONError(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func writeJSONError(w http.ResponseWriter, errorMsg string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, fmt.Sprintf("{\"message\": \"%s\"}", errorMsg))
}

func writeStreamResponse(ids []string, writer http.ResponseWriter) {
	for _, id := range ids {
		var buffer bytes.Buffer
		buffer.WriteString(fmt.Sprintf(`{"id":"%s"} `, id))
		buffer.WriteTo(writer)
	}
}
