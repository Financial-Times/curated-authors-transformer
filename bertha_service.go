package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gregjones/httpcache"
	"net/http"
)

var client = httpcache.NewMemoryCacheTransport().Client()

type berthaService struct {
	berthaUrl string
}

func newBerthaService(url string) berthaService {
	return berthaService{berthaUrl: url}
}

func (bs *berthaService) getBerthaAuthors() ([]berthaAuthor, error) {
	resp, err := bs.callBerthaService()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	fmt.Println(resp)

	var authors []berthaAuthor
	fmt.Println(authors)
	if err = json.NewDecoder(resp.Body).Decode(&authors); err != nil {
		log.Error(err)
		return nil, err
	}
	fmt.Println(authors)
	return authors, nil
}

func (bs *berthaService) getBerthaAuthorsByUuid(uuid string) (berthaAuthor, error) {
	authors, err := getBerthaAuthors()
	if err != nil {
		return nil, err
	}
	for _, author := range authors {
		if author.Uuid == uuid {
			return author, nil
		}
	}
	return nil, nil
}

func (bs *berthaService) callBerthaService() (*http.Response, error) {
	return client.Get(bs.berthaUrl)
}

func (bs *berthaService) checkConnectivity() (err error) {
	_, err = bs.callBerthaService()
	return
}
