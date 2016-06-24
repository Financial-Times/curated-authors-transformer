package main

import (
	"encoding/json"
	"fmt"
	"github.com/Financial-Times/go-fthealth/v1a"
	log "github.com/Sirupsen/logrus"
	"net/http"
)

type authorHandler struct {
	berthaService berthaService
}

func newAuthorHandler(bs berthaService) authorHandler {
	return authorHandler{berthaService: bs}
}

func (ah *authorHandler) getAuthors(writer http.ResponseWriter, req *http.Request) {
	berthaAuthors, err := ah.berthaService.getBerthaAuthors()
	if err != nil {
		writeJSONError(writer, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println(berthaAuthors)
	writeJSONResponse(berthaAuthors, len(berthaAuthors) > 0, writer)
}

func (ah *authorHandler) HealthCheck() v1a.Check {
	return v1a.Check{
		BusinessImpact:   "Unable to respond to request for curated author data from Bertha",
		Name:             "Check connectivity to Bertha",
		PanicGuide:       "https://sites.google.com/a/ft.com/ft-technology-service-transition/home/run-book-library/curated-authors-transfomer",
		Severity:         1,
		TechnicalSummary: "Cannot connect to Bertha to be able to supply curated authors",
		Checker:          ah.checker,
	}
}

func (ah *authorHandler) checker() (string, error) {
	err := ah.berthaService.checkConnectivity()
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
