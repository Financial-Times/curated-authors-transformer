package main

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gregjones/httpcache"
	"net/http"
	"sync"
)

var client = httpcache.NewMemoryCacheTransport().Client()

type berthaService struct {
	berthaUrl   string
	authorsMap  map[string]person
	transformer transformer
	mutex       *sync.Mutex
}

func newBerthaService(url string) (*berthaService, error) {
	bs := &berthaService{
		berthaUrl:   url,
		authorsMap:  map[string]person{},
		transformer: &berthaTransformer{},
		mutex:       &sync.Mutex{},
	}
	err := bs.refreshCache()
	return bs, err
}

func (bs *berthaService) refreshCache() error {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	bs.authorsMap = make(map[string]person)

	authors, err := bs.getAuthors()

	if err != nil {
		return err
	}

	for _, a := range authors {
		p, transErr := bs.transformer.authorToPerson(a)
		if transErr != nil {
			log.Error(err)
			return transErr
		}
		bs.authorsMap[p.Uuid] = p
	}
	return nil
}

func (bs *berthaService) getAuthors() ([]author, error) {
	resp, err := bs.callBerthaService()
	if err != nil {
		log.Error(err)
		return []author{}, err
	}

	var authors []author
	if err = json.NewDecoder(resp.Body).Decode(&authors); err != nil {
		log.Error(err)
		return []author{}, err
	}
	return authors, nil
}

func (bs *berthaService) callBerthaService() (res *http.Response, err error) {
	log.WithFields(log.Fields{"bertha_url": bs.berthaUrl}).Info("Calling Bertha...")
	res, err = client.Get(bs.berthaUrl)
	return
}

func (bs *berthaService) getAuthorsCount() int {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	return len(bs.authorsMap)
}

func (bs *berthaService) getAuthorsUuids() []string {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	uuids := make([]string, 0)
	for uuid, _ := range bs.authorsMap {
		uuids = append(uuids, uuid)
	}
	return uuids
}

func (bs *berthaService) getAuthorByUuid(uuid string) person {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	return bs.authorsMap[uuid]
}

func (bs *berthaService) checkConnectivity() error {
	resp, err := bs.callBerthaService()
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Bertha returns unexpected HTTP status: %d", resp.StatusCode))
	}
	return nil
}
